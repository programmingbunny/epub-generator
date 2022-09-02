package helpers

import (
	"crypto/rand"
	"strings"
)

// helper function for generating unique id for parent directory of epub file
func NumberGen() string {
	p, _ := rand.Prime(rand.Reader, 64)
	return p.String()
}

func TrimImage(imgName string) string {
	new := strings.TrimLeft(imgName, "../chapter-images/")
	return new
}
