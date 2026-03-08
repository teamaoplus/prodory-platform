/*
VMware to KubeVirt Migration Tool

A CLI tool for migrating virtual machines from VMware vSphere
to Kubernetes using KubeVirt.

Features:
- VM inventory discovery from vCenter
- VM analysis and compatibility checking
- Disk image conversion (VMDK to RAW/QCOW2)
- Hot and cold migration support
- Network mapping
- Progress tracking
- Rollback support
*/

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"vmware-migration/pkg/analyze"
	"vmware-migration/pkg/migrate"
)

var (
	version = "1.0.0"
	log     = logrus.New()

	// vCenter connection flags
	vcenterURL      string
	vcenterUser     string
	vcenterPassword string
	vcenterInsecure bool

	// Kubernetes flags
	kubeconfig string
	namespace  string

	// Migration flags
	sourceVM      string
	targetName    string
	storageClass  string
	networkMapping []string
	coldMigration bool

	rootCmd = &cobra.Command{
		Use:   "vmware-migrate",
		Short: "VMware to KubeVirt Migration Tool",
		Long: `
╔╗ ╦  ╦╔╦╗╔═╗╔═╗╔═╗  ╔╦╗╔═╗╔╦╗╔═╗╦  ╦╔╦╗╔═╗
╠╩╗║  ║║║║║╣ ║╣ ╚═╗   ║║╠═╣ ║ ╠═╣║  ║ ║ ╚═╗
╚═╝╩═╝╩╩ ╩╚═╝╚═╝╚═╝  ═╩╝╩ ╩ ╩ ╩ ╩╩═╝╩ ╩ ╚═╝

Migrate virtual machines from VMware vSphere to Kubernetes KubeVirt.

Examples:
  # Discover VMs in vCenter
  vmware-migrate discover --vcenter-url https://vcenter.local --vcenter-user admin

  # Analyze a VM for migration compatibility
  vmware-migrate analyze --vcenter-url https://vcenter.local --source-vm my-vm

  # Migrate a VM to KubeVirt
  vmware-migrate migrate --vcenter-url https://vcenter.local \
    --source-vm my-vm --target-name my-vm-kubevirt \
    --namespace default --storage-class standard
`,
		Version: version,
	}
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Global flags
	rootCmd.PersistentFlags().StringVar(&vcenterURL, "vcenter-url", "", "vCenter URL (required)")
	rootCmd.PersistentFlags().StringVar(&vcenterUser, "vcenter-user", "", "vCenter username")
	rootCmd.PersistentFlags().StringVar(&vcenterPassword, "vcenter-password", "", "vCenter password")
	rootCmd.PersistentFlags().BoolVar(&vcenterInsecure, "vcenter-insecure", false, "Skip TLS verification")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")

	// Add commands
	rootCmd.AddCommand(discoverCmd())
	rootCmd.AddCommand(analyzeCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(rollbackCmd())
}

func discoverCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "discover",
		Short: "Discover VMs in vCenter",
		Long:  `List all virtual machines available in the vCenter inventory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if vcenterURL == "" {
				return fmt.Errorf("--vcenter-url is required")
			}

			// Get credentials
			user, password, err := getCredentials()
			if err != nil {
				return err
			}

			// Create analyzer
			analyzer, err := analyze.NewVCenterAnalyzer(vcenterURL, user, password, vcenterInsecure, log)
			if err != nil {
				return err
			}
			defer analyzer.Close()

			// Discover VMs
			vms, err := analyzer.DiscoverVMs(context.Background())
			if err != nil {
				return fmt.Errorf("failed to discover VMs: %w", err)
			}

			// Display results
			fmt.Printf("\nFound %d virtual machines:\n\n", len(vms))
			for _, vm := range vms {
				powerState := "⏻"
				if vm.PowerState == "poweredOn" {
					powerState = "🟢"
				} else if vm.PowerState == "poweredOff" {
					powerState = "🔴"
				}
				fmt.Printf("  %s %s (%s, %d vCPU, %d MB RAM, %d disks)\n",
					powerState, vm.Name, vm.GuestOS, vm.CPU, vm.MemoryMB, len(vm.Disks))
			}
			fmt.Println()

			return nil
		},
	}
}

func analyzeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "analyze",
		Short: "Analyze VM for migration compatibility",
		Long:  `Analyze a VMware VM and report migration compatibility and requirements.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if vcenterURL == "" {
				return fmt.Errorf("--vcenter-url is required")
			}
			if sourceVM == "" {
				return fmt.Errorf("--source-vm is required")
			}

			// Get credentials
			user, password, err := getCredentials()
			if err != nil {
				return err
			}

			// Create analyzer
			analyzer, err := analyze.NewVCenterAnalyzer(vcenterURL, user, password, vcenterInsecure, log)
			if err != nil {
				return err
			}
			defer analyzer.Close()

			// Analyze VM
			report, err := analyzer.AnalyzeVM(context.Background(), sourceVM)
			if err != nil {
				return fmt.Errorf("failed to analyze VM: %w", err)
			}

			// Display report
			displayAnalysisReport(report)

			return nil
		},
	}
}

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate VM to KubeVirt",
		Long:  `Migrate a VMware VM to Kubernetes as a KubeVirt VM.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if vcenterURL == "" {
				return fmt.Errorf("--vcenter-url is required")
			}
			if sourceVM == "" {
				return fmt.Errorf("--source-vm is required")
			}
			if targetName == "" {
				targetName = sourceVM
			}

			// Get credentials
			user, password, err := getCredentials()
			if err != nil {
				return err
			}

			// Create migrator
			migrator, err := migrate.NewMigrator(
				vcenterURL, user, password, vcenterInsecure,
				kubeconfig, namespace, storageClass,
				log,
			)
			if err != nil {
				return err
			}
			defer migrator.Close()

			// Create migration options
			opts := &migrate.MigrationOptions{
				SourceVM:       sourceVM,
				TargetName:     targetName,
				Namespace:      namespace,
				StorageClass:   storageClass,
				NetworkMapping: parseNetworkMapping(networkMapping),
				ColdMigration:  coldMigration,
			}

			// Run migration
			green := color.New(color.FgGreen, color.Bold)
			green.Printf("\nStarting migration of '%s' to KubeVirt...\n\n", sourceVM)

			result, err := migrator.Migrate(context.Background(), opts)
			if err != nil {
				return fmt.Errorf("migration failed: %w", err)
			}

			// Display result
			green.Printf("\n✓ Migration completed successfully!\n\n")
			fmt.Printf("KubeVirt VM: %s/%s\n", result.Namespace, result.VMName)
			fmt.Printf("Status: %s\n", result.Status)
			fmt.Printf("Duration: %s\n", result.Duration.Round(time.Second))
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVarP(&sourceVM, "source-vm", "s", "", "Source VM name (required)")
	cmd.Flags().StringVarP(&targetName, "target-name", "t", "", "Target VM name (defaults to source name)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Target Kubernetes namespace")
	cmd.Flags().StringVar(&storageClass, "storage-class", "standard", "Storage class for PVCs")
	cmd.Flags().StringSliceVar(&networkMapping, "network-map", nil, "Network mapping (format: source=target)")
	cmd.Flags().BoolVar(&coldMigration, "cold", false, "Perform cold migration (VM must be powered off)")

	cmd.MarkFlagRequired("source-vm")

	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List migration jobs",
		Long:  `List all migration jobs and their status.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			migrator, err := migrate.NewMigrator(
				"", "", "", false,
				kubeconfig, namespace, "",
				log,
			)
			if err != nil {
				return err
			}
			defer migrator.Close()

			jobs, err := migrator.ListMigrations(context.Background(), namespace)
			if err != nil {
				return err
			}

			if len(jobs) == 0 {
				fmt.Println("No migration jobs found.")
				return nil
			}

			fmt.Printf("\nMigration Jobs:\n\n")
			for _, job := range jobs {
				status := "⏳"
				if job.Status == "completed" {
					status = "✓"
				} else if job.Status == "failed" {
					status = "✗"
				}
				fmt.Printf("  %s %s: %s -> %s/%s (%s)\n",
					status, job.ID, job.SourceVM, job.TargetNamespace, job.TargetName, job.Status)
			}
			fmt.Println()

			return nil
		},
	}
}

func rollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback [MIGRATION-ID]",
		Short: "Rollback a migration",
		Long:  `Rollback a failed or completed migration.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			migrationID := args[0]

			migrator, err := migrate.NewMigrator(
				"", "", "", false,
				kubeconfig, namespace, "",
				log,
			)
			if err != nil {
				return err
			}
			defer migrator.Close()

			if err := migrator.Rollback(context.Background(), migrationID); err != nil {
				return err
			}

			green := color.New(color.FgGreen, color.Bold)
			green.Printf("✓ Migration '%s' rolled back successfully!\n", migrationID)

			return nil
		},
	}
}

func getCredentials() (string, string, error) {
	user := vcenterUser
	password := vcenterPassword

	if user == "" {
		fmt.Print("vCenter Username: ")
		fmt.Scanln(&user)
	}

	if password == "" {
		fmt.Print("vCenter Password: ")
		// In production, use proper password masking
		fmt.Scanln(&password)
	}

	if user == "" || password == "" {
		return "", "", fmt.Errorf("vCenter credentials required")
	}

	return user, password, nil
}

func displayAnalysisReport(report *analyze.AnalysisReport) {
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	fmt.Printf("\n=== VM Analysis Report: %s ===\n\n", report.VMName)

	// Basic info
	fmt.Println("Configuration:")
	fmt.Printf("  vCPU: %d\n", report.CPU)
	fmt.Printf("  Memory: %d MB\n", report.MemoryMB)
	fmt.Printf("  Guest OS: %s\n", report.GuestOS)
	fmt.Printf("  Power State: %s\n", report.PowerState)
	fmt.Println()

	// Disks
	fmt.Println("Disks:")
	for _, disk := range report.Disks {
		fmt.Printf("  • %s: %d GB (%s)\n", disk.Name, disk.SizeGB, disk.Type)
	}
	fmt.Println()

	// Networks
	fmt.Println("Networks:")
	for _, net := range report.Networks {
		fmt.Printf("  • %s (%s)\n", net.Name, net.Type)
	}
	fmt.Println()

	// Compatibility
	fmt.Println("KubeVirt Compatibility:")
	if report.Compatible {
		green.Println("  ✓ VM is compatible with KubeVirt")
	} else {
		red.Println("  ✗ VM has compatibility issues")
	}
	fmt.Println()

	// Issues
	if len(report.Issues) > 0 {
		fmt.Println("Issues:")
		for _, issue := range report.Issues {
			switch issue.Severity {
			case "error":
				red.Printf("  ✗ %s\n", issue.Message)
			case "warning":
				yellow.Printf("  ⚠ %s\n", issue.Message)
			default:
				fmt.Printf("  • %s\n", issue.Message)
			}
		}
		fmt.Println()
	}

	// Recommendations
	if len(report.Recommendations) > 0 {
		fmt.Println("Recommendations:")
		for _, rec := range report.Recommendations {
			fmt.Printf("  • %s\n", rec)
		}
		fmt.Println()
	}

	// Migration estimate
	fmt.Printf("Estimated Migration Time: %s\n", report.EstimatedTime)
	fmt.Printf("Estimated Data Transfer: %d GB\n", report.TotalDiskSizeGB)
	fmt.Println()
}

func parseNetworkMapping(mappings []string) map[string]string {
	result := make(map[string]string)
	for _, m := range mappings {
		parts := splitMapping(m)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func splitMapping(s string) []string {
	for i, c := range s {
		if c == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
