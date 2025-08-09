package cli

import (
	"github.com/akamensky/argparse"
	"github.com/henilmalaviya/qr-dance/util/logger"
)

type Options struct {
	Input             string
	Duration          float64
	FrameDelay        int
	OutputPath        string
	Base64Stdout      bool
	Base64Input       bool
	ScaleFactor       int
	VerboseLevel      int
	InitialFrameDelay int
}

func ParseArgs(args []string) (*Options, error) {
	parser := argparse.NewParser("qr-dance", "A CLI tool to generate Game of Life GIF from QR Code")

	input := parser.String("i", "input", &argparse.Options{
		Required: true,
		Help:     "Path to the input image file (PNG) or base64 string if --input-base64 is used",
	})

	duration := parser.Float("d", "duration", &argparse.Options{
		Required: false,
		Default:  3.0,
		Help:     "Duration of the animation in seconds",
	})

	frameDelay := parser.Int("f", "frame-delay", &argparse.Options{
		Required: false,
		Default:  100, // Default frame delay in milliseconds
		Help:     "Delay between frames in milliseconds",
	})

	outputPath := parser.String("o", "output", &argparse.Options{
		Required: false,
		Default:  "output.gif",
		Help:     "Path to the output GIF file",
	})

	base64Stdout := parser.Flag("b", "base64", &argparse.Options{
		Required: false,
		Help:     "Output GIF as base64 to stdout instead of writing to file",
	})

	base64Input := parser.Flag("", "input-base64", &argparse.Options{
		Required: false,
		Help:     "Treat input as base64 encoded PNG data instead of file path",
	})

	scaleFactor := parser.Int("s", "scale", &argparse.Options{
		Required: false,
		Default:  20,
		Help:     "Scale factor for the output GIF (default is 20)",
	})

	verboseLevel := parser.FlagCounter("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Increase verbosity level (can be used multiple times for more verbosity)",
	})

	initialFrameDelay := parser.Int("", "initial-frame-delay", &argparse.Options{
		Required: false,
		Default:  0,
		Help:     "Initial frame delay in milliseconds for the first frame",
	})

	if err := parser.Parse(args); err != nil {
		return nil, err
	}

	logger.SetLevel(*verboseLevel)
	logger.Info("Verbosity set to %d", *verboseLevel)
	logger.Info("CLI args parsed: input=%s duration=%.2fs frameDelay=%dms scale=%d base64Stdout=%t base64Input=%t output=%s", *input, *duration, *frameDelay, *scaleFactor, *base64Stdout, *base64Input, *outputPath)

	return &Options{
		Input:             *input,
		Duration:          *duration,
		FrameDelay:        *frameDelay,
		OutputPath:        *outputPath,
		Base64Stdout:      *base64Stdout,
		Base64Input:       *base64Input,
		ScaleFactor:       *scaleFactor,
		VerboseLevel:      *verboseLevel,
		InitialFrameDelay: *initialFrameDelay,
	}, nil
}
