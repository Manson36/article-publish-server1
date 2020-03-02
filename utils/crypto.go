package utils

import (
	"crypto/md5"
	"encoding/hex"
)

type PWDInfo struct {
	Password string
	Salt     string
}

const seedString = "2345678abcdefhijkmnprstwxyzABCDEFGHJKMNPQRSTWXYZ"

func GetPwdSalt() string {
	return GetRandomString(5)
}

func HashPwdWithSalt(password, salt string) string {
	h := md5.New()
	h.Write([]byte(salt + "%$#" + password))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
func GenPwdAndSalt(password string) PWDInfo {
	salt := GetPwdSalt()
	pwd := HashPwdWithSalt(password, salt)
	return PWDInfo{Password: pwd, Salt: salt}
}
