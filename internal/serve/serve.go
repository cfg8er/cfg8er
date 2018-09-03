package serve

import (
	"fmt"

	"github.com/cfg8er/cfg8er/internal/config"
	"github.com/cfg8er/cfg8er/pkg/repository"
	"github.com/gin-gonic/gin"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/urfave/cli.v1"
)

var repoLookup map[string]*config.Repo

// Run is the cli action for the serve sub-command. It loads the config, clones
// the repos on startup, and starts the go-gin based listener.
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
	updateRepoCh := make(chan *config.Repo, 100)
	defer close(updateRepoCh)

	go cloneRepos(updateRepoCh)

	for n := range repoLookup {
		updateRepoCh <- repoLookup[n]
	}

	router := newRouter()

	return router.Run(c.String("listen"))
}

// cloneRepos iterates over over the repoLookup global cloning repos or fetching the latest objects
// if the repo has already been cloned. Ignores NoErrAlreadyUpToDate error on fetch.
func cloneRepos(updateRepoCh chan *config.Repo) {
	for r := range updateRepoCh {
		if r.ClonedRepo == (repository.Repository{}) {
			fmt.Printf("Cloning repo %s\n", r.URL)
			clonedRepo, err := repository.CloneBare(r.URL)
			if err != nil {
				fmt.Printf("Error: Cloning repo %s: %v\n", r.URL, err)
			}
			r.ClonedRepo = clonedRepo
		} else {
			fmt.Printf("Fetch latest objects from repo %s\n", r.URL)
			err := r.ClonedRepo.Fetch(&git.FetchOptions{})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				fmt.Printf("Error: Fetching objects from repo %s: %v\n", r.URL, err)
			}
		}
	}
}
