package functions

import (
	"moldy-api/utils"

	"github.com/alexedwards/argon2id"
)

func Encrypt(strToEncrypt string) string {
	hash, err := argon2id.CreateHash(strToEncrypt, argon2id.DefaultParams)

	utils.CheckErrors(err, "code 2", "The encrypt process failed", "Unknown solution")

	return hash
}
