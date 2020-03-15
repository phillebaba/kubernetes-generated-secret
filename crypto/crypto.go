package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"

	v1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

// GenerateRandomASCIIString generates a random string of given length, excludes characters that matches any of the options.
func GenerateRandomASCIIString(length int, options []v1alpha1.CharacterOption) (string, error) {
	if len(options) == 4 {
		return "", errors.New("can't exclude all character types")
	}

	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}

		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}

		n := num.Int64()
		if !(n > 32 && n < 127) {
			continue
		}

		c := string(n)
		if matchesOptions(c, options) {
			continue
		}

		result += c
	}
}

// matchesOptions hecks if chartacter matches options regex.
func matchesOptions(char string, options []v1alpha1.CharacterOption) bool {
	result := false
	for _, o := range options {
		validChar := regexp.MustCompile(o.Regex())
		if validChar.MatchString(char) {
			result = true
		}
	}

	return result
}
