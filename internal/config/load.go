package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

// LoadConfig opens the filePath using micro/go-config to scan it onto
// a map[String]Repo.
func LoadConfig(filePath string) (map[string]Repo, error) {
	// Create new config
	conf := config.NewConfig()

	repos := map[string]Repo{}

	// Load file source
	f := file.WithPath(filePath)
	s := file.NewSource(f)

	if err := conf.Load(s); err != nil {
		return repos, err
	}
	defer conf.Close()

	if err := conf.Get("repositories").Scan(&repos); err != nil {
		return repos, err
	}

	return repos, nil
}
