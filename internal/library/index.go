package library

type Entry struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Status string `yaml:"status"`
}

type Index struct {
	Version     int     `yaml:"version"`
	LastUpdated string  `yaml:"last_updated"`
	Games       []Entry `yaml:"games"`
}
