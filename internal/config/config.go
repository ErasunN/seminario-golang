package config
import "io/ioutil"
import "gopkg.in/yaml.v3"

type DbConfig struct{
	Driver string `yaml:"driver"`
}

type Config struct{
	DB DbConfig `yaml:"db"`
	Version string `yaml:"vers"`
}

func LoadConfig(filename string) (*Config, error){
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}


}