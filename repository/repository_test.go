package repository

import (
	"bytes"
	"io/ioutil"
	"testing"

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
		want    Repository
		wantErr bool
	}{
		{
			name:    "Clone a repo",
			args:    args{URL: exampleRepo},
			wantErr: false,
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
					return
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

func TestRepository_FileOpenAtCommit(t *testing.T) {
	r, err := CloneBare(exampleRepo)

	if err != nil {
		t.Errorf("Repository.FileOpenAtCommit() error = %v", err)
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
			got, err := r.FileOpenAtCommit(tt.args.path, tt.args.hash)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Repository.FileOpenAtCommit() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}
			defer got.Close()

			gotContents, err := ioutil.ReadAll(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpenAtCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(gotContents, tt.want) {
				t.Errorf("Repository.FileOpenAtCommit() = %v, want %v", gotContents, tt.want)
			}

		})
	}
}
