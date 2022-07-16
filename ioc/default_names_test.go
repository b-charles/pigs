package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Tassadar struct{}
type Zeratul struct{}

var _ = Describe("Default names", func() {

	It("should name a simple component", func() {

		name := DefaultComponentName(&Tassadar{})
		Expect(name).To(Equal("github.com/b-charles/pigs/ioc_test/Tassadar"))

	})

	It("should get name from a factory", func() {

		name := DefaultFactoryName(func() *Tassadar { return nil })
		Expect(name).To(Equal("github.com/b-charles/pigs/ioc_test/Tassadar"))

	})

	It("should get aliases from a function", func() {

		aliases := DefaultAliases(func(*Zeratul, *Tassadar) {})
		Expect(aliases).To(HaveLen(2))
		Expect(aliases).To(ContainElement("github.com/b-charles/pigs/ioc_test/Zeratul"))
		Expect(aliases).To(ContainElement("github.com/b-charles/pigs/ioc_test/Tassadar"))

	})

})
