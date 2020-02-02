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
	"math/big"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

// https://gist.github.com/denisbrodbeck/635a644089868a51eccd6ae22b2eb800
func generateRandomASCIIString(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		// Make sure that the number/byte/letter is inside
		// the range of printable ASCII characters (excluding space and DEL)
		if n > 32 && n < 127 {
			result += string(n)
		}
	}
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
		log.Error(err, "unable to fetch GeneratedSecret")
		return ctrl.Result{}, client.IgnoreNotFound(err)
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

	// Generate secrets values
	secretData := make(map[string][]byte)
	for _, d := range gs.Spec.DataList {
		randString, _ := generateRandomASCIIString(*d.Length)
		randString = base64.URLEncoding.EncodeToString([]byte(randString))
		secretData[d.Key] = []byte(randString)
	}

	// Create secret resource
	s := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        gs.Name,
			Namespace:   gs.Namespace,
		},
		Data: secretData,
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
