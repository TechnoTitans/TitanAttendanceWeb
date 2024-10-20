package utils

import (
	"TitanAttendance/src/assets"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"image"
	"image/png"
	"io"
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

var qrCodeB64 string

func CreateQRCode(pin *string) string {
	if qrCodeB64 != "" {
		return qrCodeB64
	}

	qrc, err := qrcode.NewWith(
		GetDomain()+"/pin?code="+*pin,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return ""
	}

	ttImg, err := ByteArrayToImage(assets.TTLogo)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert byte array to image.")
	}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(
		wr,
		standard.WithQRWidth(10),
		standard.WithBorderWidth(3),
		standard.WithBgTransparent(),
		standard.WithLogoImage(ttImg),
		standard.WithLogoSizeMultiplier(0),
	)
	err = qrc.Save(w2)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save QR code.")
	}

	qrCodeB64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrCodeB64
}
