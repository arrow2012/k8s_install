package utils

import (
	"github.com/satori/go.uuid"
	"strings"
)

func GenerateUUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

func GenerateSalt() string {
	sig := uuid.Must(uuid.NewV4(), nil).String()
	sig = strings.Replace(sig, "-", "", -1)
	return sig

}
