package semverref

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//SemverRef groups a Semantic Version and matching a Git Reference
type SemverRef struct {
	Ver *semver.Version
	Ref *plumbing.Reference
}

// Collection is a slice of SemverRef implimented for sorting.
// See the sort package for more details. https://golang.org/pkg/sort/
// Implimentation extended from Masterminds/semver.
type Collection []SemverRef

// Len returns the length of a collection.
func (c Collection) Len() int {
	return len(c)
}

// Less is needed for the sort interface to compare two SemverRef objects
// on the slice via the Ver attribute. If checks if the semantic version
// is less than the other.
func (c Collection) Less(i, j int) bool {
	return c[i].Ver.LessThan(c[j].Ver)
}

// Swap is needed for the sort interface to replace the SemverRef objects
// at two different positions in the slice.
func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// HighestMatch sorts the collection by the Ver *semver.Version attribute.
// Iterates over the Collection, and returns the last, thus highest, version
// matching the supplied constraint. Returns a nil Reference and an error
// if no tag is found that matches the constraint.
func (c Collection) HighestMatch(con *semver.Constraints) (*plumbing.Reference, error) {
	sort.Sort(c)

	var lastMatch *plumbing.Reference

	for _, sr := range c {
		match := con.Check(sr.Ver)

		if !match && lastMatch != nil {
			return lastMatch, nil
		} else if match {
			lastMatch = sr.Ref
		}
	}

	if lastMatch == nil {
		return nil, fmt.Errorf("No matching tag found")
	}
	return lastMatch, nil
}
