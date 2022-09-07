package slackinvite

import (
	"fmt"
	"time"

	"github.com/rusq/secure"
)

func generateToken(secret [secretSz]byte) (string, error) {
	return secure.EncryptWithPassphrase(fmt.Sprint(time.Now().Format(time.RFC3339Nano)), secret[:])
}

func verifyToken(token string, secret [secretSz]byte, timeout time.Duration) error {
	pt, err := secure.DecryptWithPassphrase(token, secret[:])
	if err != nil {
		return err
	}
	tt, err := time.Parse(time.RFC3339Nano, pt)
	if err != nil {
		return err
	}

	if age := time.Since(tt); age > timeout {
		return fmt.Errorf("expired: request time: %s, expired %s ago", tt, age-timeout)
	}
	return nil
}
