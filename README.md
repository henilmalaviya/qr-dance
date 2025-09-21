# QR Dance

A CLI tool that generates animated GIFs by running Conway's Game of Life simulation on QR code patterns. Encode any data into a QR code and watch it "dance" through cellular automaton evolution.

## Features

- Generate QR codes from any text data
- Animate QR codes using Conway's Game of Life rules
- Output animated GIFs or base64-encoded data
- Customizable animation duration, frame rate, and scaling
- Command-line interface with flexible options

## Installation

### Prerequisites

- Go 1.24.4 or later

### Build from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/henilmalaviya/qr-dance.git
   cd qr-dance
   ```

2. Build the binary:

   ```bash
   go build -o qr-dance .
   ```

3. (Optional) Move to a directory in your PATH:
   ```bash
   sudo mv qr-dance /usr/local/bin/
   ```

## Usage

```bash
qr-dance [OPTIONS] DATA
```

### Arguments

- `DATA`: The text data to encode in the QR code (required)

### Options

- `-d, --duration FLOAT`: Duration of the animation in seconds (default: 3.0)
- `-f, --frame-rate INT`: Frame rate in frames per second (default: 10)
- `-o, --output PATH`: Output GIF file path (default: output.gif)
- `-b, --base64`: Output GIF as base64 to stdout instead of file
- `-s, --scale INT`: Scale factor for the GIF (default: 20)
- `-v, --verbose`: Increase verbosity (use multiple times for more detail)
- `--initial-frame-delay INT`: Initial frame delay in milliseconds (default: 0)

## Examples

### Basic Usage

Generate a dancing QR code for "Hello World":

```bash
./qr-dance "Hello World"
```

This creates `output.gif` with a 3-second animation at 10 FPS.

### Custom Output and Settings

Generate a longer animation with higher frame rate:

```bash
./qr-dance -d 5.0 -f 15 -o hello.gif "Hello World"
```

### Base64 Output

Output the GIF as base64 to stdout:

```bash
./qr-dance -b "Hello World" > animation.b64
```

### High-Resolution Output

Increase the scale factor for larger GIFs:

```bash
./qr-dance -s 30 -o large.gif "Hello World"
```

## How It Works

1. The input data is encoded into a QR code bitmap
2. Conway's Game of Life rules are applied to the QR code pattern over time
3. Each frame of the evolution is captured as an image
4. The images are compiled into an animated GIF

The QR code remains scannable throughout the animation, as the Game of Life simulation preserves the essential pattern structure.

## Dependencies

- [github.com/akamensky/argparse](https://github.com/akamensky/argparse) - Command-line argument parsing
- [github.com/skip2/go-qrcode](https://github.com/skip2/go-qrcode) - QR code generation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
