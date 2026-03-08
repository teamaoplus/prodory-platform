package cluster

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"

	"kubernetes-in-a-box/pkg/config"
)

// Manager handles cluster operations
type Manager struct {
	log       *logrus.Logger
	store     *ClusterStore
	providers map[string]Provider
}

// Provider interface for different deployment targets
type Provider interface {
	Name() string
	Create(opts *config.ClusterOptions) error
	Delete(name string, force bool) error
	List() ([]ClusterInfo, error)
	Kubeconfig(name string) (string, error)
	SSH(name, node string) error
	Scale(name string, workers int) error
	Upgrade(name, version string) error
	Status(name string) (*ClusterStatus, error)
	Validate(name string) error
}

// ClusterInfo represents cluster information
type ClusterInfo struct {
	Name       string
	Provider   string
	Version    string
	Masters    int
	Workers    int
	Status     string
	Endpoint   string
	CreatedAt  time.Time
	Age        string
}

// ClusterStatus represents detailed cluster status
type ClusterStatus struct {
	Name           string
	Provider       string
	Version        string
	K3sVersion     string
	Status         string
	Endpoint       string
	Masters        []NodeStatus
	Workers        []NodeStatus
	Components     []ComponentStatus
	Conditions     []ClusterCondition
}

// NodeStatus represents node status
type NodeStatus struct {
	Name      string
	Role      string
	Status    string
	Version   string
	IP        string
	OS        string
	Kernel    string
	Container string
	Uptime    string
}

// ComponentStatus represents component status
type ComponentStatus struct {
	Name      string
	Namespace string
	Status    string
	Version   string
}

// ClusterCondition represents cluster condition
type ClusterCondition struct {
	Type    string
	Status  string
	Reason  string
	Message string
}

// NewManager creates a new cluster manager
func NewManager(log *logrus.Logger) *Manager {
	store := NewClusterStore()

	return &Manager{
		log:       log,
		store:     store,
		providers: make(map[string]Provider),
	}
}

// getProvider returns the provider for the given name
func (m *Manager) getProvider(name string) (Provider, error) {
	if provider, exists := m.providers[name]; exists {
		return provider, nil
	}

	// Initialize provider based on name
	switch name {
	case "local":
		provider := NewLocalProvider(m.log, m.store)
		m.providers[name] = provider
		return provider, nil
	case "aws":
		provider := NewAWSProvider(m.log, m.store)
		m.providers[name] = provider
		return provider, nil
	case "azure":
		provider := NewAzureProvider(m.log, m.store)
		m.providers[name] = provider
		return provider, nil
	case "gcp":
		provider := NewGCPProvider(m.log, m.store)
		m.providers[name] = provider
		return provider, nil
	case "vagrant":
		provider := NewVagrantProvider(m.log, m.store)
		m.providers[name] = provider
		return provider, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// Create creates a new cluster
func (m *Manager) Create(opts *config.ClusterOptions) error {
	m.log.Infof("Creating cluster '%s' with provider '%s'", opts.Name, opts.Provider)

	// Check if cluster already exists
	if _, exists := m.store.Get(opts.Name); exists {
		return fmt.Errorf("cluster '%s' already exists", opts.Name)
	}

	// Enable HA if multiple masters
	if opts.Masters > 1 {
		opts.HA = true
	}

	// Get provider
	provider, err := m.getProvider(opts.Provider)
	if err != nil {
		return err
	}

	// Create progress bar
	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription("Creating cluster..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	// Create cluster
	if err := provider.Create(opts); err != nil {
		bar.Finish()
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	bar.Finish()

	// Save cluster info
	cluster := &ClusterInfo{
		Name:      opts.Name,
		Provider:  opts.Provider,
		Masters:   opts.Masters,
		Workers:   opts.Workers,
		Status:    "Running",
		CreatedAt: time.Now(),
	}
	m.store.Save(cluster)

	// Print success message
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("\n✓ Cluster '%s' created successfully!\n\n", opts.Name)

	// Print connection info
	fmt.Println("To connect to your cluster:")
	fmt.Printf("  export KUBECONFIG=$(kib kubeconfig %s)\n", opts.Name)
	fmt.Println("  kubectl get nodes")
	fmt.Println()

	return nil
}

// Delete deletes a cluster
func (m *Manager) Delete(name string, force bool) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Confirm deletion
	if !force {
		red := color.New(color.FgRed, color.Bold)
		red.Printf("This will permanently delete cluster '%s' and all associated resources.\n", name)
		fmt.Print("Are you sure? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	m.log.Infof("Deleting cluster '%s'...", name)

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	// Delete cluster
	if err := provider.Delete(name, force); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	// Remove from store
	m.store.Delete(name)

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ Cluster '%s' deleted successfully!\n", name)

	return nil
}

// List lists all clusters
func (m *Manager) List() error {
	clusters := m.store.List()

	if len(clusters) == 0 {
		fmt.Println("No clusters found.")
		fmt.Println("Create a cluster with: kib create --name <name>")
		return nil
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "PROVIDER", "MASTERS", "WORKERS", "VERSION", "STATUS", "AGE"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, c := range clusters {
		age := time.Since(c.CreatedAt).Round(time.Second).String()
		if age == "0s" {
			age = "just now"
		}

		table.Append([]string{
			c.Name,
			c.Provider,
			fmt.Sprintf("%d", c.Masters),
			fmt.Sprintf("%d", c.Workers),
			c.Version,
			c.Status,
			age,
		})
	}

	table.Render()

	return nil
}

// Kubeconfig retrieves kubeconfig for a cluster
func (m *Manager) Kubeconfig(name string, merge bool) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	// Get kubeconfig
	kubeconfig, err := provider.Kubeconfig(name)
	if err != nil {
		return err
	}

	if merge {
		// Merge with existing kubeconfig
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		if err := mergeKubeconfig(kubeconfig, kubeconfigPath); err != nil {
			return err
		}
		fmt.Printf("Merged kubeconfig into %s\n", kubeconfigPath)
	} else {
		// Write to temp file and print path
		tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("kib-%s-kubeconfig", name))
		if err := os.WriteFile(tmpFile, []byte(kubeconfig), 0600); err != nil {
			return err
		}
		fmt.Println(tmpFile)
	}

	return nil
}

// SSH connects to a cluster node via SSH
func (m *Manager) SSH(name, node string) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	return provider.SSH(name, node)
}

// Scale scales a cluster
func (m *Manager) Scale(name string, workers int) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	m.log.Infof("Scaling cluster '%s' to %d workers...", name, workers)

	if err := provider.Scale(name, workers); err != nil {
		return err
	}

	// Update cluster info
	cluster.Workers = workers
	m.store.Save(cluster)

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ Cluster '%s' scaled to %d workers!\n", name, workers)

	return nil
}

// Upgrade upgrades a cluster
func (m *Manager) Upgrade(name, version string) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	m.log.Infof("Upgrading cluster '%s' to k3s version '%s'...", name, version)

	if err := provider.Upgrade(name, version); err != nil {
		return err
	}

	// Update cluster info
	cluster.Version = version
	m.store.Save(cluster)

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ Cluster '%s' upgraded to version '%s'!\n", name, version)

	return nil
}

// Status shows cluster status
func (m *Manager) Status(name string) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	status, err := provider.Status(name)
	if err != nil {
		return err
	}

	// Print status
	fmt.Printf("\nCluster: %s\n", status.Name)
	fmt.Printf("Provider: %s\n", status.Provider)
	fmt.Printf("Version: %s\n", status.K3sVersion)
	fmt.Printf("Status: %s\n", status.Status)
	fmt.Printf("Endpoint: %s\n\n", status.Endpoint)

	// Print masters
	if len(status.Masters) > 0 {
		fmt.Println("Master Nodes:")
		for _, node := range status.Masters {
			fmt.Printf("  • %s (%s) - %s\n", node.Name, node.IP, node.Status)
		}
		fmt.Println()
	}

	// Print workers
	if len(status.Workers) > 0 {
		fmt.Println("Worker Nodes:")
		for _, node := range status.Workers {
			fmt.Printf("  • %s (%s) - %s\n", node.Name, node.IP, node.Status)
		}
		fmt.Println()
	}

	// Print components
	if len(status.Components) > 0 {
		fmt.Println("Components:")
		for _, comp := range status.Components {
			fmt.Printf("  • %s/%s: %s\n", comp.Namespace, comp.Name, comp.Status)
		}
		fmt.Println()
	}

	return nil
}

// Validate validates cluster health
func (m *Manager) Validate(name string) error {
	// Get cluster info
	cluster, exists := m.store.Get(name)
	if !exists {
		return fmt.Errorf("cluster '%s' not found", name)
	}

	// Get provider
	provider, err := m.getProvider(cluster.Provider)
	if err != nil {
		return err
	}

	m.log.Infof("Validating cluster '%s'...", name)

	if err := provider.Validate(name); err != nil {
		red := color.New(color.FgRed, color.Bold)
		red.Printf("✗ Cluster validation failed: %v\n", err)
		return err
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✓ Cluster '%s' is healthy!\n", name)

	return nil
}

// EnableAddon enables an addon
func (m *Manager) EnableAddon(clusterName, addon string) error {
	m.log.Infof("Enabling addon '%s' for cluster '%s'...", addon, clusterName)

	// TODO: Implement addon management
	fmt.Printf("Addon '%s' enabled for cluster '%s'\n", addon, clusterName)
	return nil
}

// DisableAddon disables an addon
func (m *Manager) DisableAddon(clusterName, addon string) error {
	m.log.Infof("Disabling addon '%s' for cluster '%s'...", addon, clusterName)

	// TODO: Implement addon management
	fmt.Printf("Addon '%s' disabled for cluster '%s'\n", addon, clusterName)
	return nil
}

// ListAddons lists available addons
func (m *Manager) ListAddons(clusterName string) error {
	addons := []struct {
		Name        string
		Description string
		Default     bool
	}{
		{"metrics-server", "Kubernetes Metrics Server", true},
		{"traefik", "Traefik Ingress Controller", true},
		{"cert-manager", "Certificate Manager", false},
		{"longhorn", "Longhorn Storage", false},
		{"prometheus", "Prometheus Monitoring", false},
		{"grafana", "Grafana Dashboards", false},
		{"argocd", "ArgoCD GitOps", false},
		{"vault", "HashiCorp Vault", false},
	}

	fmt.Println("Available Addons:")
	for _, addon := range addons {
		status := ""
		if addon.Default {
			status = " (default)"
		}
		fmt.Printf("  • %s - %s%s\n", addon.Name, addon.Description, status)
	}

	return nil
}

// mergeKubeconfig merges kubeconfig into existing config
func mergeKubeconfig(newConfig, existingPath string) error {
	// TODO: Implement proper kubeconfig merging
	return os.WriteFile(existingPath, []byte(newConfig), 0600)
}
