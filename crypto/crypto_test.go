package crypto

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	v1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

var _ = Describe("Crypto", func() {
	Describe("Generate ASCII", func() {
		It("generate string with expected length", func() {
			s, _ := GenerateRandomASCIIString(100, []v1alpha1.CharacterOption{})
			Expect(len(s)).To(Equal(100))
		})
	})

	DescribeTable("Match Options",
		func(c string, os []v1alpha1.CharacterOption, expected bool) {
			Expect(matchesOptions(c, os)).To(Equal(expected))
		},
		Entry("Uppercase true", "A", []v1alpha1.CharacterOption{v1alpha1.Lowercase, v1alpha1.Numbers, v1alpha1.Symbols}, false),
		Entry("Uppercase false", "a", []v1alpha1.CharacterOption{v1alpha1.Lowercase, v1alpha1.Numbers, v1alpha1.Symbols}, true),

		Entry("Uppercase true", "a", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Numbers, v1alpha1.Symbols}, false),
		Entry("Uppercase false", "A", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Numbers, v1alpha1.Symbols}, true),

		Entry("Uppercase true", "0", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Lowercase, v1alpha1.Symbols}, false),
		Entry("Uppercase true", "z", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Lowercase, v1alpha1.Symbols}, true),

		Entry("Uppercase true", "@", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Lowercase, v1alpha1.Numbers}, false),
		Entry("Uppercase false", "0", []v1alpha1.CharacterOption{v1alpha1.Uppercase, v1alpha1.Lowercase, v1alpha1.Numbers}, true),
	)
})
