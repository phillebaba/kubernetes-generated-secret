/*

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
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GeneratedSecretReconciler reconciles a GeneratedSecret object
type GeneratedSecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.phillebaba.io,resources=generatedsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.phillebaba.io,resources=generatedsecrets/status,verbs=get;update;patch

func (r *GeneratedSecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("generatedsecret", req.NamespacedName)

	// Get the reconciled GeneratedSecret
	var gs corev1alpha1.GeneratedSecret
	if err := r.Get(ctx, req.NamespacedName, &gs); err != nil {
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	// Make sure secret does not already exist
	var es corev1.Secret
	err := r.Get(ctx, req.NamespacedName, &es)
	if err != nil && apierrs.IsNotFound(err) == false {
		return ctrl.Result{}, err
	}
	if apierrs.IsNotFound(err) == false && &es != nil {
		return ctrl.Result{}, err
	}

	// Gernete secrets
	sd := make(map[string]string)
	for _, d := range gs.Spec.DataList {
		fmt.Printf("%v\n\n", d)
		randString, _ := generateRandomString(*d.Length)
		sd[d.Key] = randString
	}

	// Create secret resource
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        gs.Name,
			Namespace:   gs.Namespace,
		},
		StringData: sd,
	}

	if err := ctrl.SetControllerReference(&gs, s, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, s); err != nil {
		log.Error(err, "unable to create Secret for GeneratedSecret", "secret", s)
		return ctrl.Result{}, err
	}

	log.V(1).Info("created Secret from GeneratedSecret", "secret", s)

	return ctrl.Result{}, nil
}

func (r *GeneratedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.GeneratedSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
