package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// NewConfig returns *Root config
func NewConfig(fileName string) (*Root, error) {
	conf, err := loadConfig(fileName)
	if err != nil {
		return nil, err
	}

	return conf, err
}

// GetEnvConfPath returns toml file path from environment variable `$GO_GOOGLE_HOME_CONF`
func GetEnvConfPath() string {
	path := os.Getenv("GO_GOOGLE_HOME_CONF")
	if strings.Contains(path, "${GOPATH}") {
		gopath := os.Getenv("GOPATH")
		path = strings.Replace(path, "${GOPATH}", gopath, 1)
	}
	return path
}

// load config file
func loadConfig(fileName string) (*Root, error) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read file: %s", fileName)
	}

	var root Root
	_, err = toml.Decode(string(d), &root)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse: %s", fileName)
	}

	// check validation of config
	if err = root.validate(); err != nil {
		return nil, err
	}

	return &root, nil
}

func (r *Root) validate() error {
	return validator.New().Struct(r)
}
