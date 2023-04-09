package ioc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIoc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ioc Suite")
}
