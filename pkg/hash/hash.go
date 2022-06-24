package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"math/rand"
)

type KeyManager interface {
	Encode(key string) string
	Decode(key string) (string, error)
}

type gcmKeyManager struct {
	aesBlock cipher.Block
	aesGcm   cipher.AEAD
	nonce    []byte
}

func NewGcmKeyManager() (KeyManager, error) {
	key, err := generateRandomBytes(aes.BlockSize)

	if err != nil {
		return nil, err
	}

	aesBlock, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(aesBlock)

	if err != nil {
		return nil, err
	}

	nonce, err := generateRandomBytes(aesGcm.NonceSize())

	if err != nil {
		return nil, err
	}

	return &gcmKeyManager{
		aesBlock: aesBlock,
		aesGcm:   aesGcm,
		nonce:    nonce,
	}, nil
}

func (m *gcmKeyManager) Encode(key string) string {
	dst := m.aesGcm.Seal(nil, m.nonce, []byte(key), nil)
	return hex.EncodeToString(dst)
}

func (m *gcmKeyManager) Decode(key string) (string, error) {
	src, err := hex.DecodeString(key)

	if err != nil {
		return "", err
	}

	result, err := m.aesGcm.Open(nil, m.nonce, src, nil)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func generateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}
