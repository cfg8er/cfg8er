package serve

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func newRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/r/:repo/:version/*path", getRepoVersionPath)

	return router
}

func getRepoVersionPath(c *gin.Context) {
	repo := c.Param("repo")
	version := c.Param("version")
	urlPath := path.Clean(c.Param("path"))

	if repo == "" || version == "" || urlPath == "" {
		c.Status(http.StatusNotFound)
		return
	}

	r, ok := repoLookup[repo]

	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	reader, size, err := r.ClonedRepo.FileOpenAtSemVer(urlPath, version)

	if err != nil || reader == nil {
		c.Status(http.StatusNotFound)
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, path.Base(urlPath)),
	}

	c.DataFromReader(http.StatusOK, size, "text/plain", reader, extraHeaders)
}
