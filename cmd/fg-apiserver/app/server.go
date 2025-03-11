package app

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/onexstack/fastgo/cmd/fg-apiserver/app/options"
)

var configFile string

func NewFastGoCommand() *cobra.Command {
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		Use:          "fg-apiserver",
		Short:        "A very lightweight full go project",
		Long:         "A very lightweight full go project, designed to help beginners guickly learn Go project development.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}

			if err := opts.Validate(); err != nil {
				return err
			}
			fmt.Printf("Read MySQL host from Viper: %s\n", viper.GetString("mysql.addr"))
			fmt.Printf("Read MySQL username from opts: %s\n", opts.MySQLOptions.Username)
			jsonData, _ := json.MarshalIndent(opts, "", " ")
			fmt.Println(string(jsonData))
			return nil
		},
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(onInitialize)
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the fg-apiserver configuration file.")
	return cmd
}
