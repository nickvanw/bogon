package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"github.com/dchest/uniuri"
)

const bucket = "i.nvw.io"
const bucketpath = "http://i.nvw.io/"

func getSite(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func bold(s string) string {
	return fmt.Sprintf("\x02%s\x02", s)
}
func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func UploadToS3(name string, data []byte, datatype string, permission s3.ACL) (string, error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return "", errors.New("Couldn't auth to Amazon AWS")
	}
	s := s3.New(auth, aws.USWest)
	bucket := s.Bucket(bucket)
	err = bucket.Put(name, data, datatype, permission, s3.Options{})
	if err != nil {
		return "", errors.New("Couldn't Upload That!")
	}
	return fmt.Sprintf("%s%s", bucketpath, name), nil
}

func GenerateRandom(n int) string {
	return uniuri.NewLen(n)
}

func urlencode(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%%%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%%%X%%%X%%%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%%%X%%%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}
