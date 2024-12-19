package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RyaiConfig struct {
	Llm LlmSet `yaml:"Llm"`
	Log LogSet `yaml:"Log"`
}

type LogSet struct {
	Level string `yaml:"level"`
	Size  string `yaml:"size"`
	Lfile string `yaml:"lfile"`
	Num   string `yaml:"num"`
	Age   string `yaml:"age"`
}

type LlmSet struct {
	Name string `yaml:"name"`
	Mod  string `yaml:"mod"`
}

func (c LogSet) String() string {
	return fmt.Sprintf("level:%s,\nsize:%s,\nlfile:%s,\nnum:%s,\nage:%s;\n\n", c.Level, c.Size, c.Lfile, c.Num, c.Age)
}

func (c LlmSet) String() string {
	return fmt.Sprintf("name:%s,\nmod:%s;\n\n", c.Name, c.Mod)
}

func (c RyaiConfig) String() string {
	return fmt.Sprintf("Llm:\n%sLog:\n%s\n", c.Llm.String(), c.Log.String())
}

func ReadCfg() (cfg RyaiConfig, err error) {
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
