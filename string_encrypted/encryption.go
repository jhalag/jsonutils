package string_encrypted

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// encrypt []byte using AES and provided keyphrase
// extradata:
//   - optional, allows for verification of plaintext data (similar to RSA signature)
//   - must be known to decrypt
//   - alteration of the plaintext data will cause decryption to fail
func Encrypt(keyphrase string, plaintext []byte, extradata []byte) ([]byte, error) {
	//hash the key using SHA256
	key := sha256.Sum256([]byte(keyphrase))

	// 1. Create a new cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// 2. Wrap the block in GCM (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 3. Generate a unique nonce (never reuse a nonce with the same key)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 4. Encrypt and authenticate (Seal)
	// The ciphertext usually includes the authentication tag appended automatically
	// note that the output is appended to nonce
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, extradata)

	return ciphertext, nil
}

// decrypt a provided []byte
func Decrypt(keyphrase string, ciphertext []byte, extradata []byte) ([]byte, error) {
	//hash the key using SHA256
	key := sha256.Sum256([]byte(keyphrase))

	// 1. Create a new cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// 2. Wrap the block in GCM (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 3. Verify the payload is at least long enough to contain a nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("malformed ciphertext: payload too short")
	}

	// 4. Split the packed payload into the nonce and the pure ciphertext
	nonce := ciphertext[:nonceSize]
	pureCiphertext := ciphertext[nonceSize:]

	// 5. Decrypt and authenticate
	plaintext, err := aesGCM.Open(nil, nonce, pureCiphertext, extradata)
	if err != nil {
		return nil, err // Authentication failure (tampered data, wrong key, or wrong nonce)
	}

	return plaintext, nil
}
