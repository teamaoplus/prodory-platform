package config

// ClusterOptions represents options for creating a cluster
type ClusterOptions struct {
	Name             string
	Provider         string
	Region           string
	Masters          int
	Workers          int
	MasterSize       string
	WorkerSize       string
	K3sVersion       string
	ExtraArgs        []string
	DisableTraefik   bool
	DisableServiceLB bool
	EnableAddons     []string
	SSHKeyPath       string
	HA               bool
}
