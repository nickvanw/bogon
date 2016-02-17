// Package dh1080 implements the FiSH encryption scheme for 1:1 communicateion over IRC
// it's based on a Diffie-Hellman key exchange and Blowfish CBC encryption, which a lot of
// implementations do not support.
// Functionality has only been tested with Textual 5
package dh1080

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

const (
	prime1080     = "++ECLiPSE+is+proud+to+present+latest+FiSH+release+featuring+even+more+security+for+you+++shouts+go+out+to+TMG+for+helping+to+generate+this+cool+sophie+germain+prime+number++++/C32L"
	validEncoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

// DH1080 implements 1:1 IRC encryption based on
// a Diffie-Hellman key exchange
type DH1080 struct {
	p       *big.Int
	g       *big.Int
	q       *big.Int
	Private *big.Int
	Public  *big.Int
	Secret  *big.Int
	State   int
}

// New begins the process of a Diffie-Hellman exchange for a key
// It also chooses a secret integer to return
func New() *DH1080 {
	data, _ := base64Decode(prime1080) // "a curious 1080 bit prime number"
	DH1080Ctx := &DH1080{}
	DH1080Ctx.p = new(big.Int).SetBytes(data) // p, from the curious number above
	DH1080Ctx.g = new(big.Int).SetInt64(2)    // g = 2, agreed on in a spec somewhere
	// p = 2q + 1, so q = (p - 1)/2
	sub := new(big.Int).Sub(DH1080Ctx.p, big.NewInt(1))
	DH1080Ctx.q = sub.Div(sub, big.NewInt(2))
	DH1080Ctx.State = 0
	DH1080Ctx.genPrivateKey()
	return DH1080Ctx
}

// GetSecret returns the base64 encoded secret
func (dh *DH1080) GetSecret() (string, error) {
	if dh.Secret == nil {
		return "", errors.New("No secret to encode")
	}
	data := dh.Secret.Bytes()
	hasher := sha256.New()
	hasher.Write(data)
	return base64Encode(hasher.Sum(nil))
}

// Pack uses the key generated when Unpack was called
func (dh *DH1080) Pack() (string, error) {
	var cmd string
	if dh.State == 0 {
		dh.State = 1
		cmd = "DH1080_INIT "
	} else {
		cmd = "DH1080_FINISH "
	}
	data, err := base64Encode(dh.Public.Bytes())
	if err != nil {
		return "", err
	}
	return cmd + data, nil
}

// Unpack accept's Alice's "g^a mod p" or "2^a mod $prime1080"
func (dh *DH1080) Unpack(msg string) error {
	if !strings.HasPrefix(msg, "DH1080_") {
		return errors.New("Invalid Message")
	}
	data := strings.Split(msg, " ")
	keyData := data[1]
	decodedKeyData, err := base64Decode(keyData)
	if err != nil {
		return err
	}
	err = dh.computeKey(decodedKeyData)
	if err != nil {
		return err
	}
	dh.State = 1
	return nil
}

// genPrivateKey calls generate until a valid key per
// is generated
func (dh *DH1080) genPrivateKey() error {
	err := dh.generate()
	if err != nil {
		return err
	}
	for !dh.validate() {
		err := dh.generate()
		if err != nil {
			return err
		}
	}
	return nil
}

// generate chooses a secret integer "b" and generates
// g^b mod p or 2^b mod $prime1080
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

// computeKey accepts the secret from Alice and creates a shared secret
// using the generated secret in generate
func (dh *DH1080) computeKey(otherKey []byte) error {
	// Alice's Secret Number
	newKey := new(big.Int).SetBytes(otherKey)
	// Secret = A^b mod P
	dh.Secret = new(big.Int).Exp(newKey, dh.Private, dh.p)
	return nil
}

func (dh *DH1080) validate() bool {
	return "1" == new(big.Int).Exp(dh.Public, dh.q, dh.p).String()
}

func base64Encode(data []byte) (string, error) {
	if len(data) < 1 {
		return "", errors.New("Zero Length String")
	}
	base := base64.NewEncoding(validEncoding)
	encodedString := base.EncodeToString(data)
	if !strings.Contains(encodedString, "=") {
		encodedString += "A"
	} else {
		encodedString = strings.TrimRight(encodedString, "=")
	}
	return encodedString, nil
}

func base64Decode(data string) ([]byte, error) {
	base := base64.NewEncoding(validEncoding)
	if len(data)%4 == 1 && data[len(data)-1] == 'A' {
		return base.DecodeString(string(data[:len(data)-1]))
	}
	if len(data)%4 != 0 {
		data = data + strings.Repeat("=", (4-(len(data)%4)))
	}
	return base.DecodeString(data)
}
