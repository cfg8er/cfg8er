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

func TestRepository_FileOpenAtRef(t *testing.T) {
	r, err := CloneBare(exampleRepo)

	if err != nil {
		t.Errorf("Repository.FileOpenAtRef() error = %v", err)
		return
	}

	type args struct {
		path    string
		refName plumbing.ReferenceName
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Open and read CHANGELOG at master",
			args:    args{path: "CHANGELOG", refName: "refs/heads/master"},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name:    "Open and read CHANGELOG at a remote branch",
			args:    args{path: "CHANGELOG", refName: "refs/remotes/origin/branch"},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
		{
			name: "Open and read vendor/foo.go at master",
			args: args{path: "vendor/foo.go", refName: "refs/heads/master"},
			want: []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n	fmt.Println(\"Hello, playground\")\n}\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FileOpenAtRef(tt.args.path, tt.args.refName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpenAtRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()

			gotContents, err := ioutil.ReadAll(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpenAtRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(gotContents, tt.want) {
				t.Errorf("Repository.FileOpenAtRef() = %v, want %v", gotContents, tt.want)
			}

		})
	}
}
