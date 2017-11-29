package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/swarmfund/api"
	"gitlab.com/swarmfund/api/assets"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/log"
)

var (
	configFile     string
	configInstance config.Config
	rootCmd        = &cobra.Command{
		Use: "api",
	}
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := api.NewApp(configInstance)
			if err != nil {
				log.WithField("service", "init").WithError(err).Fatal("failed to init app")
			}
			app.Serve()
		},
	}
	migrateCmd = &cobra.Command{
		Use:   "migrate [up|down|redo] [COUNT]",
		Short: "migrate schema",
		Long:  "performs a schema migration command",
		Run: func(cmd *cobra.Command, args []string) {
			migrateDB(cmd, args, configInstance.API().DatabaseURL, assets.Migrations.Migrate)
		},
	}
)

func main() {
	cobra.OnInitialize(func() {
		c, err := initConfig(configFile)
		if err != nil {
			log.WithField("service", "init").WithError(err).Fatal("failed to init config")
		}
		configInstance = c
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "config file")
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.Execute()
}

func initConfig(fn string) (config.Config, error) {
	c := config.NewViperConfig(fn)
	if err := c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}
