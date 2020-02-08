package crypto

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	v1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

var _ = Describe("Crypto", func() {
	Describe("Generate ASCII", func() {
		It("fail when no options are given", func() {
			_, err := GenerateRandomASCIIString(100, []v1alpha1.ValueOption{})
			Expect(err).Should(HaveOccurred())
		})

		It("generate string with expected length", func() {
			s, _ := GenerateRandomASCIIString(100, []v1alpha1.ValueOption{v1alpha1.Uppercase})
			Expect(len(s)).To(Equal(100))
		})
	})

	DescribeTable("Match Options",
		func(c string, os []v1alpha1.ValueOption, expected bool) {
			Expect(matchesOptions(c, os)).To(Equal(expected))
		},
		Entry("Uppercase true", "A", []v1alpha1.ValueOption{v1alpha1.Uppercase}, true),
		Entry("Uppercase false", "a", []v1alpha1.ValueOption{v1alpha1.Uppercase}, false),
		Entry("Lowercase true", "a", []v1alpha1.ValueOption{v1alpha1.Lowercase}, true),
		Entry("Lowercase false", "A", []v1alpha1.ValueOption{v1alpha1.Lowercase}, false),
		Entry("Numbers true", "0", []v1alpha1.ValueOption{v1alpha1.Numbers}, true),
		Entry("Numbers false", "z", []v1alpha1.ValueOption{v1alpha1.Numbers}, false),
		Entry("Symbols true", "@", []v1alpha1.ValueOption{v1alpha1.Symbols}, true),
		Entry("Symbols false", "0", []v1alpha1.ValueOption{v1alpha1.Symbols}, false),
		Entry("All", "$", []v1alpha1.ValueOption{v1alpha1.Uppercase, v1alpha1.Lowercase, v1alpha1.Numbers, v1alpha1.Symbols}, true),
	)
})
