package main

import (
	"os"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gitlab.com/swarmfund/api"
	"gitlab.com/swarmfund/api/assets"
	"gitlab.com/swarmfund/api/config"
)

var (
	// use exitCode to set process exit code instead of direct os.Exit call
	exitCode       int
	configFile     string
	configInstance config.Config
	rootCmd        = &cobra.Command{
		Use: "api",
	}
	runCmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if rvr := recover(); rvr != nil {
					configInstance.Log().WithRecover(rvr).Error("app panicked")
				}
			}()
			app, err := api.NewApp(configInstance)
			if err != nil {
				configInstance.Log().WithField("service", "init").WithError(err).Fatal("failed to init app")
			}
			app.Serve()
		},
	}
	migrateCmd = &cobra.Command{
		Use:   "migrate [up|down|redo] [COUNT]",
		Short: "migrate schema",
		Long:  "performs a schema migration command",
		Run: func(cmd *cobra.Command, args []string) {
			log := configInstance.Log().WithField("service", "migration")
			var count int
			// Allow invocations with 1 or 2 args.  All other args counts are erroneous.
			if len(args) < 1 || len(args) > 2 {
				log.WithField("arguments", args).Error("wrong argument count")
				exitCode = 1
				return
			}
			// If a second arg is present, parse it to an int and use it as the count
			// argument to the migration call.
			if len(args) == 2 {
				var err error
				if count, err = cast.ToIntE(args[1]); err != nil {
					log.WithError(err).Error("failed to parse count")
					exitCode = 1
					return
				}
			}

			applied, err := migrateDB(args[0], count, configInstance.DB().DB.DB, assets.Migrations.Migrate)
			log = log.WithField("applied", applied)
			if err != nil {
				log.WithError(err).Error("migration failed")
				exitCode = 1
				return
			}
			log.Info("migrations applied")
		},
	}
)

func main() {
	defer os.Exit(exitCode)
	cobra.OnInitialize(func() {
		c, err := initConfig(configFile)
		if err != nil {
			panic(err)
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
