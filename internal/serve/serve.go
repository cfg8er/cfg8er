package serve

import (
	"fmt"

	"github.com/cfg8er/cfg8er/internal/config"
	"github.com/cfg8er/cfg8er/pkg/repository"
	"github.com/gin-gonic/gin"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/urfave/cli.v1"
)

var repoLookup map[string]config.Repo

// Run is cli action for the serve sub-command. It loads the config, clones the repos on startup, and starts the go-gin based listener.
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
	if err := cloneFetchRepos(); err != nil {
		return err
	}

	router := newRouter()

	return router.Run(c.String("listen"))
}

// cloneFetchRepos iterates over over the repoLookup global cloning repos or fetching the latest objects
// if the repo has already been cloned. Ignores NoErrAlreadyUpToDate error on fetch.
func cloneFetchRepos() error {
	for n, r := range repoLookup {
		if r.ClonedRepo == (repository.Repository{}) {
			fmt.Printf("Cloning repo %s, %s\n", n, r.URL)
			clonedRepo, err := repository.CloneBare(r.URL)
			if err != nil {
				return err
			}
			r.ClonedRepo = clonedRepo
			repoLookup[n] = r
		} else {
			fmt.Printf("Fetch latest objects from repo %s, %s\n", n, r.URL)
			err := r.ClonedRepo.Fetch(&git.FetchOptions{})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return err
			}
		}
	}

	return nil
}
