package semverref

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func TestSemverref_HighestMatch(t *testing.T) {

	verOneZeroZero, _ := semver.NewVersion("v1.0.0")
	verOneZeroOne, _ := semver.NewVersion("v1.0.1")
	verOneOneZero, _ := semver.NewVersion("v1.1.0")
	verTwoZeroOneAlpha, _ := semver.NewVersion("v2.0.1-alpha")
	verTwoZeroZero, _ := semver.NewVersion("v2.0.0")

	constOneZeroZero, _ := semver.NewConstraint("1.0.0")
	constOneZeroZeroPatch, _ := semver.NewConstraint("~1.0.0")
	constOneZeroZeroMajor, _ := semver.NewConstraint("^1.0.0")
	constTwoZeroZeroPrePatch, _ := semver.NewConstraint("~2.0.0-0")
	constNineNineNinePatch, _ := semver.NewConstraint("~9.9.9")

	coll1 := Collection{
		SemverRef{
			Ver: verOneZeroZero,
			Ref: plumbing.NewReferenceFromStrings("refs/tags/v1.0.0", "1111111111111111111111111111111111111111"),
		},
		SemverRef{
			Ver: verOneZeroOne,
			Ref: plumbing.NewReferenceFromStrings("refs/tags/v1.0.1", "2222222222222222222222222222222222222222"),
		},
		SemverRef{
			Ver: verOneOneZero,
			Ref: plumbing.NewReferenceFromStrings("refs/tags/v1.1.0", "3333333333333333333333333333333333333333"),
		},
		SemverRef{
			Ver: verTwoZeroOneAlpha,
			Ref: plumbing.NewReferenceFromStrings("refs/tags/v2.0.1-alpha", "4444444444444444444444444444444444444444"),
		},
		SemverRef{
			Ver: verTwoZeroZero,
			Ref: plumbing.NewReferenceFromStrings("refs/tags/v2.0.0", "5555555555555555555555555555555555555555"),
		},
	}
	tests := []struct {
		name       string
		c          Collection
		constraint *semver.Constraints
		want       *plumbing.Reference
		wantErr    bool
	}{
		{
			name:       "1.0.0",
			c:          coll1,
			constraint: constOneZeroZero,
			want:       plumbing.NewReferenceFromStrings("refs/tags/v1.0.0", "1111111111111111111111111111111111111111"),
			wantErr:    false,
		},
		{
			name:       "~1.0.1",
			c:          coll1,
			constraint: constOneZeroZeroPatch,
			want:       plumbing.NewReferenceFromStrings("refs/tags/v1.0.1", "2222222222222222222222222222222222222222"),
			wantErr:    false,
		},
		{
			name:       "^1.0.0",
			c:          coll1,
			constraint: constOneZeroZeroMajor,
			want:       plumbing.NewReferenceFromStrings("refs/tags/v1.1.0", "3333333333333333333333333333333333333333"),
			wantErr:    false,
		},
		{
			name:       "~2.0.0-0",
			c:          coll1,
			constraint: constTwoZeroZeroPrePatch,
			want:       plumbing.NewReferenceFromStrings("refs/tags/v2.0.1-alpha", "4444444444444444444444444444444444444444"),
			wantErr:    false,
		},
		{
			name:       "Non-existent version ~9.9.9",
			c:          coll1,
			constraint: constNineNineNinePatch,
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.HighestMatch(tt.constraint)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Collection.HighestMatch() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection.HighestMatch() = %v, want %v", got.Name(), tt.want.Name())
			}
		})
	}
}
