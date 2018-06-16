package yaml_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/config/confsources/file/yaml"
)

func TestConfsources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Yaml config sources Suite")
}

var _ = Describe("FileYaml", func() {

	It("should parse correctly a Yaml file", func() {

		testContent := `
my:
    property: "MyValue"
    number: 42
    float: 1.61803398
    array:
        - one
        - two
        - 3
    substruct:
        toto: tata
       `

		env := ParseYamlToEnv(testContent)

		Expect(env).Should(HaveLen(7))
		Expect(env).Should(HaveKeyWithValue("my.property", "MyValue"))
		Expect(env).Should(HaveKeyWithValue("my.number", "42"))
		Expect(env).Should(HaveKeyWithValue("my.float", "1.61803398"))
		Expect(env).Should(HaveKeyWithValue("my.array.0", "one"))
		Expect(env).Should(HaveKeyWithValue("my.array.1", "two"))
		Expect(env).Should(HaveKeyWithValue("my.array.2", "3"))
		Expect(env).Should(HaveKeyWithValue("my.substruct.toto", "tata"))

	})

})
