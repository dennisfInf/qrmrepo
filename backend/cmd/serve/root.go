package serve

import (
	"fmt"
	"github.com/enclaive/backend/config"
	"github.com/enclaive/backend/server"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the backend",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)

			var cfg config.Config
			err := config.FromEnv(&cfg)
			if err != nil {
				return err
			}

			err = validator.New().Struct(cfg)
			if err != nil {
				return fmt.Errorf("validation of config failed: %w", err)
			}

			s := server.New(cfg)
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
