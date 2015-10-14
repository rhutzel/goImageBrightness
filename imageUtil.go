package goImageBrightness

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"sync"
)

const SCALE_MIN float64 = 0.0
const SCALE_MAX float64 = 100.0

type SegmentResult struct {
	nPixels int
	sum float64
}

func ImageFromFile(imgPath string) (image.Image, string, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, "", errors.New("Could not open [" + imgPath + "]")
	}
	defer file.Close()
	return image.Decode(bufio.NewReader(file))
}

func AnalyseImage(img image.Image) int {
	bounds := img.Bounds()
	segmentResult := sumPixels(img, bounds.Min.X, bounds.Max.X, relativeLuminanceRec709)
	return normalizeColorChannelAsPct(segmentResult.sum / float64(segmentResult.nPixels))
}

func ParallelAnalyseImage(img image.Image, splits int) int {
	bounds := img.Bounds()
	ceiling := bounds.Max.X - bounds.Min.X
	segmentLength := int(math.Ceil(float64(ceiling) / float64(splits)))
	log.Printf("Ceiling is %d and segmentLength is %d", ceiling, segmentLength)

	var xMin int
	var waiters sync.WaitGroup
	segmentResults := make(chan SegmentResult, splits)

	for ceiling > 0 {
		waiters.Add(1)
		xMin = ceiling - segmentLength
		if xMin < 0 {
			xMin = 0
		}

		go func(img image.Image, xMin int, ceiling int, relativeLuminanceRec709 func(color.Color) float64) {
			defer waiters.Done()
			log.Printf("Covering from %d up to %d...", xMin, ceiling)
			segmentResults <- sumPixels(img, xMin, ceiling, relativeLuminanceRec709)
		}(img, xMin, ceiling, relativeLuminanceRec709)

		ceiling -= segmentLength
	}

	waiters.Wait()
	close(segmentResults)

	log.Printf("All splits have completed. Reducing results...")
	nPixels := 0
	sum := 0.0

	for segmentResult := range segmentResults {
		nPixels += segmentResult.nPixels
		sum += segmentResult.sum
		log.Printf("Running total pixels [%d] with sum [%f]...", nPixels, sum)
	}

	if nPixels > 0 {
		return normalizeColorChannelAsPct(sum / float64(nPixels))
	} else {
		return 0
	}
}

func sumPixels(img image.Image, xMin int, xMax int, fn func(color.Color) float64) SegmentResult {
	nPixels := 0
	sum := 0.0
	bounds := img.Bounds()
	for x := xMin; x < xMax; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			nPixels++
			sum += fn(img.At(x, y))
		}
	}
	return SegmentResult{ nPixels: nPixels, sum: sum }
}

func normalizeColorChannelAsPct(channel float64) int {
	// alpha-premultiplied 16-bits per channel
	return int((SCALE_MIN + (channel - 0.0) * (SCALE_MAX - SCALE_MIN)) / (float64(math.MaxUint16) - 0.0))
}

func relativeLuminanceRec709(color color.Color) float64 {
	r, g, b, _ := color.RGBA()
	return 0.2126 * float64(r) + 0.7152 * float64(g) + 0.0722 * float64(b)
}
