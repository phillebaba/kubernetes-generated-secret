package controllers

import (
	"context"
	"github.com/pkg/errors"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

	// Get the reconciled GeneratedSecret
	var gs corev1alpha1.GeneratedSecret
	if err := r.Get(ctx, req.NamespacedName, &gs); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	s := corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: gs.Name, Namespace: gs.Namespace}}
	_, err := ctrl.CreateOrUpdate(ctx, r, &s, func() error {
		sm := *gs.Spec.SecretMeta.DeepCopy()
		sm.Name = gs.Name
		sm.Namespace = gs.Namespace
		s.ObjectMeta = sm

		// Generate secret data
		sd := make(map[string][]byte)
		for _, d := range gs.Spec.DataList {
			// Skip if value already exists
			if val, ok := s.Data[d.Key]; ok {
				sd[d.Key] = val
				continue
			}

			randString, err := crypto.GenerateRandomASCIIString(d.Length, d.ValueOptions)
			if err != nil {
				return err
			}

			sd[d.Key] = []byte(randString)
		}
		s.Data = sd
		return controllerutil.SetControllerReference(&gs, &s, r.Scheme)
	})

	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "unable to create or update Secret")
	}

	return ctrl.Result{}, nil
}

func (r *GeneratedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.GeneratedSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
