package tesla

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	AccessToken  string `yaml:"accessToken"`
	RefreshToken string `yaml:"refreshToken"`
	MainVehicle  string `yaml:"mainVehicle"`
	filePath     string
	log          *zap.SugaredLogger
}

func NewConfig(filePath string, log *zap.SugaredLogger) *Config {
	return &Config{
		filePath: filePath,
		log:      log,
	}
}

func (c *Config) Load() *Config {
	cfgYaml, err := os.ReadFile(c.filePath)

	if err != nil {
		c.log.Fatal(err)
	}

	err = yaml.Unmarshal(cfgYaml, &c)

	if err != nil {
		c.log.Fatal(err)
	}

	return c
}

func (c *Config) Save() *Config {
	cfgYaml, err := yaml.Marshal(&c)

	if err != nil {
		c.log.Error(err)
		return nil
	}

	err = os.WriteFile(c.filePath, cfgYaml, 0)

	if err != nil {
		c.log.Error(err)
		return nil
	}

	return c
}
