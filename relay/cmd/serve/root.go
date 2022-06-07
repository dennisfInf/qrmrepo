package serve

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/enclaive/relay/config"
	"github.com/enclaive/relay/persistence"
	"github.com/enclaive/relay/server"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the relay",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			var cfg config.GlobalConfig
			err := config.FromEnv(&cfg)
			if err != nil {
				return err
			}

			err = validator.New().Struct(cfg)
			if err != nil {
				return fmt.Errorf("validation of config failed: %w", err)
			}

			db := persistence.New(cfg.Postgres)

			s, err := server.New(cfg.Server, db)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}
			go func() {
				err := s.Run()
				if err != nil {
					panic(err)
				}
			}()

			<-stop

			return s.Stop()
		},
	}
}

func RegisterCommands(parent *cobra.Command) {
	cmd := newServeCommand()
	parent.AddCommand(cmd)
}
