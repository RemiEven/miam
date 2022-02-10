package testutils

import (
	"math"
	"math/rand"
	"os"
	"path"
	"time"
)

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().Unix()))
}

func getRandomFileSuffix() string {
	return int64ToSeed(rng.Int63())
}

const radix = '9' - '0' + 'z' - 'a' + 2

func int64ToSeed(value int64) string {
	e := int(math.Floor(math.Log(float64(value)) / math.Log(radix)))
	seed := make([]rune, 0, e)
	posValue := int64(math.Pow(radix, float64(e)))
	for e >= 0 {
		digit := value / posValue
		seed = append(seed, toRune(digit))
		value %= posValue
		posValue /= radix
		e--
	}
	return string(seed)
}

func toRune(value int64) rune {
	if 0 <= value && value <= '9'-'0' {
		return '0' + rune(value)
	}
	if '9'-'0'+1 <= value && value <= radix-1 {
		return 'a' - ('9' - '0' + 1) + rune(value)
	}
	return -1
}

// GetRandomDBFileName returns a randomly generated filename, in the tmp directory, which can be used for the sqlite database in unit tests
func GetRandomDBFileName() string {
	return path.Join(os.TempDir(), "miam_unit_test_"+getRandomFileSuffix()+".db")
}
