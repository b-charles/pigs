package envvar_test

import (
	"testing"

	. "github.com/b-charles/pigs/config/confsources/envvar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfsources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Env var config sources Suite")
}

var _ = Describe("Envvar config source", func() {

	It("should read env vars", func() {

		env := ParseEnvVar([]string{"hello=it's me"})

		Expect(env).Should(HaveKeyWithValue("hello", "it's me"))

	})

	It("should convert keys", func() {

		env := ParseEnvVar([]string{"ONE_TWO_THREE_FOUR=FIVE_SIX_SEVEN_EIGHT"})

		Expect(env).Should(HaveKeyWithValue("one.two.three.four", "FIVE_SIX_SEVEN_EIGHT"))

	})

})
