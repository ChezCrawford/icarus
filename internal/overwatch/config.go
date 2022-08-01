package overwatch

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	EnphaseSystemId    string `split_words:"true" required:"true"`
	EnphaseApiKey      string `split_words:"true" required:"true"`
	EnphaseAccessToken string `split_words:"true" required:"true"`
	PdApiKey           string `split_words:"true" required:"true"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
