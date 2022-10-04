package misc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"runtime"
)

/*

  File:    encrypt.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

	https://crackstation.net/hashing-security.htm
		and
	https://cryptobook.nakov.com/cryptography-overview
		Argon2d – provides strong GPU resistance, but has potential side-channel attacks
			(possible in very special situations).
		Argon2i – provides less GPU resistance, but has no side-channel attacks.
		Argon2id – recommended (combines the Argon2d and Argon2i).

*/

type HashAlgorithmType string

const (
	PBKDF2   = "PBKDF2"
	MD5      = "MD5"
	SHA256   = "SHA256"
	ARGON2ID = "ARGON2ID"
)

// ArgonMemoryKB is a settable memory size. 64*1024 sets the memory cost to ~64 MB
var ArgonMemoryKB uint32 = 8 * 1024

// CreateHash passphrase id password + salt (pre or post appended)
//goland:noinspection SpellCheckingInspection
//goland:noinspection ALL
func CreateHash(algorithmtype HashAlgorithmType, password, salt []byte, cost int) ([]byte, error) {
	switch algorithmtype {
	case PBKDF2: // auto salted byte[16]. cost: 4 >= cost <= 31. default = 10
		// PBKDF2 is used in WPA-2 and TrueCrypt
		return createPBKDFHash(password, cost)
	case MD5: // creates 64 byte block
		return createMD5Hash(password), nil
	case SHA256: // creates 256 bit hash value
		// prepend application constant
		return createSHA256Hash(append(salt, password...)), nil
	case ARGON2ID: // create a 32 byte (256 bit) key
		return argon2.IDKey(password, salt, uint32(cost), ArgonMemoryKB, uint8(runtime.NumCPU()), 32), nil
	}
	return nil, errors.New(fmt.Sprintf("invalid hashType (%s).", algorithmtype))
}
func createMD5Hash(key []byte) []byte {
	hash := md5.New()
	hash.Write(key)
	return []byte(hex.EncodeToString(hash.Sum(nil)))
}

// CreateSHA256Hash aka SHA2
func createSHA256Hash(key []byte) []byte {
	hash := sha256.Sum256(key) // [32]byte
	return hash[:]             // slice
}
func createPBKDFHash(key []byte, cost int) ([]byte, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	h, err := bcrypt.GenerateFromPassword(key, cost)
	if err != nil {
		return nil, err
	}
	return h, err
}

// CompareHashAndPassword compare a saved hash with a computed hash of supplied password
//  (password is supplied from user via volatile method)
// create a   PBKDF   hash of each and then compare
// this compare includes the auto salt from GCM
// MUST use the CreatePBKDFHash method
//goland:noinspection GoUnusedExportedFunction
func CompareHashAndPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		// log.Println("compare error", err.Error(), "len", len(hashedPassword))
		return false
	}
	return true
}

/*

Galois/Counter Mode (GCM)
From: https://en.wikipedia.org/wiki/Galois/Counter_Mode

Authenticated encryption with associated data (AEAD)
From: https://en.wikipedia.org/wiki/Authenticated_encryption

A typical programming interface for an AE implementation provides
	the following functions:
	Encryption
		Input: plaintext, key, and optionally a header
			that will not be encrypted,
			but will be covered by authenticity protection.
		Output: ciphertext and authentication tag (message authentication code).
	Decryption
		Input: ciphertext, key, authentication tag, and optionally a header (if used during the encryption).
		Output: plaintext, or an error if the authentication tag does not match the supplied ciphertext or header.
		The header part is intended to provide authenticity and integrity
			protection for networking or storage metadata
			for which confidentiality is unnecessary,
			but authenticity is desired.
*/

// EncryptAEAD encrypt with hash prepended and header (associated data)
// AES-256 needs a 32-byte key)
//goland:noinspection ALL,SpellCheckingInspection,SpellCheckingInspection,SpellCheckingInspection
func EncryptAEAD(plainText, key, header []byte) ([]byte, error) {
	//	log.Printf("EncryptAEAD plainText %d key %d header %d :", len(plainText), len(key), len(header))
	dst := make([]byte, 0)
	if len(plainText) < 1 {
		return dst, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return dst, err
	}
	gcm, err2 := cipher.NewGCM(block)
	if err2 != nil {
		return dst, err2
	}
	// compute the GCM with standard nonce length
	// populate nonce with secure random sequence
	// On Linux and FreeBSD, Reader uses getrandom(2) if available, /dev/urandom otherwise.
	// On OpenBSD, Reader uses getentropy(2).
	// On other Unix-like systems, Reader reads from /dev/urandom.
	// On Windows systems, Reader uses the RtlGenRandom API.
	// On Wasm, Reader uses the Web Crypto API.
	nonce := make([]byte, gcm.NonceSize())
	if _, e := io.ReadFull(rand.Reader, nonce); e != nil {
		log.Println("EncryptAEAD err", e)
		return dst, e
	}
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data, returning the returned slice.
	// The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	// pre-pend the nonce
	ciphertext := gcm.Seal(nonce, nonce, plainText, header)
	return ciphertext, nil
}

// DecryptAEAD decrypt from hash prepended and header
// both must match, but header is not used in encryption
//goland:noinspection GoUnusedExportedFunction,SpellCheckingInspection
func DecryptAEAD(cipherText, key, header []byte) ([]byte, error) {
	//	log.Printf("DecryptAEADGCM ciperText %d key %d header %d :", len(cipherText), len(key), len(header))
	dst := make([]byte, 0)
	if len(cipherText) < 1 {
		return dst, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return dst, err
	}
	gcm, err2 := cipher.NewGCM(block)
	if err != nil {
		log.Println("DecryptAEADGCM err", 2)
		return dst, err2
	}
	// the nonce is prepended
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	// authentication requires hash and header
	plaintext, err3 := gcm.Open(nil, nonce, ciphertext, header)
	if err3 != nil {
		log.Println("DecryptAEAD:Open err", err3)
		return dst, err3
	}
	return plaintext, err3
}
