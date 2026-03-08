/*
Kubernetes-in-a-Box (kib) - Lightweight Kubernetes Cluster Provisioning Tool

A CLI tool for quickly provisioning lightweight Kubernetes clusters using k3s.
Supports multiple deployment scenarios including single-node, HA masters,
and multi-node worker configurations.

Features:
- Single-node k3s deployment
- High-availability master setup
- Worker node joining
- Cloud provider integration (AWS, Azure, GCP)
- On-premises bare metal support
- Automated TLS certificate management
- Integrated ingress controller
- Local storage provisioner
- Metrics server
*/

package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kubernetes-in-a-box/pkg/cluster"
	"kubernetes-in-a-box/pkg/config"
)

var (
	version = "1.0.0"
	commit  = "unknown"
	date    = "unknown"

	log = logrus.New()

	rootCmd = &cobra.Command{
		Use:   "kib",
		Short: "Kubernetes-in-a-Box - Lightweight K8s provisioning",
		Long: `
╦╔═╔═╗╦ ╦╔═╗╦ ╦╔═╗╦═╗╔═╗
╠╩╗╠═╣╠═╣║╣ ║║║║ ║╠╦╝╚═╗
╩ ╩╩ ╩╩ ╩╚═╝╚╩╝╚═╝╩╚═╚═╝

Kubernetes-in-a-Box (kib) is a CLI tool for quickly provisioning
lightweight Kubernetes clusters using k3s. It supports single-node,
HA master, and multi-node configurations on various platforms.

Examples:
  # Create a single-node cluster
  kib create --name my-cluster --provider local

  # Create an HA cluster with 3 masters
  kib create --name ha-cluster --provider aws --masters 3 --workers 3

  # List all clusters
  kib list

  # Delete a cluster
  kib delete my-cluster

  # Get cluster kubeconfig
  kib kubeconfig my-cluster
`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.kib/config.yaml)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Add commands
	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(kubeconfigCmd())
	rootCmd.AddCommand(sshCmd())
	rootCmd.AddCommand(scaleCmd())
	rootCmd.AddCommand(upgradeCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(addonsCmd())
	rootCmd.AddCommand(validateCmd())
}

func initConfig() {
	cfgFile := rootCmd.PersistentFlags().Lookup("config").Value.String()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Warn("Failed to get home directory")
		} else {
			viper.AddConfigPath(home + "/.kib")
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		}
	}

	viper.SetEnvPrefix("KIB")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Set log level
	level, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	if viper.GetBool("verbose") {
		log.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// createCmd creates the 'create' command
func createCmd() *cobra.Command {
	var opts config.ClusterOptions

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Kubernetes cluster",
		Long:  `Create a new Kubernetes cluster using k3s on the specified provider.`,
		Example: `  # Local single-node cluster
  kib create --name local-dev

  # AWS cluster with 3 masters and 5 workers
  kib create --name production --provider aws --region us-west-2 --masters 3 --workers 5

  # Azure cluster with specific VM sizes
  kib create --name azure-cluster --provider azure --master-size Standard_D4s_v3 --worker-size Standard_D8s_v3`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Name == "" {
				return fmt.Errorf("cluster name is required")
			}

			manager := cluster.NewManager(log)
			return manager.Create(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Cluster name (required)")
	cmd.Flags().StringVarP(&opts.Provider, "provider", "p", "local", "Provider (local, aws, azure, gcp, vagrant)")
	cmd.Flags().StringVar(&opts.Region, "region", "", "Cloud region (auto-detected if not specified)")
	cmd.Flags().IntVar(&opts.Masters, "masters", 1, "Number of master nodes")
	cmd.Flags().IntVar(&opts.Workers, "workers", 0, "Number of worker nodes")
	cmd.Flags().StringVar(&opts.MasterSize, "master-size", "", "Master node size/flavor")
	cmd.Flags().StringVar(&opts.WorkerSize, "worker-size", "", "Worker node size/flavor")
	cmd.Flags().StringVar(&opts.K3sVersion, "k3s-version", "latest", "k3s version to install")
	cmd.Flags().StringSliceVar(&opts.ExtraArgs, "extra-args", nil, "Extra arguments for k3s server/agent")
	cmd.Flags().BoolVar(&opts.DisableTraefik, "disable-traefik", false, "Disable Traefik ingress")
	cmd.Flags().BoolVar(&opts.DisableServiceLB, "disable-servicelb", false, "Disable ServiceLB")
	cmd.Flags().StringSliceVar(&opts.EnableAddons, "enable", []string{"metrics-server"}, "Addons to enable")
	cmd.Flags().StringVar(&opts.SSHKeyPath, "ssh-key", "", "Path to SSH private key")
	cmd.Flags().BoolVar(&opts.HA, "ha", false, "Enable HA mode (automatic with >1 master)")

	cmd.MarkFlagRequired("name")

	return cmd
}

// deleteCmd creates the 'delete' command
func deleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete [NAME]",
		Short: "Delete a Kubernetes cluster",
		Long:  `Delete a Kubernetes cluster and all associated resources.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 0 {
				name = args[0]
			}

			manager := cluster.NewManager(log)
			return manager.Delete(name, force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force deletion without confirmation")

	return cmd
}

// listCmd creates the 'list' command
func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all clusters",
		Long:    `List all Kubernetes clusters managed by kib.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.List()
		},
	}
}

// kubeconfigCmd creates the 'kubeconfig' command
func kubeconfigCmd() *cobra.Command {
	var merge bool

	cmd := &cobra.Command{
		Use:   "kubeconfig [NAME]",
		Short: "Get cluster kubeconfig",
		Long:  `Retrieve and display the kubeconfig for a cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.Kubeconfig(args[0], merge)
		},
	}

	cmd.Flags().BoolVarP(&merge, "merge", "m", false, "Merge with existing kubeconfig")

	return cmd
}

// sshCmd creates the 'ssh' command
func sshCmd() *cobra.Command {
	var node string

	cmd := &cobra.Command{
		Use:   "ssh [NAME]",
		Short: "SSH into a cluster node",
		Long:  `SSH into a node in the specified cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.SSH(args[0], node)
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Node to SSH into (defaults to first master)")

	return cmd
}

// scaleCmd creates the 'scale' command
func scaleCmd() *cobra.Command {
	var workers int

	cmd := &cobra.Command{
		Use:   "scale [NAME]",
		Short: "Scale cluster workers",
		Long:  `Scale the number of worker nodes in a cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.Scale(args[0], workers)
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 0, "Target number of workers (required)")
	cmd.MarkFlagRequired("workers")

	return cmd
}

// upgradeCmd creates the 'upgrade' command
func upgradeCmd() *cobra.Command {
	var k3sVersion string

	cmd := &cobra.Command{
		Use:   "upgrade [NAME]",
		Short: "Upgrade cluster k3s version",
		Long:  `Upgrade the k3s version on all nodes in the cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.Upgrade(args[0], k3sVersion)
		},
	}

	cmd.Flags().StringVar(&k3sVersion, "k3s-version", "latest", "Target k3s version")

	return cmd
}

// statusCmd creates the 'status' command
func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [NAME]",
		Short: "Show cluster status",
		Long:  `Display detailed status information about a cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.Status(args[0])
		},
	}
}

// addonsCmd creates the 'addons' command
func addonsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addons",
		Short: "Manage cluster addons",
		Long:  `Enable or disable addons for a cluster.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "enable [CLUSTER] [ADDON]",
		Short: "Enable an addon",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.EnableAddon(args[0], args[1])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "disable [CLUSTER] [ADDON]",
		Short: "Disable an addon",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.DisableAddon(args[0], args[1])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list [CLUSTER]",
		Short: "List available addons",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			clusterName := ""
			if len(args) > 0 {
				clusterName = args[0]
			}
			return manager.ListAddons(clusterName)
		},
	})

	return cmd
}

// validateCmd creates the 'validate' command
func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [NAME]",
		Short: "Validate cluster health",
		Long:  `Run health checks and validation on a cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := cluster.NewManager(log)
			return manager.Validate(args[0])
		},
	}
}
