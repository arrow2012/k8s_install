package utils

import (
	"github.com/sethvargo/go-password/password"
	"github.com/toolkits/str"
)


func HashIt(passwd, salt string) string {
	return str.Md5Encode(salt + passwd)
}

func GeneratePass(length int) string {
	res, err := password.Generate(length, 5, 0, false, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}
