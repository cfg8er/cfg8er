package repository

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var exampleRepo = "https://github.com/git-fixtures/basic.git"

func TestClone(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    Repository
		wantErr bool
	}{
		{
			name:    "Clone a repo",
			args:    args{p: exampleRepo},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Clone(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_FileOpen(t *testing.T) {
	r, err := Clone(exampleRepo)

	if err != nil {
		t.Errorf("Repository.FileOpen() error = %v", err)
		return
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Open and read CHANGELOG",
			args:    args{path: "CHANGELOG"},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.FileOpen(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer got.Close()

			gotContents, err := ioutil.ReadAll(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.FileOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !bytes.Equal(gotContents, tt.want) {
				t.Errorf("Repository.FileOpen() = %s, want %s", gotContents, tt.want)
			}
		})
	}
}
