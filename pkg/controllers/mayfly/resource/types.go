package resource

type Resource struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
}

type Resources struct {
	Resources []Resource `yaml:"resources"`
}
