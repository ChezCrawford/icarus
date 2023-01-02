package overwatch

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	EnphaseSystemId      string        `split_words:"true" required:"true"`
	EnphaseApiKey        string        `split_words:"true" required:"true"`
	EnphaseAccessToken   string        `split_words:"true" required:"true"`
	EnphasePollFrequency time.Duration `split_words:"true" default:"15s"`
	PdApiKey             string        `split_words:"true" required:"true"`
}

func LoadConfig() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
