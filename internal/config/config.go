package config

// Tab represents a single tab in a project workspace.
type Tab struct {
	Name    string `toml:"name"`
	Command string `toml:"command"`
}

// Config represents a project's TOML configuration.
type Config struct {
	Name        string   `toml:"name"`
	Path        string   `toml:"path"`
	DefaultMode string   `toml:"default_mode"`
	Backend     string   `toml:"backend"`
	PreOpen     []string `toml:"pre_open"`
	PostOpen    []string `toml:"post_open"`
	Tabs        []Tab    `toml:"tabs"`
}
