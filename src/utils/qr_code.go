package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func ByteArrayToImage(imageData []byte) (image.Image, error) {
	reader := bytes.NewReader(imageData)

	decode, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}

	return decode, nil
}

type QRCode struct {
	Base64 string
	Pin    string
}

var cachedQRCode QRCode

func CreateQRCode(pin *string) QRCode {
	if cachedQRCode.Pin == *pin && cachedQRCode.Base64 != "" {
		return cachedQRCode
	}

	qrc, err := qrcode.NewWith(
		fmt.Sprintf("https://%s/login?code=%s", Domain, *pin),
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return QRCode{}
	}

	//ttImg, err := ByteArrayToImage(assets.TTLogo)
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to convert byte array to image.")
	//}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(
		wr,
		standard.WithBgTransparent(),
		standard.WithQRWidth(10),
		standard.WithBorderWidth(3),
		//standard.WithLogoImage(ttImg),
		//standard.WithLogoSizeMultiplier(0),
	)
	err = qrc.Save(w2)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save QR code.")
	}

	cachedQRCode = QRCode{
		Base64: base64.StdEncoding.EncodeToString(buf.Bytes()),
		Pin:    *pin,
	}
	return cachedQRCode
}
