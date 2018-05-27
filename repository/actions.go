package repository

import "gopkg.in/src-d/go-git.v4/plumbing"

// Register adds a new repository to the service
func Register(p string) *Repository {

	return nil

}

// Clone downloads a repository
func (r *Repository) Clone() error {

	return nil

}

// Fetch downloads the latest commits to a repository
func (r *Repository) Fetch() error {

	return nil

}

// Checkout makes the given reference available to the service
func (r *Repository) Checkout(ref *plumbing.Reference) error {

	return nil

}
