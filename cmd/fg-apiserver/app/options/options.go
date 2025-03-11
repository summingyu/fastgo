package options

import (
	"github.com/onexstack/fastgo/internal/apiserver"
	genericoptions "github.com/onexstack/fastgo/pkg/options"
)

type ServerOptions struct {
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: genericoptions.NewMySQLOptions(),
	}
}

func (o *ServerOptions) Validate() error {
	if err := o.MySQLOptions.Validate(); err != nil {
		return err
	}
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		MySQLOptions: o.MySQLOptions,
	}, nil
}
