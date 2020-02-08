package controllers

import (
	"context"
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
				Name:      "default",
				Namespace: "default",
			}

			created := &corev1alpha1.GeneratedSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: corev1alpha1.GeneratedSecretSpec{
					SecretMeta: metav1.ObjectMeta{
						Name:        "try to override",
						Namespace:   "dummy",
						Labels:      map[string]string{"test": "label"},
						Annotations: map[string]string{"test": "annotations"},
					},
					DataList: []corev1alpha1.GeneratedSecretData{{Key: "foo", Length: 50, ValueOptions: []corev1alpha1.ValueOption{"Uppercase", "Lowercase", "Numbers", "Symbols"}}},
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())

			// Get generated secret
			By("Expecting secret to be created")
			f := &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).Should(Succeed())
			Expect(f.ObjectMeta.Name).To(Equal(key.Name))
			Expect(f.ObjectMeta.Namespace).To(Equal(key.Namespace))
			Expect(f.ObjectMeta.Labels).To(Equal(created.Spec.SecretMeta.Labels))
			Expect(f.ObjectMeta.Annotations).To(Equal(created.Spec.SecretMeta.Annotations))

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
