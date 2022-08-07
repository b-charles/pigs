package config_test

import (
	. "github.com/b-charles/pigs/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func parse(args []string) map[string]string {
	if env, err := ParseArgs(args); err != nil {
		panic(err)
	} else {
		return env
	}
}

var _ = Describe("Args", func() {

	It("should parse simple value", func() {
		source := parse([]string{"--oasis=Wonderwall"})
		Expect(source).To(HaveKeyWithValue("oasis", "Wonderwall"))
	})

	It("should parse simple value with one dash", func() {
		source := parse([]string{"-moriarty=Jimmy"})
		Expect(source).To(HaveKeyWithValue("moriarty", "Jimmy"))
	})

	It("should parse value between simple quote", func() {
		source := parse([]string{"--jamiroquai='Virtual Insanity'"})
		Expect(source).To(HaveKeyWithValue("jamiroquai", "Virtual Insanity"))
	})

	It("should parse value between double quote", func() {
		source := parse([]string{"--santana=\"Flor D'Luna\""})
		Expect(source).To(HaveKeyWithValue("santana", "Flor D'Luna"))
	})

	It("should parse boolean", func() {
		source := parse([]string{"--yes"})
		Expect(source).To(HaveKeyWithValue("yes", "true"))
	})

	It("should parse false boolean", func() {
		source := parse([]string{"--no-yes"})
		Expect(source).To(HaveKeyWithValue("yes", "false"))
	})

	It("should returns an error for unknown pattern", func() {
		Expect(func() {
			parse([]string{"hello=goodbye"})
		}).To(Panic())
	})

})
