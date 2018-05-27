package repository

import (
	"bytes"
	"io/ioutil"
	"testing"

	git "gopkg.in/src-d/go-git.v4"
)

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
			args:    args{p: "https://github.com/git-fixtures/basic.git"},
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
	r, _ := Clone("https://github.com/git-fixtures/basic.git")

	type fields struct {
		r *git.Repository
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "Open and read CHANGELOG",
			fields:  fields{r: r.r},
			args:    args{path: "CHANGELOG"},
			want:    []byte("Initial changelog\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				r: tt.fields.r,
			}
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
