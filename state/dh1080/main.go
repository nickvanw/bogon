package dh1080

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

const DH1080_STR = "++ECLiPSE+is+proud+to+present+latest+FiSH+release+featuring+even+more+security+for+you+++shouts+go+out+to+TMG+for+helping+to+generate+this+cool+sophie+germain+prime+number++++/C32L"
const ENCODING = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

type DH1080 struct {
	p       *big.Int
	g       *big.Int
	q       *big.Int
	Private *big.Int
	Public  *big.Int
	Secret  *big.Int
	State   int
}

func DH1080_Init() *DH1080 {
	data, _ := DH1080_Base64Decode(DH1080_STR)
	DH1080Ctx := &DH1080{}
	DH1080Ctx.p = new(big.Int).SetBytes(data)
	DH1080Ctx.g = new(big.Int).SetInt64(2)
	sub := new(big.Int).Sub(DH1080Ctx.p, big.NewInt(1))
	DH1080Ctx.q = sub.Div(sub, big.NewInt(2))
	DH1080Ctx.State = 0
	DH1080Ctx.GeneratePrivateKey()
	return DH1080Ctx
}

func (dh *DH1080) GeneratePrivateKey() error {
	err := dh.generate()
	if err != nil {
		return err
	}
	for !dh.Validate() {
		err := dh.generate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (dh *DH1080) generate() error {
	data := make([]byte, 1080/8)
	_, err := rand.Read(data)
	if err != nil {
		return err
	}
	dh.Private = new(big.Int).SetBytes(data)
	dh.Public = new(big.Int).Exp(dh.g, dh.Private, dh.p)
	return nil
}

func (dh *DH1080) ComputeKey(otherKey []byte) error {
	newKey := new(big.Int).SetBytes(otherKey)
	dh.Secret = new(big.Int).Exp(newKey, dh.Private, dh.p)
	return nil
}

func (dh *DH1080) Validate() bool {
	return "1" == new(big.Int).Exp(dh.Public, dh.q, dh.p).String()
}

func (dh *DH1080) Pack() (string, error) {
	var cmd string
	if dh.State == 0 {
		dh.State = 1
		cmd = "DH1080_INIT "
	} else {
		cmd = "DH1080_FINISH "
	}
	data, err := DH1080_Base64Encode(dh.Public.Bytes())
	if err != nil {
		return "", err
	}
	return cmd + data, nil
}

func (dh *DH1080) GetSecret() (string, error) {
	if dh.Secret == nil {
		return "", errors.New("No secret to encode")
	}
	data := dh.Secret.Bytes()
	hasher := sha256.New()
	hasher.Write(data)
	return DH1080_Base64Encode(hasher.Sum(nil))
}

func (dh *DH1080) Unpack(msg string) error {
	if !strings.HasPrefix(msg, "DH1080_") {
		return errors.New("Invalid Message")
	}
	data := strings.Split(msg, " ")
	keyData := data[1]
	decodedKeyData, err := DH1080_Base64Decode(keyData)
	if err != nil {
		return err
	}
	err = dh.ComputeKey(decodedKeyData)
	if err != nil {
		return err
	}
	dh.State = 1
	return nil
}

func DH1080_Base64Encode(data []byte) (string, error) {
	if len(data) < 1 {
		return "", errors.New("Zero Length String")
	}
	base := base64.NewEncoding(ENCODING)
	encodedString := base.EncodeToString(data)
	if !strings.Contains(encodedString, "=") {
		encodedString += "A"
	} else {
		encodedString = strings.TrimRight(encodedString, "=")
	}
	return encodedString, nil
}

func DH1080_Base64Decode(data string) ([]byte, error) {
	base := base64.NewEncoding(ENCODING)
	if len(data)%4 == 1 && data[len(data)-1] == 'A' {
		return base.DecodeString(string(data[:len(data)-1]))
	}
	if len(data)%4 != 0 {
		data = data + strings.Repeat("=", (4-(len(data)%4)))
	}
	return base.DecodeString(data)
}
