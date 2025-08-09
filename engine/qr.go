package engine

import (
	"github.com/henilmalaviya/qr-dance/util/logger"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeBitmap(data string) ([][]bool, error) {
	logger.Trace("Entering GenerateQRCodeBitmap")
	defer logger.Trace("Exiting GenerateQRCodeBitmap")

	q, err := qrcode.New(data, qrcode.Medium)

	if err != nil {
		logger.Error("Failed to generate QR code: %v", err)
		return nil, err
	}

	return q.Bitmap(), nil
}
