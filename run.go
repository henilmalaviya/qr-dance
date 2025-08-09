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
		// No logger configured yet, so use standard fmt
		return fmt.Errorf("options cannot be nil")
	}
	logger.Info("Starting QR dance with options: %+v", *opts)

	logger.Info("Generating QR code bitmap for data: %s", opts.Data)
	bitMap, err := engine.GenerateQRCodeBitmap(opts.Data)
	if err != nil {
		logger.Error("Error generating QR code bitmap: %v", err)
		return fmt.Errorf("error generating QR code bitmap: %w", err)
	}
	logger.Debug("QR code bitmap generated successfully, size: %dx%d", len(bitMap), len(bitMap[0]))

	logger.Info("Creating new game from bitmap")
	g := engine.NewGameFromBitmap(bitMap)

	totalTicks := int(math.Ceil(opts.Duration * float64(opts.FrameRate)))
	if totalTicks <= 0 {
		logger.Error("Total ticks must be greater than 0, got %d", totalTicks)
		return fmt.Errorf("total ticks must be greater than 0, got %d", totalTicks)
	}
	logger.Info("Total ticks calculated: %d", totalTicks)

	logger.Info("Running simulation for %d ticks", totalTicks)
	images := engine.RunTicks(g, totalTicks, opts.ScaleFactor)
	logger.Debug("Simulation finished, %d images generated", len(images))

	frameDelay := int(math.Round(1000.0 / float64(opts.FrameRate)))
	if frameDelay <= 0 {
		logger.Error("Frame delay must be greater than 0, got %d", frameDelay)
		return fmt.Errorf("frame delay must be greater than 0, got %d", frameDelay)
	}
	logger.Info("Frame delay calculated: %d", frameDelay)

	logger.Info("Preparing GIF from %d images", len(images))
	gifImg, err := engine.PrepareGIFFromImages(images, frameDelay, opts.InitialFrameDelay)
	if err != nil {
		logger.Error("Error preparing GIF from images: %v", err)
		return fmt.Errorf("error preparing GIF from images: %w", err)
	}

	if opts.Base64Stdout {
		logger.Info("Encoding GIF to base64")
		base64, err := engine.GIFToBase64(gifImg)
		if err != nil {
			logger.Error("Error encoding GIF to base64: %v", err)
			return fmt.Errorf("error encoding GIF to base64: %w", err)
		}
		// Do not log the base64 content; just write to stdout
		fmt.Print(base64)
		logger.Info("Base64 encoded GIF written to stdout")
	} else {
		logger.Info("Writing GIF to file: %s", opts.OutputPath)
		err := engine.WriteGIFToFile(gifImg, opts.OutputPath)
		if err != nil {
			logger.Error("Error writing GIF to file: %v", err)
			return fmt.Errorf("error writing GIF to file: %w", err)
		}
		logger.Info("GIF written to file successfully")
	}

	return nil
}

func Run(args []string) error {
	logger.Info("Parsing command line arguments")
	opts, err := cli.ParseArgs(args)
	if err != nil {
		logger.Error("Error parsing arguments: %v", err)
		return fmt.Errorf("error parsing arguments: %w", err)
	}

	return RunWithOptions(opts)
}
