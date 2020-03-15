package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *GeneratedSecretReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	ctx := context.Background()

	// Fetch generated secret
	gs := &corev1alpha1.GeneratedSecret{}
	if err := r.Get(ctx, req.NamespacedName, gs); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Set GeneratedSecret status
	defer func() {
		if err != nil {
			gs.Status.State = corev1alpha1.Failed
		}

		r.Status().Update(ctx, gs)
	}()

	if gs.Status.State == "" {
		gs.Status.State = corev1alpha1.Generating
		err := r.Status().Update(ctx, gs)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Set defaults
	if ok := checkAndSetDefaults(gs); ok == false {
		if err := r.Update(ctx, gs); err != nil {
			return ctrl.Result{}, err
		}

		r.Log.Info("Setting default values and requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if secret already exists
	s := &corev1.Secret{}
	err = r.Get(ctx, req.NamespacedName, s)
	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}
	if apierrors.IsNotFound(err) == false && hasOwnerReference(gs.UID, s.ObjectMeta.OwnerReferences) == false {
		gs.Status.State = corev1alpha1.Conflict
		return ctrl.Result{Requeue: true}, nil
	}

	// Update or create secret
	s = &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: gs.Name, Namespace: gs.Namespace}}
	_, err = ctrl.CreateOrUpdate(ctx, r, s, func() error {
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

			randString, err := crypto.GenerateRandomASCIIString(*d.Length, d.Exclude)
			if err != nil {
				return err
			}

			sd[d.Key] = []byte(randString)
		}
		s.Data = sd

		gs.Status.State = corev1alpha1.Generated
		return controllerutil.SetControllerReference(gs, s, r.Scheme)
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

// Checks if OwnerReferences contains UID
func hasOwnerReference(UID types.UID, ownerReferences []metav1.OwnerReference) bool {
	has := false
	for _, ownerReference := range ownerReferences {
		if ownerReference.UID == UID {
			has = true
		}
	}

	return has
}

// Checks if all values are set and sets them if not
func checkAndSetDefaults(gs *corev1alpha1.GeneratedSecret) bool {
	ok := true

	for i := 0; i < len(gs.Spec.DataList); i++ {
		d := &gs.Spec.DataList[i]

		if d.Length == nil {
			length := int(10)
			d.Length = &length
			ok = false
		}
	}

	return ok
}
