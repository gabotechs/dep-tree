package check

type Config struct {
	Path                      string
	Entrypoints               []string                    `yaml:"entrypoints"`
	AllowCircularDependencies bool                        `yaml:"allowCircularDependencies"`
	Aliases                   map[string][]string         `yaml:"aliases"`
	WhiteList                 map[string]WhiteListEntries `yaml:"allow"`
	BlackList                 map[string][]BlackListEntry `yaml:"deny"`
}

func (c *Config) Init(path string) {
	c.Path = path
	c.expandAliases()
}

func (c *Config) expandAliases() {
	for k, entries := range c.WhiteList {
		newV := make([]string, 0)
		for _, entry := range entries.To {
			if aliases, ok := c.Aliases[entry]; ok {
				newV = append(newV, aliases...)
			} else {
				newV = append(newV, entry)
			}
		}
		c.WhiteList[k] = WhiteListEntries{
			To:     newV,
			Reason: entries.Reason,
		}
	}

	for k, entries := range c.BlackList {
		newV := make([]BlackListEntry, 0)
		for _, entry := range entries {
			if aliases, ok := c.Aliases[entry.To]; ok {
				for _, alias := range aliases {
					newV = append(newV, BlackListEntry{
						To:     alias,
						Reason: entry.Reason,
					})
				}
			} else {
				newV = append(newV, entry)
			}
		}
		c.BlackList[k] = newV
	}
}

type BlackListEntry struct {
	To     string `yaml:"to"`
	Reason string `yaml:"reason"`
}

func (v *BlackListEntry) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err == nil {
		v.To = str
		return nil
	}
	temp := struct {
		To     string `yaml:"to"`
		Reason string `yaml:"reason"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}
	v.To = temp.To
	v.Reason = temp.Reason
	return nil
}

type WhiteListEntries struct {
	To     []string `yaml:"to"`
	Reason string   `yaml:"reason"`
}

func (v *WhiteListEntries) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strList []string
	if err := unmarshal(&strList); err == nil {
		v.To = strList
		return nil
	}

	temp := struct {
		To     []string `yaml:"to"`
		Reason string   `yaml:"reason"`
	}{}
	err := unmarshal(&temp)
	if err != nil {
		return err
	}
	v.To = temp.To
	v.Reason = temp.Reason
	return nil
}
