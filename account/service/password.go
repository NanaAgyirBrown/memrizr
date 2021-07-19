package service

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math/rand"
	"strings"
)

func hashPassword(password string) (string, error){
	// Creating 32 byte salt from go-random package
	salt := make([]byte, 32) // we could have used bcrypt.MaxCost instead of 32
	_, err := rand.Read(salt)

	if err != nil{
		return "",err
	}

	// creating a hash key from byte slice
	shash, err := scrypt.Key([]byte(password), salt,32768, 8, 1, 32)
	if err != nil{
		return "", nil
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW,nil
}

func comparePasswords(storedPassword string, suppliedPassword string) (bool, error)  {
	pwsalt := strings.Split(storedPassword, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("Unable to verify user password")
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)

	return hex.EncodeToString(shash) == pwsalt[0], nil
}