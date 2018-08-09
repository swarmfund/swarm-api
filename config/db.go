package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2"
)

func (c *ViperConfig) DB() *db2.Repo {
	c.Lock()
	defer c.Unlock()

	if c.db != nil {
		return c.db
	}

	config := struct {
		URL     string `fig:"url,required"`
		MaxIdle int    `fig:"max_idle"`
		MaxOpen int    `fig:"max_open"`
	}{
		MaxIdle: 4,
		MaxOpen: 12,
	}

	if err := figure.Out(&config).From(c.GetStringMap("db")).Please(); err != nil {
		panic(errors.Wrap(err, "failed to db"))
	}

	repo, err := db2.Open(config.URL)
	if err != nil {
		panic(errors.Wrap(err, "failed to open database"))
	}

	repo.DB.SetMaxIdleConns(config.MaxIdle)
	repo.DB.SetMaxOpenConns(config.MaxOpen)

	c.db = repo
	return c.db
}
