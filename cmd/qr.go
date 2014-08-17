package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"strings"

	"github.com/crowdmob/goamz/s3"
	"github.com/qpliu/qrencode-go/qrencode"
)

func init() {
	AddPlugin("QRCode", "(?i)^\\.qr(code)?$", MessageHandler(CreateQR), false, false)
}

func CreateQR(msg *Message) {
	text := strings.Join(msg.Params[1:], " ")
	grid, err := qrencode.Encode(text, qrencode.ECLevelQ)
	if err != nil {
		msg.Return("Error encoding ur text!")
		return
	}
	var b bytes.Buffer
	png.Encode(&b, grid.Image(8))
	data, err := ioutil.ReadAll(&b)
	if err != nil {
		msg.Return("Error reading from my PNG buffer, this should not happen!")
		return
	}
	name := GenerateRandom(6)
	uploaded, err := UploadToS3(name, data, "image/png", s3.PublicRead)
	if err != nil {
		msg.Return("Error Uploading QR Code!")
		return
	}
	out := fmt.Sprintf("%s: %s", bold("Your QR"), uploaded)
	msg.Return(out)
}
