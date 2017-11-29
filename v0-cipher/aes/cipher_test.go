package aes

import (
	"testing"

	"mobilehound/v0-test"
)

func TestAES(t *testing.T) {
	test.CipherTest(t, NewCipher128)
}
