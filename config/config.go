package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

)

var (
	Options AppOptions
)

func InitOption(cfgPath string) error {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &Options); err != nil {
		return err
	}
	Options.FillWithDefaults()
	return nil
}

