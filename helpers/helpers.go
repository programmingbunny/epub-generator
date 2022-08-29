package helpers

import "crypto/rand"

// helper function for generating unique id for parent directory of epub file
func NumberGen() string {
	p, _ := rand.Prime(rand.Reader, 64)
	return p.String()
}
