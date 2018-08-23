package bitmex

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

const (
	HashSHA1 = iota
	HashSHA256
	HashSHA512
	HashSHA512_384
	HashMD5
)

// GetHMAC returns a keyed-hash message authentication code using the desired
// hashtype
func GetHMAC(hashType int, input, key []byte) []byte {
	var hash func() hash.Hash

	switch hashType {
	case HashSHA1:
		{
			hash = sha1.New
		}
	case HashSHA256:
		{
			hash = sha256.New
		}
	case HashSHA512:
		{
			hash = sha512.New
		}
	case HashSHA512_384:
		{
			hash = sha512.New384
		}
	case HashMD5:
		{
			hash = md5.New
		}
	}

	hmac := hmac.New(hash, []byte(key))
	hmac.Write(input)
	return hmac.Sum(nil)
}
