package goImageBrightness_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGoImageBrightness(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoImageBrightness Suite")
}
