package cmd

import (
	"github.com/enclaive/relay/cmd/serve"
	"github.com/enclaive/relay/cmd/setup"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use: "relay",
	}
}

func Execute() error {
	cmd := newRootCmd()
	serve.RegisterCommands(cmd)
	setup.RegisterCommands(cmd)

	return cmd.Execute()
}
