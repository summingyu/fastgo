package options

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/onexstack/fastgo/internal/apiserver"
	genericoptions "github.com/onexstack/fastgo/pkg/options"
)

type ServerOptions struct {
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	Addr         string                       `json:"addr" mapstructure:"addr"`
	JWTKey       string                       `json:"jwt-key" mapstructure:"jwt-key"`
	Expiration   time.Duration                `json:"expiration" mapstructure:"expiration"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MySQLOptions: genericoptions.NewMySQLOptions(),
		Addr:         "0.0.0.0:6666",
		Expiration:   2 * time.Hour,
	}
}

func (o *ServerOptions) Validate() error {
	if err := o.MySQLOptions.Validate(); err != nil {
		return err
	}
	if o.Addr == "" {
		return fmt.Errorf("server address cannot be empty")
	}
	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invalid server address format '%s': %w", o.Addr, err)
	}
	port, err := strconv.Atoi(portStr)

	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid server port '%s'", portStr)
	}
	if len(o.JWTKey) < 6 {
		return fmt.Errorf("jwt key length must be at least 6 characters")
	}
	return nil
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		MySQLOptions: o.MySQLOptions,
		Addr:         o.Addr,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
	}, nil
}
