package config

type Tab struct {
	Name    string `toml:"name"`
	Command string `toml:"command"`
}

type Config struct {
	Name        string   `toml:"name"`
	Path        string   `toml:"path"`
	DefaultMode string   `toml:"default_mode"`
	Backend     string   `toml:"backend"`
	PreOpen     []string `toml:"pre_open"`
	PostOpen    []string `toml:"post_open"`
	Tabs        []Tab    `toml:"tabs"`
}
