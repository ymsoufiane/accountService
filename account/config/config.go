package config

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `default:"localhost"`
		Port     string `default:"3306"`
		Username string
		Password string
		Dbname   string
		Driver   string `default:"mysql"`
	}

	File struct {
		Logpath string `default:"Logs"`
	}

	Jwt struct {
		Secret string `default:"secretTest"`
	}

	Server struct {
		Port string `default:"8080"`
	}

}

func Load(filename string) *Config {

	file, er := os.Open(filename)
	if er != nil {
		log.Fatal(er)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, fileinfo.Size())

	count, err := file.Read(data)

	file.Close()

	config := &Config{}

	if count > 0 {

		if err := yaml.Unmarshal(data, config); err != nil {
			log.Fatal(err)
		}

	}

	return config
}
