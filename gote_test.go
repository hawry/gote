package main_test

import (
	_ "github.com/hawry/gote"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gote", func() {
	Describe("Testing the testing framework", func() {
		Context("And with a context", func() {
			It("Should compare strings", func() {
				Expect("hello").To(Equal("hello"))
			})
		})
	})
})
