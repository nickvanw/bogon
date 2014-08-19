package dh1080

import (
	"bytes"
	"crypto/cipher"
	"errors"

	"code.google.com/p/go.crypto/blowfish"
)

func blowfishChecksizeAndPad(pt []byte) []byte {
	modulus := len(pt) % blowfish.BlockSize
	if modulus != 0 {
		padlen := blowfish.BlockSize - modulus
		for i := 0; i < padlen; i++ {
			pt = append(pt, 0)
		}
	}
	return pt
}
func Dec(et, key []byte) ([]byte, error) {
	dcipher, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	div := et[:blowfish.BlockSize]
	decrypted := et[blowfish.BlockSize:]
	if len(decrypted)%blowfish.BlockSize != 0 {
		return nil, errors.New("decrypted is not a multiple of blowfish.BlockSize")
	}
	dcbc := cipher.NewCBCDecrypter(dcipher, div)
	dcbc.CryptBlocks(decrypted, decrypted)
	return bytes.Trim(decrypted, "\x00"), nil
}
func Enc(ppt, key []byte) ([]byte, error) {
	ppt = blowfishChecksizeAndPad(ppt)
	ecipher, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, blowfish.BlockSize+len(ppt))
	eiv := ciphertext[:blowfish.BlockSize]
	ecbc := cipher.NewCBCEncrypter(ecipher, eiv)
	ecbc.CryptBlocks(ciphertext[blowfish.BlockSize:], ppt)
	return ciphertext, nil
}
