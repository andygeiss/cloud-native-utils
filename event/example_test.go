package event_test

import (
	"fmt"

	"github.com/andygeiss/cloud-native-utils/event"
)

// UserCreated is a domain event. Implementing Topic is all it takes.
type UserCreated struct {
	Email string
}

func (e UserCreated) Topic() string { return "user.created" }

func Example() {
	// Domain code speaks in events; the broker adapter stays out of it.
	var e event.Event = UserCreated{Email: "alice@example.com"}
	fmt.Println(e.Topic())
	// Output: user.created
}
