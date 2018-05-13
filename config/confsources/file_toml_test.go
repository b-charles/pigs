package confsources_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/l3eegbee/pigs/config/confsources"
)

var _ = Describe("FileToml", func() {

	It("should parse correctly a Toml file", func() {

		testContent := `

		[my]
		property = "MyValue"
		number = 42
		float = 1.61803398
		array = [ "one", "two", "3" ]

		[my.substruct]
		toto = "tata"

		`

		env := ParseTomlToEnv(testContent)

		Expect(env).Should(HaveLen(7))
		Expect(env).Should(HaveKeyWithValue("my.property", "MyValue"))
		Expect(env).Should(HaveKeyWithValue("my.number", "42"))
		Expect(env).Should(HaveKeyWithValue("my.float", "1.61803398"))
		Expect(env).Should(HaveKeyWithValue("my.array[0]", "one"))
		Expect(env).Should(HaveKeyWithValue("my.array[1]", "two"))
		Expect(env).Should(HaveKeyWithValue("my.array[2]", "3"))
		Expect(env).Should(HaveKeyWithValue("my.substruct.toto", "tata"))

	})

})
