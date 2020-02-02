/*

licensed under the apache license, version 2.0 (the "license");
you may not use this file except in compliance with the license.
you may obtain a copy of the license at

    http://www.apache.org/licenses/license-2.0

unless required by applicable law or agreed to in writing, software
distributed under the license is distributed on an "as is" basis,
without warranties or conditions of any kind, either express or implied.
see the license for the specific language governing permissions and
limitations under the license.
*/

package controllers

import (
	"context"
	"encoding/base64"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	corev1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

var _ = Describe("Generated Secret Controller", func() {
	const timeout = time.Second * 30
	const interval = time.Second * 1

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	Context("Simple Secret", func() {
		It("Should create successfully", func() {
			key := types.NamespacedName{
				Name:      "foo",
				Namespace: "default",
			}

			length := int(8)
			created := &corev1alpha1.GeneratedSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: corev1alpha1.GeneratedSecretSpec{
					DataList: []corev1alpha1.GeneratedSecretData{corev1alpha1.GeneratedSecretData{Length: &length, Key: "foo"}},
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())

			// Get generated secret
			By("Expecting secret to be created")
			Eventually(func() bool {
				f := &corev1.Secret{}
				k8sClient.Get(context.Background(), key, f)
				decoded, _ := base64.StdEncoding.DecodeString(string(f.Data["foo"]))
				return len(string(decoded)) == length
			}, timeout, interval).Should(BeTrue())

			// Delete
			By("Expecting to delete successfully")
			Eventually(func() error {
				f := &corev1alpha1.GeneratedSecret{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			By("Expecting to delete finish")
			Eventually(func() error {
				f := &corev1alpha1.GeneratedSecret{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
