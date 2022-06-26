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

		container.Put(NewMap("Hello", "World"), "MyMap")

		var myMap MyMap
		Expect(container.CallInjected(func(injected MyMap) {
			myMap = injected
		})).Should(Succeed())

		Expect(myMap).Should(HaveLen(1))
		Expect(myMap).Should(HaveKeyWithValue("Hello", "World"))

	})

	It("should inject not struct components", func() {

		container.Put(NewMap("Hello", "World"), "MyMap")
		container.PutFactory(CreateWrapMap, "Wrap")

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

		container.Put(NewMap("Hello", "World"), "MyMap1", "Maps")
		container.Put(NewMap("Hi", "Everybody"), "MyMap2", "Maps")

		container.PutFactory(CreateSuperMap, "MySuperMap")

		var super MySuperMap
		Expect(container.CallInjected(func(injected MySuperMap) {
			super = injected
		})).Should(Succeed())

		Expect(super).Should(HaveLen(2))
		Expect(super).Should(HaveKeyWithValue("Hello", "World"))
		Expect(super).Should(HaveKeyWithValue("Hi", "Everybody"))

	})

	It("should inject not struct components in a map", func() {

		container.Put(NewMap("Hello", "World"), "MyMap1", "Maps")
		container.Put(NewMap("Hi", "Everybody"), "MyMap2", "Maps")

		container.PutFactory(CreateHyperMap, "MyHyperMap")

		var hyper MyHyperMap
		Expect(container.CallInjected(func(injected MyHyperMap) {
			hyper = injected
		})).Should(Succeed())

		Expect(hyper).Should(HaveLen(2))
		Expect(hyper).Should(HaveKeyWithValue("MyMap1.Hello", "World"))
		Expect(hyper).Should(HaveKeyWithValue("MyMap2.Hi", "Everybody"))

	})

})
