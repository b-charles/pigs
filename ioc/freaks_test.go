package ioc_test

import (
	. "github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IOC freaks show", func() {

	var (
		container *Container
	)

	BeforeEach(func() {
		container = NewContainer()
	})

	It("should handle not struct component", func() {

		Expect(container.Put(NewMap("Hello", "World"))).Should(BeNil())

		var myMap MyMap
		Expect(container.CallInjected(func(injected MyMap) {
			myMap = injected
		})).Should(Succeed())

		Expect(myMap).Should(HaveLen(1))
		Expect(myMap).Should(HaveKeyWithValue("Hello", "World"))

	})

	It("should inject not struct components", func() {

		Expect(container.PutNamed(NewMap("Hello", "World"), "MyMap")).Should(BeNil())
		Expect(container.PutNamedFactory(CreateWrapMap, "Wrap")).Should(BeNil())

		var wrap MyMap
		Expect(container.CallInjected(func(injected struct {
			Map MyMap `inject:"Wrap"`
		}) {
			wrap = injected.Map
		})).Should(Succeed())

		Expect(wrap).Should(HaveLen(1))
		Expect(wrap).Should(HaveKeyWithValue("Hello", "World"))

	})

	It("should inject not struct components in a slice", func() {

		Expect(container.PutNamed(NewMap("Hello", "World"), "MyMap1", "Maps")).Should(BeNil())
		Expect(container.PutNamed(NewMap("Hi", "Everybody"), "MyMap2", "Maps")).Should(BeNil())

		Expect(container.PutFactory(CreateSuperMap)).Should(BeNil())

		var super MySuperMap
		Expect(container.CallInjected(func(injected MySuperMap) {
			super = injected
		})).Should(Succeed())

		Expect(super).Should(HaveLen(2))
		Expect(super).Should(HaveKeyWithValue("Hello", "World"))
		Expect(super).Should(HaveKeyWithValue("Hi", "Everybody"))

	})

	It("should inject not struct components in a map", func() {

		Expect(container.PutNamed(NewMap("Hello", "World"), "MyMap1", "Maps")).Should(BeNil())
		Expect(container.PutNamed(NewMap("Hi", "Everybody"), "MyMap2", "Maps")).Should(BeNil())

		Expect(container.PutFactory(CreateHyperMap)).Should(BeNil())

		var hyper MyHyperMap
		Expect(container.CallInjected(func(injected MyHyperMap) {
			hyper = injected
		})).Should(Succeed())

		Expect(hyper).Should(HaveLen(2))
		Expect(hyper).Should(HaveKeyWithValue("MyMap1.Hello", "World"))
		Expect(hyper).Should(HaveKeyWithValue("MyMap2.Hi", "Everybody"))

	})

})
