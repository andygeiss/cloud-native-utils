package service_test

import (
	"context"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/service"
)

func ExampleWrap() {
	// Wrap turns a plain function into one that honours a context.
	double := service.Wrap(func(n int) (int, error) { return n * 2, nil })

	out, err := double(context.Background(), 21)
	fmt.Println(out, err)

	// A cancelled context stops the call before the function runs.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = double(ctx, 21)
	fmt.Println(err)
	// Output:
	// 42 <nil>
	// context canceled
}
