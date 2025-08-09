package cli

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/henilmalaviya/qr-dance/util/logger"
)

type Options struct {
	Data              string
	Duration          float64
	FrameRate         int
	OutputPath        string
	Base64Stdout      bool
	ScaleFactor       int
	VerboseLevel      int
	InitialFrameDelay int
}

func ParseArgs(args []string) (*Options, error) {
	parser := argparse.NewParser("qr-dance", "A CLI tool to generate Game of Life GIF inside a QR code")

	data := parser.StringPositional(&argparse.Options{
		Required: true,
		Help:     "Data to be encoded in the QR code",
	})

	duration := parser.Float("d", "duration", &argparse.Options{
		Required: false,
		Default:  3.0,
		Help:     "Duration of the animation in seconds",
	})

	frameRate := parser.Int("f", "frame-rate", &argparse.Options{
		Required: false,
		Default:  10,
		Help:     "Frame rate of the animation in frames per second (default is 10)",
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
		// print help
		fmt.Print(parser.Usage(""))

		return nil, err
	}

	logger.SetLevel(*verboseLevel)

	return &Options{
		Data:              *data,
		Duration:          *duration,
		FrameRate:         *frameRate,
		OutputPath:        *outputPath,
		Base64Stdout:      *base64Stdout,
		ScaleFactor:       *scaleFactor,
		VerboseLevel:      *verboseLevel,
		InitialFrameDelay: *initialFrameDelay,
	}, nil
}
