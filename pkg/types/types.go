package types

// PodSelectors represents what selectors to use when listing pods
type PodSelectors struct {
	Label string
	Field string
	Names []string
}

// EventSelectors represents what selectors to use when listing pods
type EventSelectors struct {
	Label string
	Field string
}

// ClientConfig represents configuration for the Kubernetes Client
type ClientConfig struct {
	ConfigFile    string
	Namespace     string
	Context       string
	AllNamespaces bool
}

// PodCondition is a wrapper around a Kubernetes Pod Condition
type PodCondition struct {
	Type       string
	Successful bool
	Reason     string
	Message    string
}

// ContainerStatus is a wrapper around a Kubernetes Pod's Container Status
type ContainerStatus struct {
	Name  string
	Ready bool
}
