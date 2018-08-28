package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/Masterminds/semver"
	"github.com/cfg8er/cfg8er/pkg/repository"
	"github.com/gin-gonic/gin"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/urfave/cli.v1"
)

var configPath string
var listen string

func serve(c *cli.Context) error {

	if !c.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	repo, err := repository.CloneBare("https://github.com/cfg8er/fixture.git")
	if err != nil {
		fmt.Printf("Cannot clone repo, exiting: %s\n", err)
	}

	router.GET("/r/:repo/:ref/*path", func(c *gin.Context) {
		urlPath := path.Clean(c.Param("path"))

		var reader io.ReadCloser
		var size int64
		var fOpenErr error

		constraint, err := semver.NewConstraint(c.Param("ref"))
		if err != nil {
			reader, size, fOpenErr = repo.FileOpenAtRev(urlPath, plumbing.Revision(c.Param("ref")))
		} else {
			if ref, fOpenErr := repo.FindSemverTag(constraint); fOpenErr == nil {
				reader, size, fOpenErr = repo.FileOpenAtRef(urlPath, *ref)
			}
		}

		if fOpenErr != nil || reader == nil {
			fmt.Printf("fOpenErr: %v\n", fOpenErr)
			c.Status(http.StatusNotFound)
			return
		}

		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, path.Base(urlPath)),
		}

		c.DataFromReader(http.StatusOK, size, "text/plain", reader, extraHeaders)
	})

	return router.Run(c.String("listen"))
}

func main() {
	app := cli.NewApp()
	app.Name = "Cfg8er"
	app.Usage = "git based configuration hosting service"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "Start http service and serve configured repositories",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "config.yml",
					Usage: "Configuration file path",
				},
				cli.StringFlag{
					Name:  "listen, l",
					Value: "127.0.0.1:8080",
					Usage: "IP address and port to listen on",
				},
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "Enable debug mode",
				},
			},
			Action: serve,
		},
	}

	app.Run(os.Args)
}
