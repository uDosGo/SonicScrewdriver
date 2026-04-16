package types

type ContainerSpec struct {
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Ports     map[string]string `json:"ports"`
	Env       map[string]string `json:"env"`
	Volumes   map[string]string `json:"volumes"`
	AutoStart bool              `json:"auto_start"`
}
