package casts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCastsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Casts Test Suite")
}