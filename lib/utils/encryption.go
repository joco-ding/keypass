package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"keypass/lib/stores"
	"log"
	random "math/rand"
	"time"
)

func JsonEncode(v interface{}) []byte {
	_res, _err := json.Marshal(v)
	if _err == nil {
		return _res
	}
	Logging("jsonEncode1", _err.Error())
	return []byte{}
}

func RandThis(_len int, _str string) string {
	rn1 := []rune(_str)

	l1 := len(_str)
	d1 := int64(time.Now().UnixNano())

	s1 := random.NewSource(d1)
	r1 := random.New(s1)
	e1 := ""
	for i1 := 0; i1 < _len; i1++ {
		acak1 := r1.Intn(l1)
		e1 = e1 + string(rn1[acak1:(acak1+1)])

	}
	return e1
}

func GenPass1() string {
	return GenNum(16)
}

func GenPass2() string {
	return GenAbj(16)
}

func GenPass3() string {
	return RandThis(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func GenPass4() string {
	return RandThis(16, "!@#$%^&*()_+=-[]}{\\|;:'\",<.>/?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func GenNum(_len int) string {
	return RandThis(_len, "1234567890")
}

func GenAbj(_len int) string {
	return RandThis(_len, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func GenToken() string {
	return RandThis(24, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
}

func Genckey(_len int) string {
	bytes := make([]byte, _len)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
}

func Genkey() string {
	bytes := make([]byte, 32) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
}

// Encrypt function to encrypt string
func Encrypt(stringToEncrypt string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(stores.Config.KeyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return "Error"
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
		return "Error"
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err.Error())
		return "Error"
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

// Decrypt function to decrypt string
func Decrypt(encryptedString string) (decryptedString string) {

	key, _ := hex.DecodeString(stores.Config.KeyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
		return "Error"
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
		return "Error"
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if nonceSize > len(enc) {
		return "Error"
	}

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err.Error())
		return "Error"
	}

	return fmt.Sprintf("%s", plaintext)
}
