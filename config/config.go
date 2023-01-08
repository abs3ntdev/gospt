package config

import (
	"path/filepath"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

var Values struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	DeviceId     string `yaml:"device_id"`
	Port         int    `yaml:"port"`
	DeviceName   string `yaml:"device_name"`
}

func LoadConfig(configDir string) {
	yamlDecoder := aconfigyaml.New()

	loader := aconfig.LoaderFor(&Values, aconfig.Config{
		AllowUnknownFields: true,
		AllowUnknownEnvs:   true,
		AllowUnknownFlags:  true,
		SkipFlags:          true,
		DontGenerateTags:   true,
		MergeFiles:         true,
		EnvPrefix:          "",
		FlagPrefix:         "",
		Files: []string{
			filepath.Join(configDir, "client.yml"),
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": yamlDecoder,
		},
	})
	if err := loader.Load(); err != nil {
		panic(err)
	}
}
