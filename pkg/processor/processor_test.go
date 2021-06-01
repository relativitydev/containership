package processor

import (
	"fmt"
	"reflect"
	"testing"

	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
)

func Test_populateTagArrays(t *testing.T) {
	type args struct {
		destinationTagsList []string
		supportedTags       []string
	}

	tests := []struct {
		name         string
		args         args
		wantToDelete []string
		wantToAdd    []string
	}{
		{
			name: "Test supported and existing tags match",
			args: args{
				destinationTagsList: []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
				supportedTags:       []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
			},
			wantToDelete: []string{},
			wantToAdd:    []string{},
		},
		{
			name: "Test no supported tags",
			args: args{
				destinationTagsList: []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
				supportedTags:       []string{},
			},
			wantToDelete: []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
			wantToAdd:    []string{},
		},
		{
			name: "Test no destination tags",
			args: args{
				destinationTagsList: []string{},
				supportedTags:       []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
			},
			wantToDelete: []string{},
			wantToAdd:    []string{"a", "b", "c", "1", " ", "", "#", "\\", "^", "=", ".3", "`", "~", "<", ".", "'", ";", "_", "@", "()"},
		},
		{
			name: "Test empty arrays",
			args: args{
				destinationTagsList: []string{},
				supportedTags:       []string{},
			},
			wantToDelete: []string{},
			wantToAdd:    []string{},
		},
		{
			name: "Test typical busybox example",
			args: args{
				destinationTagsList: []string{"1.32.0", "1.32", "latest", "1-musl", "1.32.0-musl"},
				supportedTags:       []string{"1.32.0", "latest", "musl", "glibc", "1.32.0-uclibc"},
			},
			wantToDelete: []string{"1.32", "1-musl", "1.32.0-musl"},
			wantToAdd:    []string{"musl", "glibc", "1.32.0-uclibc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toDelete, toAdd := populateTagArrays(tt.args.destinationTagsList, tt.args.supportedTags)
			if !reflect.DeepEqual(toDelete, tt.wantToDelete) {
				t.Errorf("populateTagArrays() toDelete = %v, wantToDelete %v", toDelete, tt.wantToDelete)
			}
			if !reflect.DeepEqual(toAdd, tt.wantToAdd) {
				t.Errorf("populateTagArrays() toAdd = %v, wantToDelete %v", toAdd, tt.wantToDelete)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		images     []containershipappv1beta2.Image
		registries []RegistryCredentials
	}

	tests := []struct {
		name        string
		args        args
		mockReturns []interface{}
		wantErr     bool
	}{
		{
			name: "busybox",
			args: args{
				images: []containershipappv1beta2.Image{
					{
						SourceRepository: "busybox",
						TargetRepository: "relativitydev/busybox",
						SupportedTags: []string{
							"latest",
							"musl",
							"glibc",
						},
					},
				},
				registries: []RegistryCredentials{
					{
						Hostname: "index.docker.io",
						Username: "",
						Password: "",
					},
				},
			},
			mockReturns: []interface{}{
				[]string{"latest", "musl", "tagToDelete"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		mockRegistryClient := &MockRegistryClient{}

		mockRegistryClient.On("listTags", tt.args.images[0].TargetRepository, tt.args.registries[0]).Return(tt.mockReturns[0], nil)

		mockRegistryClient.On("copy", fmt.Sprintf("%s:%s", tt.args.images[0].SourceRepository, "glibc"), fmt.Sprintf("%s/%s:%s", tt.args.registries[0].Hostname, tt.args.images[0].TargetRepository, "glibc"), tt.args.registries[0]).Return(nil)

		mockRegistryClient.On("delete", fmt.Sprintf("%s/%s:%s", tt.args.registries[0].Hostname, tt.args.images[0].TargetRepository, "tagToDelete"), tt.args.registries[0]).Return(nil)

		t.Run(tt.name, func(t *testing.T) {
			if err := Run(mockRegistryClient, tt.args.images, tt.args.registries); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setTargetRepository(t *testing.T) {
	type args struct {
		targetRepository *string
		sourceRepository string
	}

	empty := ""
	set := "hello-world"

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "TargetRepository is empty",
			args: args{
				targetRepository: &empty,
				sourceRepository: "library/busybox",
			},
			wantErr: false,
			want:    "library/busybox",
		},
		{
			name: "TargetRepository is set",
			args: args{
				targetRepository: &set,
				sourceRepository: "library/busybox",
			},
			wantErr: false,
			want:    set,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setTargetRepository(tt.args.targetRepository, tt.args.sourceRepository); (err != nil) != tt.wantErr {
				t.Errorf("setTargetRepository() error = %v, wantErr %v", err, tt.wantErr)
			}

			if *tt.args.targetRepository != tt.want {
				t.Errorf("target repository was incorrectly set")
			}
		})
	}
}
