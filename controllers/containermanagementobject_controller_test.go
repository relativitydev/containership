/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"testing"

	"github.com/relativitydev/containership/pkg/processor"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getRegistryCredentials(t *testing.T) {
	type args struct {
		ctx    context.Context
		secret *corev1.Secret
		obj    map[string]processor.RegistryCredentials
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get registry credentials",
			args: args{
				ctx: context.TODO(),
				secret: &corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name: "test-secret",
					},
					Data: map[string][]byte{
						".dockerconfigjson": []byte("{\"auths\":{\"containership-docker\":{\"auth\":\"dGlnZXI6cGFzczExMw==\"}}}"),
					},
				},
				obj: make(map[string]processor.RegistryCredentials),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getRegistryCredentials(tt.args.ctx, tt.args.secret, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("getRegistryCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
