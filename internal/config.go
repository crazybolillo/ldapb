package internal

type Config struct {
	Backends map[string]Backend `yaml:"backends"`
}

type Backend struct {
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
}
