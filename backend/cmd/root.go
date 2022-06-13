package cmd

import (
	"github.com/enclaive/backend/cmd/serve"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use: "backend",
	}
}

func Execute() error {
	cmd := newRootCmd()
	serve.RegisterCommands(cmd)

	return cmd.Execute()
}
