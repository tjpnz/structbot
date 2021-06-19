package config

import "github.com/imdario/mergo"

// Config defines fields for configuring factories upon creation. Currently this
// is limited to transformers for patching. Currently it should be considered
// internal but it may take on a larger role in future iterations to allow for
// user configuration.
type Config struct {
	Transformers mergo.Transformers
}
