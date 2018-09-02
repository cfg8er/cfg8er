package config

import (
	"github.com/cfg8er/cfg8er/pkg/repository"
)

// Repo represents the contents of cfg8er-server configuration file
// with the additional of tracking the cloned repo. This is the primary
// type that is passed around the serve package.
type Repo struct {
	URL               string
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
	ClonedRepo        repository.Repository
}
