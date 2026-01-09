package event_test

import (
	"context"
	"errors"
	"testing"

	"github.com/andygeiss/cloud-native-utils/assert"
	"github.com/andygeiss/cloud-native-utils/event"
)

// -----------------------------------------------------------------------------
// Test implementations for interfaces
// -----------------------------------------------------------------------------

// testEvent is a simple event implementation for testing.
type testEvent struct {
	topic   string
	payload string
}

func (e testEvent) Topic() string {
	return e.topic
}

// mockEventPublisher is a mock implementation of EventPublisher for testing.
type mockEventPublisher struct {
	events    []event.Event
	shouldErr bool
}

func (m *mockEventPublisher) Publish(_ context.Context, e event.Event) error {
	if m.shouldErr {
		return errors.New("mock publish error")
	}
	m.events = append(m.events, e)
	return nil
}

// mockEventSubscriber is a mock implementation of EventSubscriber for testing.
type mockEventSubscriber struct {
	subscriptions map[string]struct {
		factory event.EventFactoryFn
		handler event.EventHandlerFn
	}
	shouldErr bool
}

func newMockEventSubscriber() *mockEventSubscriber {
	return &mockEventSubscriber{
		subscriptions: make(map[string]struct {
			factory event.EventFactoryFn
			handler event.EventHandlerFn
		}),
	}
}

func (m *mockEventSubscriber) Subscribe(_ context.Context, topic string, factory event.EventFactoryFn, handler event.EventHandlerFn) error {
	if m.shouldErr {
		return errors.New("mock subscribe error")
	}
	m.subscriptions[topic] = struct {
		factory event.EventFactoryFn
		handler event.EventHandlerFn
	}{factory, handler}
	return nil
}

// -----------------------------------------------------------------------------
// Event interface tests
// -----------------------------------------------------------------------------

func Test_Event_Topic_Should_ReturnTopic(t *testing.T) {
	// Arrange
	e := testEvent{topic: "test.topic", payload: "test payload"}

	// Act
	topic := e.Topic()

	// Assert
	assert.That(t, "topic must match", topic, "test.topic")
}

// -----------------------------------------------------------------------------
// EventPublisher interface tests
// -----------------------------------------------------------------------------

func Test_EventPublisher_Publish_Should_PublishEvent(t *testing.T) {
	// Arrange
	pub := &mockEventPublisher{}
	e := testEvent{topic: "user.created", payload: "user123"}
	ctx := context.Background()

	// Act
	err := pub.Publish(ctx, e)

	// Assert
	assert.That(t, "must not return error", err, nil)
	assert.That(t, "must have 1 published event", len(pub.events), 1)
	assert.That(t, "event topic must match", pub.events[0].Topic(), "user.created")
}

func Test_EventPublisher_Publish_With_Error_Should_ReturnError(t *testing.T) {
	// Arrange
	pub := &mockEventPublisher{shouldErr: true}
	e := testEvent{topic: "user.created", payload: "user123"}
	ctx := context.Background()

	// Act
	err := pub.Publish(ctx, e)

	// Assert
	assert.That(t, "must return error", err != nil, true)
}

// -----------------------------------------------------------------------------
// EventSubscriber interface tests
// -----------------------------------------------------------------------------

func Test_EventSubscriber_Subscribe_Should_RegisterSubscription(t *testing.T) {
	// Arrange
	sub := newMockEventSubscriber()
	ctx := context.Background()
	factory := func() event.Event { return testEvent{topic: "user.created"} }
	handler := func(_ event.Event) error { return nil }

	// Act
	err := sub.Subscribe(ctx, "user.created", factory, handler)

	// Assert
	assert.That(t, "must not return error", err, nil)
	_, exists := sub.subscriptions["user.created"]
	assert.That(t, "subscription must exist", exists, true)
}

func Test_EventSubscriber_Subscribe_With_Error_Should_ReturnError(t *testing.T) {
	// Arrange
	sub := newMockEventSubscriber()
	sub.shouldErr = true
	ctx := context.Background()
	factory := func() event.Event { return testEvent{topic: "user.created"} }
	handler := func(_ event.Event) error { return nil }

	// Act
	err := sub.Subscribe(ctx, "user.created", factory, handler)

	// Assert
	assert.That(t, "must return error", err != nil, true)
}

// -----------------------------------------------------------------------------
// EventFactoryFn tests
// -----------------------------------------------------------------------------

func Test_EventFactoryFn_Should_CreateEvent(t *testing.T) {
	// Arrange
	factory := func() event.Event {
		return testEvent{topic: "order.created", payload: "order456"}
	}

	// Act
	e := factory()

	// Assert
	assert.That(t, "topic must match", e.Topic(), "order.created")
}

// -----------------------------------------------------------------------------
// EventHandlerFn tests
// -----------------------------------------------------------------------------

func Test_EventHandlerFn_Should_HandleEvent(t *testing.T) {
	// Arrange
	var handledEvent event.Event
	var handler event.EventHandlerFn = func(e event.Event) error {
		handledEvent = e
		return nil
	}
	e := testEvent{topic: "order.shipped", payload: "order789"}

	// Act
	err := handler(e)

	// Assert
	assert.That(t, "must not return error", err, nil)
	assert.That(t, "event must be handled", handledEvent.Topic(), "order.shipped")
}

func Test_EventHandlerFn_With_Error_Should_ReturnError(t *testing.T) {
	// Arrange
	handler := func(_ event.Event) error {
		return errors.New("handler error")
	}
	e := testEvent{topic: "order.shipped", payload: "order789"}

	// Act
	err := handler(e)

	// Assert
	assert.That(t, "must return error", err != nil, true)
}
