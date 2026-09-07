package extensibility_test

import (
	"fmt"
	"log"

	"github.com/andygeiss/cloud-native-utils/extensibility"
)

// This example has no Output comment: it needs a plugin built beside it.
func ExampleLoadPlugin() {
	// The plugin is built separately, by the same toolchain as this program:
	//
	//	go build -buildmode=plugin -o adapter.so adapter.go
	//
	// The consumer declares the interface it wants back, so the plugin never
	// has to import this package.
	type Repository interface {
		FindByID(id string) (name string, err error)
	}

	adapter, err := extensibility.LoadPlugin[Repository]("adapter.so", "Adapter")
	if err != nil {
		log.Fatal(err)
	}

	name, err := adapter.FindByID("1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}
