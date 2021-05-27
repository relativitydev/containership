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
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/relativitydev/containership/api/v1beta2"
	"github.com/relativitydev/containership/pkg/processor"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	secret = &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "registry-auth-creds",
			Namespace: "default",
		},
		Type: corev1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte("{\"auths\":{\"dockerhub-containership\":{\"auth\":\"dGlnZXI6cGFzczExMw==\"}}}"),
		},
	}

	registriesConfig = &v1beta2.RegistriesConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      "registries",
			Namespace: "default",
		},
		Spec: v1beta2.RegistriesConfigSpec{
			Registries: []v1beta2.Registry{
				{
					Name:       "dockerhub-containership",
					URI:        "docker.io/relativitydev",
					SecretName: "registry-auth-creds",
				},
				{
					Name: "dockerhub-public",
					URI:  "docker.io",
					// SecretName: "", // no secret needed. This is a public registry
				},
			},
		},
	}

	containerManagementObject = &v1beta2.ContainerManagementObject{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-images",
			Namespace: "default",
		},
		Spec: v1beta2.ContainerManagementObjectSpec{
			Images: []v1beta2.Image{
				{
					SourceRepository: "docker.io/library/busybox",
					SupportedTags: []string{
						"latest",
						"musl",
						"glibc",
					},
					Destinations: []string{
						"dockerhub-containership",
					},
				},
			},
		},
	}
)

var _ = Describe("CMO Controller", func() {
	const timeout = time.Second * 20

	ctx := context.Background()

	Context("CMO", func() {
		It("Should read CMO and RegistryConfig", func() {

			By("Creating test secret")
			Expect(k8sClient.Create(ctx, secret)).Should(Succeed())

			By("Creating regsitries config")
			Expect(k8sClient.Create(ctx, registriesConfig)).Should(Succeed())

			rc := &v1beta2.RegistriesConfig{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{
					Name:      registriesConfig.Name,
					Namespace: registriesConfig.Namespace,
				}, rc)
			}, timeout, interval).Should(Succeed())

			By("Creating cmo")
			Expect(k8sClient.Create(ctx, containerManagementObject)).Should(Succeed())

			cmo := &v1beta2.ContainerManagementObject{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{
					Name:      containerManagementObject.Name,
					Namespace: containerManagementObject.Namespace,
				}, cmo)
			}, timeout, interval).Should(Succeed())
		})
	})
})

func Test_getRegistryCredentials(t *testing.T) {
	type args struct {
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
			if err := getRegistryCredentials(tt.args.secret, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("getRegistryCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
