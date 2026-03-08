module vmware-migration

go 1.21

require (
	github.com/spf13/cobra v1.8.0
	github.com/sirupsen/logrus v1.9.3
	github.com/vmware/govmomi v0.33.1
	k8s.io/client-go v0.28.4
	k8s.io/api v0.28.4
	k8s.io/apimachinery v0.28.4
	kubevirt.io/api v1.1.0
	kubevirt.io/containerized-data-importer-api v1.57.0
	github.com/fatih/color v1.16.0
	github.com/olekukonko/tablewriter v0.0.5
	gopkg.in/yaml.v3 v3.0.1
)
