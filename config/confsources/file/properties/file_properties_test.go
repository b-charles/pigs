package properties_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/config/confsources/file/properties"
)

func TestConfsources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Properties config sources Suite")
}

var _ = Describe("FileProperties", func() {

	It("should parse a Properties file without replace anything", func() {

		testContent := `

# This is a comment
! This is an other one

my.property=MyValue
who.do.that Nobody, I guess
psycho:Only psychopaths use :

my.replaced.property=${my.property}

		`

		env := ParsePropertiesToEnv(testContent)

		Expect(env).Should(HaveLen(4))
		Expect(env).Should(HaveKeyWithValue("my.property", "MyValue"))
		Expect(env).Should(HaveKeyWithValue("who.do.that", "Nobody, I guess"))
		Expect(env).Should(HaveKeyWithValue("psycho", "Only psychopaths use :"))
		Expect(env).Should(HaveKeyWithValue("my.replaced.property", "${my.property}"))

	})

})
