package repository

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var exampleRepo = "https://github.com/git-fixtures/basic.git"

func TestCloneBare(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Clone a repo",
			args:    args{URL: exampleRepo},
			wantErr: false,
		},
		{
			name:    "Clone a non-existent repo",
			args:    args{URL: "https://example.com/example.git"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CloneBare(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneBare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_FileOpenAtRev(t *testing.T) {
	r, err := CloneBare(exampleRepo)

	if err != nil {
		t.Errorf("Repository.FileOpenAtRev() error = %v", err)
		return
	}

	type args struct {
		path string
		rev  plumbing.Revision
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Reference not found",
			args:    args{path: "CHANGELOG", rev: plumbing.Revision("asdf1234")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Open and read CHANGELOG at refs/heads/master",
			args:    args{path: "CHANGELOG", rev: plumbing.Revision("refs/heads/master")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name:    "Open and read CHANGELOG at refs/remotes/origin/branch",
			args:    args{path: "CHANGELOG", rev: plumbing.Revision("refs/remotes/origin/branch")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name: "Open and read vendor/foo.go at refs/heads/master",
			args: args{path: "vendor/foo.go", rev: plumbing.Revision("refs/heads/master")},
			want: []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Hello, playground\")\n}\n"),
			wantErr: false,
		},
		{
			name:    "Open and read CHANGELOG at 6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			args:    args{path: "CHANGELOG", rev: plumbing.Revision("6ecf0ef2c2dffb796033e5a02219af86ec6584e5")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name:    "Open and read CHANGELOG at short ref master",
			args:    args{path: "CHANGELOG", rev: plumbing.Revision("master")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FileOpenAtRev(tt.args.path, tt.args.rev)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Repository.FileOpenAtRev() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			defer got.Close()

			gotContents, err := ioutil.ReadAll(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpenAtRev() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(gotContents, tt.want) {
				t.Errorf("Repository.FileOpenAtRev() = %v, want %v", gotContents, tt.want)
			}

		})
	}
}

func TestRepository_fileOpenAtHash(t *testing.T) {
	r, err := CloneBare(exampleRepo)

	if err != nil {
		t.Errorf("Repository.fileOpenAtHash() error = %v", err)
		return
	}

	type args struct {
		path string
		hash plumbing.Hash
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Non-existent commit hash",
			args:    args{path: "CHANGELOG", hash: plumbing.NewHash("0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Non-existent file path",
			args:    args{path: "asdf/ghjk", hash: plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Open and read CHANGELOG at 6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			args:    args{path: "CHANGELOG", hash: plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name:    "Open and read CHANGELOG at a remote branch commit e8d3ffab552895c19b9fcf7aa264d277cde33881",
			args:    args{path: "CHANGELOG", hash: plumbing.NewHash("e8d3ffab552895c19b9fcf7aa264d277cde33881")},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name: "Open and read vendor/foo.go at 6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			args: args{path: "vendor/foo.go", hash: plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5")},
			want: []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Hello, playground\")\n}\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.fileOpenAtHash(tt.args.path, tt.args.hash)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Repository.fileOpenAtHash() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}
			defer got.Close()

			gotContents, err := ioutil.ReadAll(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.fileOpenAtHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(gotContents, tt.want) {
				t.Errorf("Repository.fileOpenAtHash() = %v, want %v", gotContents, tt.want)
			}

		})
	}
}

func TestRepository_FindSemverTag(t *testing.T) {
	r, err := CloneBare("https://github.com/cfg8er/fixture.git")

	if err != nil {
		t.Errorf("Repository.FindSemverTag() error = %v", err)
		return
	}

	zeroZeroOne, _ := semver.NewConstraint("0.0.1")
	zeroZeroOnePatch, _ := semver.NewConstraint("~0.0.1")
	zeroOnePatch, _ := semver.NewConstraint("~0.1")
	oneZeroZeroPatch, _ := semver.NewConstraint("~1.0.0")
	oneZeroOnePrePatch, _ := semver.NewConstraint("~1.0.1-0")
	nineNineNinePatch, _ := semver.NewConstraint("~9.9.9")

	tests := []struct {
		name       string
		constraint *semver.Constraints
		wantHash   plumbing.Hash
		wantErr    bool
	}{
		{
			name:       "0.0.1",
			constraint: zeroZeroOne,
			wantHash:   plumbing.NewHash("bdbfd0ebc52195e74f4d748bed9adde12a275c75"),
			wantErr:    false,
		},
		{
			name:       "~0.0.1",
			constraint: zeroZeroOnePatch,
			wantHash:   plumbing.NewHash("f9beb3bc5e04eb1a33f85805e1f2c5541e6661fc"),
			wantErr:    false,
		},
		{
			name:       "~0.1",
			constraint: zeroOnePatch,
			wantHash:   plumbing.NewHash("48e8f899ab1cf3a3be36371f4161cfb897659c45"),
			wantErr:    false,
		},
		{
			name:       "~1.0.0",
			constraint: oneZeroZeroPatch,
			wantHash:   plumbing.NewHash("c204415aafecf7dd22513f3d7158d224a32763f4"),
			wantErr:    false,
		},
		{
			name:       "~1.0.0-0",
			constraint: oneZeroOnePrePatch,
			wantHash:   plumbing.NewHash("774270d020ae8e17836bc399f238b77cda990e77"),
			wantErr:    false,
		},
		{
			name:       "~9.9.9 non-existant version",
			constraint: nineNineNinePatch,
			wantHash:   plumbing.Hash{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindSemverTag(tt.constraint)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Repository.FindSemverTag() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.Hash() != tt.wantHash {
				t.Errorf("Repository.FindSemverTag() = %v, want %v", got, tt.wantHash)
			}
		})
	}
}
