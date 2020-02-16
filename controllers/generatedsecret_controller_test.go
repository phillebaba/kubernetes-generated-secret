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
					DataList: []corev1alpha1.GeneratedSecretData{{Key: "foo"}},
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), created)).Should(Succeed())
			By("Expecting secret to be created")
			f := &corev1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).Should(Succeed())
			Expect(f.ObjectMeta.Name).To(Equal(key.Name))
			Expect(f.ObjectMeta.Namespace).To(Equal(key.Namespace))
			Expect(f.ObjectMeta.Labels).To(Equal(created.Spec.SecretMeta.Labels))
			Expect(f.ObjectMeta.Annotations).To(Equal(created.Spec.SecretMeta.Annotations))
			Expect(len(f.Data)).To(Equal(len(created.Spec.DataList)))

			// Update
			updated := &corev1alpha1.GeneratedSecret{}
			Expect(k8sClient.Get(context.Background(), key, updated)).Should(Succeed())

			updated.Spec.SecretMeta = metav1.ObjectMeta{
				Name:        "override-name",
				Namespace:   "override-namespace",
				Labels:      map[string]string{"test": "label"},
				Annotations: map[string]string{"test": "annotations"},
			}
			updated.Spec.DataList = append(updated.Spec.DataList, corev1alpha1.GeneratedSecretData{Key: "bar"})
			Expect(k8sClient.Update(context.Background(), updated)).Should(Succeed())

			time.Sleep(100 * time.Millisecond)

			By("Expecting secret to be updated")
			Eventually(func() error {
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).Should(Succeed())
			Expect(f.ObjectMeta.Name).To(Equal(key.Name))
			Expect(f.ObjectMeta.Namespace).To(Equal(key.Namespace))
			Expect(f.ObjectMeta.Labels).To(Equal(updated.Spec.SecretMeta.Labels))
			Expect(f.ObjectMeta.Annotations).To(Equal(updated.Spec.SecretMeta.Annotations))
			Expect(len(f.Data)).To(Equal(len(updated.Spec.DataList)))

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
