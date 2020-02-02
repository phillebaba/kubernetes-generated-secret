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

package v1alpha1

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var generatedsecretlog = logf.Log.WithName("generatedsecret-resource")

func (r *GeneratedSecret) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-core-phillebaba-io-v1alpha1-generatedsecret,mutating=true,failurePolicy=fail,groups=core.phillebaba.io,resources=generatedsecrets,verbs=create;update,versions=v1alpha1,name=mgeneratedsecret.kb.io

var _ webhook.Defaulter = &GeneratedSecret{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *GeneratedSecret) Default() {
	generatedsecretlog.Info("default", "name", r.Name)

	for _, d := range r.Spec.DataList {
		if d.Length == nil {
			length := int(8)
			d.Length = &length
		}

		if d.Letters == nil {
			d.Letters = new(bool)
		}
	}
}
