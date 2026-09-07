package env_test

import (
	"fmt"
	"os"
	"time"

	"github.com/andygeiss/cloud-native-utils/env"
)

func ExampleGet() {
	_ = os.Setenv("PORT", "8443")
	defer func() { _ = os.Unsetenv("PORT") }()

	// The default value decides the type the variable is parsed into.
	port := env.Get("PORT", 8080)

	// An unset variable yields the default.
	timeout := env.Get("SERVER_READ_TIMEOUT", 5*time.Second)

	fmt.Println(port)
	fmt.Println(timeout)
	// Output:
	// 8443
	// 5s
}
