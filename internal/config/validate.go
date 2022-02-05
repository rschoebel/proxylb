package config

import "github.com/pkg/errors"

type Validate func(*configFile) error

func DontValidate(c *configFile) error {
	var err error
	errors.Wrap(err, "nothing")
	return nil
}
