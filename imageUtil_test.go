package goImageBrightness_test

import (
	. "github.com/rhutzel/goImageBrightness"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImageUtil", func() {
	Context("Successful file", func() {
		It("should determine a mid-level brightness", func() {
			img, _, err := ImageFromFile("./testImage.png")
			if err != nil {
				Fail("Failed to decode image.")
			}

			avgBrightness := AnalyseImage(img)
			Expect(err).To(BeNil())
			Expect(avgBrightness).To(BeNumerically(">", 25))
			Expect(avgBrightness).To(BeNumerically("<", 40))
		})
	})

	Context("Missing file", func() {
		It("should handle trying to access a missing file", func() {
			_, _, err := ImageFromFile("MISSING_FILENAME")
			Expect(err).NotTo(BeNil())
		})
	})
})
