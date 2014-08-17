package cmd

import (
	"fmt"
	"strings"

	"github.com/crowdmob/goamz/s3"
)

func init() {
	AddPlugin("StreetView", "(?i)^\\.s(treet)?v(iew)?$", MessageHandler(GoogleSV), false, false)
}

const svapi = "http://maps.googleapis.com/maps/api/streetview?size=640x640&sensor=false"

func GoogleSV(msg *Message) {
	addr := strings.Join(msg.Params[1:], " ")
	googleapi, avail := GetConfig("Streetview")
	if avail != true {
		fmt.Println("Google Streetview API Key not available")
		return
	}
	url := fmt.Sprintf("%s&location=%s&key=%s", svapi, urlencode(addr), googleapi)
	data, err := getSite(url)
	if err != nil {
		fmt.Println(err)
		msg.Return("Unable to get that streetview site!")
	} else {
		name := fmt.Sprintf("%s.png", GenerateRandom(6))
		url, err = UploadToS3(name, data, "image/png", s3.PublicRead)
		if err != nil {
			msg.Return("Unable to upload!")
			return
		}
		msg.Return(fmt.Sprintf("Streetview Image: %s", url))
	}
}
