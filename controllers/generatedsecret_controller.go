package controllers

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
	"github.com/phillebaba/kubernetes-generated-secret/crypto"
)

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
		randString, err := crypto.GenerateRandomASCIIString(d.Length, d.ValueOptions)
		if err != nil {
			return ctrl.Result{}, err
		}

		secretData[d.Key] = []byte(randString)
	}

	// Create secret resource
	omc := *gs.Spec.SecretMeta.DeepCopy()
	omc.Name = gs.Name
	omc.Namespace = gs.Namespace
	s := &corev1.Secret{
		ObjectMeta: omc,
		Data:       secretData,
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
