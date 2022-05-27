package setup

import (
	"fmt"
	"github.com/enclaive/relay/config"
	"github.com/enclaive/relay/persistence"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func newSetupCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Setup the database",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var cfg config.PostgresConfig
			err := config.FromEnv(&cfg)
			if err != nil {
				return err
			}

			err = validator.New().Struct(cfg)
			if err != nil {
				return fmt.Errorf("validation of config failed: %w", err)
			}

			db := persistence.New(cfg)
			return db.ApplySchema()
		},
	}
}

func RegisterCommands(parent *cobra.Command) {
	cmd := newSetupCommand()
	parent.AddCommand(cmd)
}
