package main

import (
	"fmt"
	"math"

	"github.com/henilmalaviya/qr-dance/cli"
	"github.com/henilmalaviya/qr-dance/engine"
	"github.com/henilmalaviya/qr-dance/util/logger"
)

func RunWithOptions(opts *cli.Options) error {
	if opts == nil {
		return fmt.Errorf("options cannot be nil")
	}

	logger.Info("Starting run with options: duration=%.2fs frameDelay=%dms scale=%d base64Stdout=%t base64Input=%t output=%s", opts.Duration, opts.FrameDelay, opts.ScaleFactor, opts.Base64Stdout, opts.Base64Input, opts.OutputPath)

	src, err := engine.ReadInput(opts.Input, opts.Base64Input)
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	logger.Debug("Input image loaded: bounds=%v", (*src).Bounds())

	srcRGBA := engine.DrawImageToRGBA(*src)
	logger.Debug("Converted input image to RGBA format")

	unitRGBA := engine.ExtractUnitQR(srcRGBA)
	logger.Debug("Detected QR region")

	matrix := engine.GenerateQRMatrix(unitRGBA)
	logger.Debug("Generated matrix: height=%d width=%d", len(matrix), len(matrix[0]))

	g := engine.NewGameFromQRMatrix(matrix)
	logger.Info("Game initialized: grid=%dx%d", len(matrix[0]), len(matrix))

	totalTicks := int(math.Round(opts.Duration * 1000 / float64(opts.FrameDelay)))
	logger.Info("Running ticks: total=%d", totalTicks)

	images := engine.RunTicks(g, totalTicks, opts.ScaleFactor)
	logger.Info("Ticks completed: frames=%d", len(images))

	gifImg, err := engine.PrepareGIFFromImages(images, opts.FrameDelay, opts.InitialFrameDelay)
	if err != nil {
		return fmt.Errorf("error preparing GIF from images: %w", err)
	}
	logger.Info("GIF prepared: frames=%d delay=%d", len(gifImg.Image), opts.FrameDelay)

	if opts.Base64Stdout {
		logger.Info("Encoding GIF to base64 and writing to stdout")
		base64, err := engine.GIFToBase64(gifImg)
		if err != nil {
			return fmt.Errorf("error encoding GIF to base64: %w", err)
		}
		// Do not log the base64 content; just write to stdout
		fmt.Print(base64)
	} else {
		logger.Info("Writing GIF to file: %s", opts.OutputPath)
		err := engine.WriteGIFToFile(gifImg, opts.OutputPath)
		if err != nil {
			return fmt.Errorf("error writing GIF to file: %w", err)
		}
	}

	logger.Info("Run completed successfully")
	return nil
}

func Run(args []string) error {
	opts, err := cli.ParseArgs(args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	logger.Debug("Parsed CLI options: %+v", opts)
	return RunWithOptions(opts)
}
