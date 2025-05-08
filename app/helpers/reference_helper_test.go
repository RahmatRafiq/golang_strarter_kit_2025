package helpers_test

import (
	"sync"

	"golang_starter_kit_2025/app/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateReference", func() {
	Context("when reference is generated", func() {
		It("should return reference with code", func() {
			code := "INV"
			ref := helpers.GenerateReference(code)
			Expect(ref).To(HavePrefix(code + "-"))
		})
	})

	Context("when reference is generated massively is not duplicated", func() {
		It("should return unique reference", func() {
			code := "INV"
			refs := make(map[string]bool)
			for i := 0; i < 100000; i++ {
				ref := helpers.GenerateReference(code)
				Expect(refs[ref]).To(BeFalse())
				refs[ref] = true
			}
		})
	})

	Context("when reference is generated massively in parallel with optimal Big O notation", func() {
		It("should return unique reference", func() {
			code := "INV"
			refs := sync.Map{}
			wg := sync.WaitGroup{}
			for i := 0; i < 100000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					ref := helpers.GenerateReference(code)
					if _, loaded := refs.LoadOrStore(ref, true); loaded {
						Fail("duplicate reference found: " + ref)
					}
				}()
			}
			wg.Wait()
		})
	})

	Context("when code contains special characters", func() {
		It("should return reference with special characters in code", func() {
			code := "INV-123_ABC"
			ref := helpers.GenerateReference(code)
			Expect(ref).To(HavePrefix(code + "-"))
		})
	})

	Context("when multiple references are generated sequentially", func() {
		It("should return unique references", func() {
			code := "SEQ"
			ref1 := helpers.GenerateReference(code)
			ref2 := helpers.GenerateReference(code)
			Expect(ref1).NotTo(Equal(ref2))
		})
	})
})
