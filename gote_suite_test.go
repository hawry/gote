package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gote Suite")
}
