package confsources_test

import (
	"os"

	. "github.com/l3eegbee/pigs/config/confsources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envvar config source", func() {

	It("should read env vars", func() {

		os.Setenv("hello", "it's me")

		env := NewEnvVarConfigSource().LoadEnv()

		Expect(env).Should(HaveKeyWithValue("hello", "it's me"))

	})

	It("should convert keys", func() {

		os.Setenv("ONE_TWO_THREE_FOUR", "FIVE_SIX_SEVEN_EIGHT")

		env := NewEnvVarConfigSource().LoadEnv()

		Expect(env).Should(HaveKeyWithValue("one.two.three.four", "FIVE_SIX_SEVEN_EIGHT"))

	})

})
