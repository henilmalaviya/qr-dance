package engine

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/henilmalaviya/qr-dance/util/logger"
)

func isPixelBlack(r, g, b uint8) bool {
	return r < 2 && g < 2 && b < 2
}

func isColorBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	isBlack := isPixelBlack(r8, g8, b8)
	logger.Trace("isColorBlack: r=%d g=%d b=%d isBlack=%v", r8, g8, b8, isBlack)
	return isBlack
}

func DrawImageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba
}

func ExtractPureQRRegion(img image.Image) (*image.RGBA, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	logger.Debug("ExtractPureQRRegion: image=%dx%d", width, height)

	diagonal := int(math.Floor(math.Sqrt(float64(width*width + height*height))))
	diagonal = int(math.Min(float64(diagonal), float64(width)))
	diagonal = int(math.Min(float64(diagonal), float64(height)))
	logger.Debug("ExtractPureQRRegion: calculated diagonal=%d", diagonal)

	d := 0

	for ; d < diagonal; d++ {
		c := img.At(d, d)
		isBlack := isColorBlack(c)

		if isBlack {
			logger.Debug("ExtractPureQRRegion: found black pixel at d=%d", d)
			break
		}
	}

	// Create a new image to hold the extracted region

	newWidth, newHeight := diagonal-(d*2), diagonal-(d*2)

	extracted := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := range newHeight {
		for x := range newWidth {
			c := img.At(x+d, y+d)
			extracted.Set(x, y, c)
			logger.Trace("ExtractPureQRRegion: copying pixel at (%d, %d) -> (%d, %d)", x+d, y+d, x, y)
		}
	}

	return extracted, nil
}

func ExtractUnitQR(img *image.RGBA) (unitImg *image.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	logger.Debug("ExtractUnitQR: image=%dx%d", width, height)

	diagonal := int(math.Floor(math.Sqrt(math.Pow(float64(width), 2) + math.Pow(float64(height), 2))))
	diagonal = int(math.Min(math.Min(float64(diagonal), float64(width)), float64(height)))
	logger.Debug("ExtractUnitQR: calculated diagonal=%d", diagonal)

	d := 0
	for ; d < diagonal; d++ {
		c := img.At(d, d)
		isBlack := isColorBlack(c)

		logger.Trace("ExtractUnitQR: checking pixel at (%d, %d) isBlack=%v", d, d, isBlack)

		if isBlack {
			logger.Debug("ExtractUnitQR: found black pixel at d=%d", d)
			break
		}
	}

	bw := d
	bh := d
	width, height = width-(d*2), height-(d*2)

	for {
		if bw >= width || bh >= height {
			logger.Warn("ExtractUnitQR: reached diagonal limit without finding a black pixel")
			break
		}

		c := img.At(bw, bh)
		isBlack := isColorBlack(c)

		logger.Trace("ExtractUnitQR: bw=%d bh=%d isBlack=%v", bw, bh, isBlack)

		if !isBlack {
			logger.Debug("ExtractUnitQR: found non-black pixel at bw=%d bh=%d", bw, bh)
			break
		}

		bw++
		bh++
	}

	bw -= d
	bh -= d

	nx := int(math.Ceil(float64(width) / float64(bw)))
	ny := int(math.Ceil(float64(height) / float64(bh)))

	unitRect := image.Rect(0, 0, nx, ny)
	unitImg = image.NewRGBA(unitRect)
	for bx := range nx {
		for by := range ny {

			srcX := d + (bx * bw)
			srcY := d + (by * bh)

			isBlack := isColorBlack(img.At(srcX, srcY))

			logger.Trace("ExtractUnitQR: bx=%d by=%d srcX=%d srcY=%d isBlack=%v", bx, by, srcX, srcY, isBlack)

			if isBlack {
				unitImg.Set(bx, by, color.Black)
			} else {
				unitImg.Set(bx, by, color.White)
			}
		}
	}

	return unitImg
}

func GenerateQRMatrix(unitImg *image.RGBA) [][]int {
	nx, ny := unitImg.Bounds().Dx(), unitImg.Bounds().Dy()
	logger.Debug("GenerateQRMatrix: unitImg size=%dx%d", nx, ny)

	matrix := make([][]int, ny)
	for i := range matrix {
		matrix[i] = make([]int, nx)
	}

	for by := range ny {
		for bx := range nx {
			c := unitImg.At(bx, by)
			isBlack := isColorBlack(c)
			if isBlack {
				matrix[by][bx] = 1
			} else {
				matrix[by][bx] = 0
			}
			logger.Trace("GenerateQRMatrix: bx=%d by=%d isBlack=%v", bx, by, isBlack)
		}
	}

	logger.Debug("GenerateQRMatrix: matrix built size=%dx%d", len(matrix[0]), len(matrix))
	return matrix
}

func UpscalePNGImage(img *image.RGBA, scale int) *image.RGBA {
	width := img.Bounds().Dx() * scale
	height := img.Bounds().Dy() * scale
	logger.Trace("UpscalePNGImage: src=%dx%d scale=%d -> dst=%dx%d", img.Bounds().Dx(), img.Bounds().Dy(), scale, width, height)
	upscaledImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			origX := x / scale
			origY := y / scale
			upscaledImg.Set(x, y, img.At(origX, origY))
		}
	}

	return upscaledImg
}

func DrawMatrixToPNG(matrix [][]int) *image.RGBA {
	logger.Trace("DrawMatrixToPNG: size=%dx%d", len(matrix[0]), len(matrix))

	img := image.NewRGBA(image.Rect(0, 0, len(matrix[0]), len(matrix)))
	for y, row := range matrix {
		for x, cell := range row {
			if cell == 1 {
				img.Set(x, y, image.Black)
			} else {
				img.Set(x, y, image.White)
			}
		}
	}

	return img
}
