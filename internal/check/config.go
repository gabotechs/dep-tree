package check

type Config struct {
	Path                      string
	Entrypoints               []string            `yaml:"entrypoints"`
	AllowCircularDependencies bool                `yaml:"allowCircularDependencies"`
	Aliases                   map[string][]string `yaml:"aliases"`
	WhiteList                 map[string][]string `yaml:"allow"`
	BlackList                 map[string][]string `yaml:"deny"`
}

func (c *Config) Init(path string) {
	c.Path = path
	c.expandAliases()
}

func (c *Config) expandAliases() {
	lists := []map[string][]string{
		c.WhiteList,
		c.BlackList,
	}
	for _, list := range lists {
		for k, v := range list {
			newV := make([]string, 0)
			for _, entry := range v {
				if alias, ok := c.Aliases[entry]; ok {
					newV = append(newV, alias...)
				} else {
					newV = append(newV, entry)
				}
			}
			list[k] = newV
		}
	}
}
