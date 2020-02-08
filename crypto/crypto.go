package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"

	v1alpha1 "github.com/phillebaba/kubernetes-generated-secret/api/v1alpha1"
)

// https://gist.github.com/denisbrodbeck/635a644089868a51eccd6ae22b2eb800
func GenerateRandomASCIIString(length int, options []v1alpha1.ValueOption) (string, error) {
	if len(options) == 0 {
		return "", errors.New("option list can't be empty")
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
		if matchesOptions(c, options) == false {
			continue
		}

		result += c
	}
}

// Checks if chartacter matches options regex
func matchesOptions(char string, options []v1alpha1.ValueOption) bool {
	result := false
	for _, o := range options {
		validChar := regexp.MustCompile(o.Regex())
		if validChar.MatchString(char) {
			result = true
		}
	}

	return result
}
