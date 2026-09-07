package resource_test

import (
	"context"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/resource"
)

// Every backend implements the same Access[K, V] interface, so swapping
// NewInMemoryAccess for a file, SQLite or Postgres backend changes nothing else.
func Example() {
	store := resource.NewInMemoryAccess[string, string]()
	ctx := context.Background()

	_ = store.Create(ctx, "user-1", "Alice")
	value, _ := store.Read(ctx, "user-1")
	fmt.Println(*value)

	_ = store.Update(ctx, "user-1", "Alice Smith")
	value, _ = store.Read(ctx, "user-1")
	fmt.Println(*value)

	_ = store.Delete(ctx, "user-1")
	_, err := store.Read(ctx, "user-1")
	fmt.Println(err)
	// Output:
	// Alice
	// Alice Smith
	// resource not found
}

func ExampleNewIndexedAccess() {
	type user struct {
		Email string
		Role  string
	}

	store := resource.NewIndexedAccess(resource.NewInMemoryAccess[string, user]())
	store.AddIndex("role", func(u user) string { return u.Role })
	ctx := context.Background()

	_ = store.Create(ctx, "1", user{Email: "alice@example.com", Role: "admin"})
	_ = store.Create(ctx, "2", user{Email: "bob@example.com", Role: "admin"})
	_ = store.Create(ctx, "3", user{Email: "carol@example.com", Role: "guest"})

	admins, err := store.FindByIndex(ctx, "role", "admin")
	fmt.Println(len(admins), err)
	// Output: 2 <nil>
}
