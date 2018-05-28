package repository

import (
	billy "gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Repository is an extended go-git Repository
type Repository struct{ *git.Repository }

// Clone downloads the repository
func Clone(path string) (Repository, error) {
	// Filesystem abstraction based on memory
	fs := memfs.New()
	// Git objects storer based on memory
	storer := memory.NewStorage()
	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: path,
	})
	if err != nil {
		return Repository{}, err
	}
	return Repository{repo}, nil
}

//FileOpen opens a given path within the repository's worktree
func (r *Repository) FileOpen(path string) (billy.File, error) {
	wt, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	fs := wt.Filesystem
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Fetch downloads the latest commits to a repository
func (r *Repository) Fetch() error {

	return nil
}

// Checkout makes the given reference available to the service
func (r *Repository) Checkout(ref *plumbing.Reference) error {

	return nil
}
