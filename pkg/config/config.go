package config

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
)

type Repository struct {
	Name              string
	URI               string
	UpdateFrequency   int      `json:"update_frequency"`
	EnableUpdateAPI   bool     `json:"enable_update_api"`
	EnableSemversTags bool     `json:"enable_semvers_tags"`
	EnableTags        bool     `json:"enable_tags"`
	EnableCommits     bool     `json:"enable_commits"`
	WhitelistRefs     []string `json:"whitelist_refs"`
	BlacklistRefs     []string `json:"blacklist_refs"`
	AllowHosts        []string `json:"allow_hosts"`
	GpgVerifyCommit   bool     `json:"gpg_verify_commit"`
	GpgVerifyTag      bool     `json:"gpg_verify_tag"`
	GpgAllowIds       []string `json:"gpg_allow_ids"`
}

type Repositories []Repository

func LoadConfig(filePath string) (Repositories, error) {
	// Create new config
	conf := config.NewConfig()

	repos := Repositories{}

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
