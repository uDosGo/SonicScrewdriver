package types

type GameManifest struct {
	Name        string      `yaml:"name" json:"name"`
	Version     string      `yaml:"version" json:"version"`
	Title       string      `yaml:"title" json:"title"`
	Description string      `yaml:"description" json:"description"`
	Type        string      `yaml:"type" json:"type"`
	Runtime     Runtime     `yaml:"runtime" json:"runtime"`
	Assets      []Asset     `yaml:"assets" json:"assets"`
	Integration Integration `yaml:"integration" json:"integration"`
}

type Runtime struct {
	Engine    string `yaml:"engine" json:"engine"`
	Container string `yaml:"container" json:"container"`
	Port      int    `yaml:"port" json:"port"`
	WebUI     bool   `yaml:"web_ui" json:"web_ui"`
}

type Asset struct {
	Name       string `yaml:"name" json:"name"`
	Required   bool   `yaml:"required" json:"required"`
	Source     string `yaml:"source" json:"source"`
	Location   string `yaml:"location" json:"location"`
	Validation string `yaml:"validation" json:"validation"`
}

type Integration struct {
	Skin DefaultSkin `yaml:"skin" json:"skin"`
	Lens LensConfig  `yaml:"lens" json:"lens"`
}

type DefaultSkin struct {
	Default   string   `yaml:"default" json:"default"`
	Available []string `yaml:"available" json:"available"`
}

type LensConfig struct {
	Variables []LensVariable `yaml:"variables" json:"variables"`
}

type LensVariable struct {
	Name   string `yaml:"name" json:"name"`
	Source string `yaml:"source" json:"source"`
	Offset string `yaml:"offset,omitempty" json:"offset,omitempty"`
}
