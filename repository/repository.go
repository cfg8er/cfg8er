package repository

import (
	"fmt"
	"io"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Repository is an extended go-git Repository
type Repository struct{ *git.Repository }

// CloneBare downloads the repository as a bare repo
func CloneBare(URL string) (Repository, error) {
	// Git objects storer based on memory
	storer := memory.NewStorage()
	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	repo, err := git.Clone(storer, nil, &git.CloneOptions{
		URL: URL,
	})
	if err != nil {
		return Repository{}, err
	}
	return Repository{repo}, nil
}

//FileOpenAtRef opens a file at a given path at a given reference
func (r *Repository) FileOpenAtRef(path string, refName plumbing.ReferenceName) (io.ReadCloser, error) {
	ref, err := r.Reference(refName, true)
	if err != nil {
		return nil, fmt.Errorf("refName Lookup of %s: %s", refName, err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("Commit object lookup of %v: %s", ref.Hash(), err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("Get tree of commit %v: %s", commit.TreeHash, err)
	}

	te, err := tree.FindEntry(path)
	if err != nil {
		return nil, fmt.Errorf("Find path in tree %s: %s", path, err)
	}

	o, err := r.BlobObject(te.Hash)
	if err != nil {
		return nil, fmt.Errorf("Blob object lookup of %v: %s", te.Hash, err)

	}

	return o.Reader()
}

// Fetch downloads the latest commits to a repository
func (r *Repository) Fetch() error {

	return nil
}

// Checkout makes the given reference available to the service
func (r *Repository) Checkout(ref *plumbing.Reference) error {

	return nil
}
