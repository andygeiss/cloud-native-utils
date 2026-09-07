package env_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native-utils/env"
)

func FuzzGet(f *testing.F) {
	f.Add("8080")
	f.Add("")
	f.Add("not a number")
	f.Add("99999999999999999999999") // overflows every integer type
	f.Add("-0")
	f.Add("5s")
	f.Add("true")
	f.Add("1e309") // overflows float64

	f.Fuzz(func(t *testing.T, value string) {
		// The operating system cannot store a NUL byte in an environment
		// variable, so such an input can never reach Get in a real program.
		if strings.ContainsRune(value, 0) {
			t.Skip("an environment variable cannot hold a NUL byte")
		}
		t.Setenv("FUZZ_VALUE", value)

		// Whatever the variable holds, an unparsable value must fall back to
		// the default rather than panic or return a half-parsed one.
		if got := env.Get("FUZZ_VALUE", -1); got != -1 && value == "" {
			t.Errorf("empty value should have yielded the default, got %d", got)
		}
		_ = env.Get("FUZZ_VALUE", 0.5)
		_ = env.Get("FUZZ_VALUE", false)
		_ = env.Get("FUZZ_VALUE", time.Second)

		// A string default is returned as-is.
		if got := env.Get("FUZZ_VALUE", "fallback"); value != "" && got != value {
			t.Errorf("string value should pass through: %q became %q", value, got)
		}
	})
}
