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
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"

	containershipappv1beta2 "github.com/relativitydev/containership/api/v1beta2"
	"github.com/relativitydev/containership/pkg/processor"
)

// Time in seconds for regular looping of CMO configurations
var interval = 600

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

	/*
		  Get registry credentials. This will only look at the first RegistryConfig returned.
			Supporting multiple RegistriesConfigs and a consistent registry promotion order will
			require more complex logic. For now, just one RegistriesConfig should exist.
	*/
	for _, registry := range registryConfigs.Items[0].Spec.Registries {
		config := &processor.RegistryCredentials{
			LoginURI: registry.URI,
		}

		if registry.SecretName != "" {
			secretInstance := &corev1.Secret{}

			err = r.Client.Get(ctx, types.NamespacedName{
				Name:      registry.SecretName,
				Namespace: req.Namespace,
			}, secretInstance)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, fmt.Sprintf("Error reading secret %s - requeue the request", registry.SecretName))
			}

			err = getRegistryCredentials(secretInstance, config)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, "Error reading secret context - requeue the request")
			}
		}

	}

	// TODO: Call image promotion processor logic

	return ctrl.Result{RequeueAfter: time.Second * time.Duration(interval)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContainerManagementObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&containershipappv1beta2.ContainerManagementObject{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}

// There is no need to reconcile deleted CMOs
func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// getRegistryCredentials get the creds for each registry's auth config
func getRegistryCredentials(secret *corev1.Secret, obj *processor.RegistryCredentials) error {
	// Base64 decode the secret property value
	rawCreds := secret.Data[".dockerconfigjson"]
	if len(rawCreds) == 0 {
		rawCreds = []byte(secret.StringData[".dockerconfigjson"])
	}

	/*
		Example docker credentials
		https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#log-in-to-docker
			{
				"auths": {
					"containership-docker": {
						"auth": "putyourauthhere" // should be base64 encoded
					}
				}
			}
	*/
	authsWrapper := struct {
		Auths map[string]map[string]string `json:"auths"`
	}{}

	err := json.Unmarshal(rawCreds, &authsWrapper)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading secret %s contents", secret.Name))
	}

	for _, value := range authsWrapper.Auths {
		authRaw, err := base64.StdEncoding.DecodeString(value["auth"])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error base64 decoding secret %s", secret.Name))
		}

		auth := string(authRaw)
		parts := strings.Split(auth, ":")

		obj.Username = parts[0]
		obj.Password = parts[1]
	}

	return nil
}
