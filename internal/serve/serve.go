package serve

import (
	"fmt"

	"github.com/cfg8er/cfg8er/internal/config"
	"github.com/cfg8er/cfg8er/pkg/repository"
	"github.com/gin-gonic/gin"
	"gopkg.in/urfave/cli.v1"
)

var repoLookup map[string]config.Repo

func Run(c *cli.Context) error {
	if !c.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	configPath := c.String("config")

	var err error
	repoLookup, err = config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// Clone all the repos
	if err := cloneRepos(); err != nil {
		return err
	}

	router := newRouter()

	return router.Run(c.String("listen"))
}

func cloneRepos() error {
	for n, r := range repoLookup {
		if r.ClonedRepo == (repository.Repository{}) {
			fmt.Printf("Cloning repo %s, %s\n", n, r.URL)
			clonedRepo, err := repository.CloneBare(r.URL)
			if err != nil {
				return err
			}

			r.ClonedRepo = clonedRepo
			repoLookup[n] = r
		}
	}

	return nil
}
