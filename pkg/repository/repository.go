package repository

import (
	"fmt"
	"io"
	"path"

	"github.com/Masterminds/semver"
	"github.com/cfg8er/cfg8er/pkg/repository/semverref"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// Repository is an extended go-git Repository
type Repository struct{ *git.Repository }

// CloneBare downloads the repository as a bare repo including all tags
func CloneBare(URL string) (Repository, error) {
	// Git objects storer based on memory
	storer := memory.NewStorage()

	repo, err := git.Clone(storer, nil, &git.CloneOptions{
		URL:  URL,
		Tags: git.TagMode(2),
	})
	if err != nil {
		return Repository{}, err
	}
	return Repository{repo}, nil
}

// FileOpenAtRev opens a file at a given path at a given Git revision, eg.
// https://kernel.org/pub/software/scm/git/docs/gitrevisions.html. Returns an
// open io.ReadCloser, file size, and error.
func (r *Repository) FileOpenAtRev(filePath string, rev plumbing.Revision) (io.ReadCloser, int64, error) {
	ref, err := r.ResolveRevision(rev)
	if err != nil {
		return nil, 0, fmt.Errorf("Revision resolve of %s: %s", ref, err)
	}
	return r.fileOpenAtHash(filePath, *ref)
}

// FileOpenAtRef opens a file at a given path at given reference. Returns an open io.ReadCloser,
// file size, and error.
func (r *Repository) FileOpenAtRef(filePath string, ref plumbing.Reference) (io.ReadCloser, int64, error) {
	return r.fileOpenAtHash(filePath, ref.Hash())
}

// fileOpenAtHash opens a file at a given path at a given hash. Returns an open io.ReadCloser,
// file size, and error.
func (r *Repository) fileOpenAtHash(filePath string, hash plumbing.Hash) (io.ReadCloser, int64, error) {
	commit, err := r.CommitObject(hash)
	if err != nil {
		return nil, 0, fmt.Errorf("Commit object of %v: %s", hash, err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, 0, fmt.Errorf("Tree of commit %v: %s", commit.TreeHash, err)
	}

	// If filePath has a leading slash remove it as tree entries don't have a leading slash.
	if path.IsAbs(filePath) {
		filePath = filePath[1:]
	}

	entry, err := tree.FindEntry(filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("Path in tree %s: %s", filePath, err)
	}

	object, err := r.BlobObject(entry.Hash)
	if err != nil {
		return nil, 0, fmt.Errorf("Blob object of %v: %s", entry.Hash, err)
	}

	reader, err := object.Reader()

	return reader, object.Size, err
}

// FindSemverTag iterates through the repository's tags looking for tags that
// follow semantic versioning (https://semver.org). Returns the highest version
// tag that meets the supplied contraint. Silently ignores tags that aren't
// parsable as a semantic version.
func (r *Repository) FindSemverTag(c *semver.Constraints) (*plumbing.Reference, error) {
	tagsIter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	coll := semverref.Collection{}

	if err := tagsIter.ForEach(func(t *plumbing.Reference) error {
		v, err := semver.NewVersion(t.Name().Short())
		if err != nil {
			return nil // Ignore errors and thus tags that aren't parsable as a semver
		}

		// No way to a priori find the length of tagsIter so append to the collection.
		coll = append(coll, semverref.SemverRef{Ver: v, Ref: t})
		return nil
	}); err != nil {
		return nil, err
	}

	return coll.HighestMatch(c)
}
