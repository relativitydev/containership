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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"

	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
	"github.com/relativitydev/containership/pkg/processor"
)

// ContainerManagementObjectReconciler reconciles a ContainerManagementObject object
type ContainerManagementObjectReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=containership.app,resources=containermanagementobjects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=containership.app,resources=containermanagementobjects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=containership.app,resources=containermanagementobjects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ContainerManagementObject object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ContainerManagementObjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("containermanagementobject", req.NamespacedName)

	// Get the ContainerManagementObject
	instance := &containershipappv1beta2.ContainerManagementObject{}

	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, errors.Wrap(err, "Error reading containermanagementobject - requeue the request")
	}

	// Get the registryConfigs
	registryConfigs := &containershipappv1beta2.RegistriesConfigList{}

	err = r.Client.List(ctx, registryConfigs, client.InNamespace(req.Namespace))
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, errors.Wrap(err, "Error reading registriesConfig list - requeue the request")
	}

	// Get registry credentials
	secretNames := []string{}

	for _, registryConfig := range registryConfigs.Items {
		for _, registry := range registryConfig.Spec.Registries {
			if registry.SecretName != "" {
				secretNames = append(secretNames, registry.SecretName)
			}
		}
	}

	registryCredentialsDict := map[string]processor.RegistryCredentials{}

	for _, secretName := range secretNames {
		secretInstance := &corev1.Secret{}

		err = r.Client.Get(ctx, types.NamespacedName{
			Name:      secretName,
			Namespace: req.Namespace,
		}, secretInstance)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, fmt.Sprintf("Error reading secret %s - requeue the request", secretName))
		}

		getRegistryCredentials(ctx, secretInstance, registryCredentialsDict)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContainerManagementObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&containershipappv1beta2.ContainerManagementObject{}).
		Complete(r)
}

// getRegistryCredentials get the creds for each registry's auth config
func getRegistryCredentials(ctx context.Context, secret *corev1.Secret, obj map[string]processor.RegistryCredentials) error {
	// Base64 decode the secret property value
	encodedCreds := secret.StringData[".dockerconfigjson"]
	rawCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error base64 decoding secret %s - requeue the request", secret.Name))
	}

	// Unmarshall the json to a golang struct
	config := &DockerConfigSecret{}
	err = json.Unmarshal(rawCreds, config)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading secret %s contents - requeue the request", secret.Name))
	}

	// Read each auth value and add to a registry dictionary
	for _, creds := range config.Auths {
		authRaw, err := base64.RawStdEncoding.DecodeString(creds.Auth)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error base64 decoding secret %s - requeue the request", secret.Name))
		}

		auth := string(authRaw)
		parts := strings.Split(auth, ":")

		obj[""] = processor.RegistryCredentials{
			Username: parts[0],
			Password: parts[1],
		}
	}

	return nil
}
