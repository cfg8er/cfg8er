package repository

import git "gopkg.in/src-d/go-git.v4"

// Repository is a go-git Repository
type Repository struct {
	r git.Repository
}
