package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewFastGoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "fg-apiserver",
		Short:        "A very lightweight full go project",
		Long:         "A very lightweight full go project, designed to help beginners guickly learn Go project development.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello FastGO!")
			return nil
		},
		Args: cobra.NoArgs,
	}
	return cmd
}
