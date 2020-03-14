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

	Context("New Cluster", func() {
		ctx := context.TODO()
		ns := SetupTest(ctx)

		It("Should create, update, and delete successfully", func() {
			key := types.NamespacedName{
				Name:      "default",
				Namespace: ns.Name,
			}

			By("Expecting secret to be created")
			created := &corev1alpha1.GeneratedSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: corev1alpha1.GeneratedSecretSpec{
					DataList: []corev1alpha1.GeneratedSecretData{{Key: "foo"}},
				},
			}
			Expect(k8sClient.Create(ctx, created)).Should(Succeed())
			Eventually(func() *corev1.Secret {
				s := &corev1.Secret{}
				_ = k8sClient.Get(ctx, key, s)
				return s
			}, timeout, interval).Should(SatisfyAll(
				WithTransform(func(e *corev1.Secret) string { return e.Name }, Equal(key.Name)),
				WithTransform(func(e *corev1.Secret) string { return e.Namespace }, Equal(key.Namespace)),
				WithTransform(func(e *corev1.Secret) map[string]string { return e.Labels }, Equal(created.Spec.SecretMeta.Labels)),
				WithTransform(func(e *corev1.Secret) map[string]string { return e.Annotations }, Equal(created.Spec.SecretMeta.Annotations)),
				WithTransform(func(e *corev1.Secret) int { return len(e.Data) }, Equal(len(created.Spec.DataList))),
			))

			By("Expecting secret to be updated")
			updated := &corev1alpha1.GeneratedSecret{}
			Expect(k8sClient.Get(ctx, key, updated)).Should(Succeed())
			updated.Spec.SecretMeta = metav1.ObjectMeta{
				Name:        "override-name",
				Namespace:   "override-namespace",
				Labels:      map[string]string{"test": "label"},
				Annotations: map[string]string{"test": "annotations"},
			}
			updated.Spec.DataList = append(updated.Spec.DataList, corev1alpha1.GeneratedSecretData{Key: "bar"})
			Expect(k8sClient.Update(ctx, updated)).Should(Succeed())

			Eventually(func() *corev1.Secret {
				s := &corev1.Secret{}
				_ = k8sClient.Get(ctx, key, s)
				return s
			}, timeout, interval).Should(SatisfyAll(
				WithTransform(func(e *corev1.Secret) string { return e.Name }, Equal(key.Name)),
				WithTransform(func(e *corev1.Secret) string { return e.Namespace }, Equal(key.Namespace)),
				WithTransform(func(e *corev1.Secret) map[string]string { return e.Labels }, Equal(updated.Spec.SecretMeta.Labels)),
				WithTransform(func(e *corev1.Secret) map[string]string { return e.Annotations }, Equal(updated.Spec.SecretMeta.Annotations)),
				WithTransform(func(e *corev1.Secret) int { return len(e.Data) }, Equal(len(updated.Spec.DataList))),
			))

			By("Expecting to delete successfully")
			Eventually(func() error {
				f := &corev1alpha1.GeneratedSecret{}
				k8sClient.Get(ctx, key, f)
				return k8sClient.Delete(ctx, f)
			}, timeout, interval).Should(Succeed())
			Eventually(func() error {
				f := &corev1alpha1.GeneratedSecret{}
				return k8sClient.Get(ctx, key, f)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})

	Context("Cluster with existing secret", func() {
		ctx := context.TODO()
		ns := SetupTest(ctx)

		It("Should not override existing secrets", func() {
			key := types.NamespacedName{
				Name:      "default",
				Namespace: ns.Name,
			}

			By("Adding secret that will conflict")
			existingSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
			}
			Expect(k8sClient.Create(ctx, existingSecret)).Should(Succeed())

			By("Creating generated secret that conflicts")
			generatedSecret := &corev1alpha1.GeneratedSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: corev1alpha1.GeneratedSecretSpec{
					DataList: []corev1alpha1.GeneratedSecretData{{Key: "foo"}},
				},
			}
			Expect(k8sClient.Create(context.Background(), generatedSecret)).Should(Succeed())

			By("Making sure secret is not overriden")
			Eventually(func() string {
				gs := &corev1alpha1.GeneratedSecret{}
				_ = k8sClient.Get(ctx, key, gs)
				return gs.Status.State
			}, timeout, interval).Should(Equal("Conflict"))
			s := &corev1.Secret{}
			Expect(k8sClient.Get(ctx, key, s)).Should(Succeed())
			Expect(len(s.Data)).To(Equal(0))

			By("Deleting existing secret")
			s = &corev1.Secret{}
			Expect(k8sClient.Get(ctx, key, s)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, s)).Should(Succeed())
			Eventually(func() string {
				gs := &corev1alpha1.GeneratedSecret{}
				_ = k8sClient.Get(ctx, key, gs)
				return gs.Status.State
			}, timeout, interval).Should(Equal("Generated"))
		})
	})
})
