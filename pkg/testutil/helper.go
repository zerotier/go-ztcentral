package testutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TokenFile is the filename of where the token resides. Whitespace will be
// trimmed at InitToken() time on read. You can override this before calling
// InitToken() or InitTokenFromEnv().
var TokenFile = "test-token.txt"

// TokenEnv is the name of the environment variable InitTokenFromEnv() reads.
// Also overridable.
var TokenEnv = "ZEROTIER_CENTRAL_TOKEN"

// InitTokenFromEnv initializes the token from either the file or the
// environment variable. See TokenEnv and TokenFile.
func InitTokenFromEnv() string {
	return InitToken(os.Getenv(TokenEnv))
}

// InitToken will attempt to read the string passed; if it is empty, it will
// attempt to read it from the file (see TokenFile).
func InitToken(controllerToken string) string {
	if controllerToken == "" {
		if fi, err := os.Stat("test-token.txt"); err != nil {
			fmt.Fprintln(os.Stderr, "test-token.txt not present in tree; ZEROTIER_CENTRAL_TOKEN is required in environment for many tests.")
		} else if fi.Mode()&os.ModeIrregular != 0 {
			panic("test-token.txt is not a regular file; not sure what to do here, so bailing")
		} else {
			content, err := ioutil.ReadFile("test-token.txt")
			if err != nil {
				panic(err)
			}

			controllerToken = strings.TrimSpace(string(content))
		}
	}

	return controllerToken
}

// NeedsToken runs InitTokenFromEnv and skips the test if it returns an empty string.
func NeedsToken(t *testing.T) {
	if InitTokenFromEnv() == "" {
		t.Skipf("This test requires %q be set in the environment, or %q exists in the repository root with the token inside.", TokenEnv, TokenFile)
	}
}
