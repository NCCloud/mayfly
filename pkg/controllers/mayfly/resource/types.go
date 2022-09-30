package resource

type Resource struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
}

type ResourceList struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Items      []Resource
}

type Resources struct {
	Resources []Resource `yaml:"resources"`
}
