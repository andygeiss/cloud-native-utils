package messaging_test

import (
	"context"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/messaging"
	"github.com/andygeiss/cloud-native-utils/service"
)

func ExampleNewInternalDispatcher() {
	dispatcher := messaging.NewInternalDispatcher()
	ctx := context.Background()

	_ = dispatcher.Subscribe(ctx, "user.created", service.Wrap(func(m messaging.Message) (messaging.MessageState, error) {
		fmt.Printf("%s: %s\n", m.Topic, m.Data)
		return messaging.MessageStateCompleted, nil
	}))

	// Publish returns once every subscriber has finished.
	_ = dispatcher.Publish(ctx, messaging.NewMessage("user.created", []byte("alice")))
	// Output: user.created: alice
}
