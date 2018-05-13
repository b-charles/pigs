package confsources_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfsources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Sources Suite")
}
