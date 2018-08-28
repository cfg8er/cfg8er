package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Repositories
		wantErr bool
	}{
		{
			name:    "Non-existant path",
			args:    args{filePath: "/non-existant/path"},
			want:    Repositories{},
			wantErr: true,
		},
		/*{
			name:    "Empty json config",
			args:    args{filePath: "../../fixtures/config/empty.json"},
			want:    Repositories{},
			wantErr: false,
		},*/
		{
			name: "Load example yaml config",
			args: args{filePath: "../../fixtures/config/config.yaml"},
			want: Repositories{
				{
					Name:              "cfg8er-fixture",
					URI:               "https://github.com/cfg8er/fixture.git",
					UpdateFrequency:   600,
					EnableUpdateAPI:   false,
					EnableSemversTags: true,
					EnableTags:        true,
					EnableCommits:     false,
					WhitelistRefs:     []string{"v1.0.*"},
					BlacklistRefs:     []string{"v0.*"},
					AllowHosts:        []string{"127.0.0.1/8"},
					GpgVerifyCommit:   true,
					GpgVerifyTag:      true,
					GpgAllowIds:       []string{"29DF880B"},
				},
			},
			wantErr: false,
		},
		{
			name: "Load example json config",
			args: args{filePath: "../../fixtures/config/config.json"},
			want: Repositories{
				{
					Name:              "cfg8er-fixture",
					URI:               "https://github.com/cfg8er/fixture.git",
					UpdateFrequency:   600,
					EnableUpdateAPI:   false,
					EnableSemversTags: true,
					EnableTags:        true,
					EnableCommits:     false,
					WhitelistRefs:     []string{"v1.0.*"},
					BlacklistRefs:     []string{"v0.*"},
					AllowHosts:        []string{"127.0.0.1/8"},
					GpgVerifyCommit:   true,
					GpgVerifyTag:      true,
					GpgAllowIds:       []string{"29DF880B"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
