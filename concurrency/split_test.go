package concurrency

import (
	"testing"
	"time"
)

func TestSplit_DistributeValues(t *testing.T) {
	// Create an input channel.
	in := make(chan int, 6)
	// Populate the input channel.
	for i := 1; i <= 6; i++ {
		in <- i
	}
	close(in)
	// Call the Split function with 3 output channels.
	outChannels := Split(in, 3)
	// Collect values from the output channels.
	results := make([][]int, 3)
	for i, ch := range outChannels {
		go func(i int, ch <-chan int) {
			for val := range ch {
				results[i] = append(results[i], val)
			}
		}(i, ch)
	}
	// Allow time for all values to be distributed.
	time.Sleep(100 * time.Millisecond)
	// Ensure all values from the input are in the output channels.
	expectedValues := []int{1, 2, 3, 4, 5, 6}
	collected := make(map[int]bool)
	for _, res := range results {
		for _, val := range res {
			collected[val] = true
		}
	}
	for _, val := range expectedValues {
		if !collected[val] {
			t.Fatalf("value %d from input channel is missing in output channels", val)
		}
	}
}

func TestSplit_EmptyInputChannel(t *testing.T) {
	// Create an empty input channel.
	in := make(chan int)
	close(in)
	// Call the Split function with 3 output channels.
	outChannels := Split(in, 3)
	// Verify all output channels are empty and closed.
	for i, ch := range outChannels {
		select {
		case _, ok := <-ch:
			if ok {
				t.Fatalf("expected channel %d to be closed, but received a value", i)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("channel %d did not close in time", i)
		}
	}
}

func TestSplit_ZeroOutputChannels(t *testing.T) {
	// Create an input channel.
	in := make(chan int, 3)
	in <- 1
	in <- 2
	in <- 3
	close(in)
	// Call the Split function with 0 output channels.
	outChannels := Split(in, 0)
	// Verify no output channels are created.
	if len(outChannels) != 0 {
		t.Fatalf("expected 0 output channels, got %d", len(outChannels))
	}
}

func TestSplit_SingleOutputChannel(t *testing.T) {
	// Create an input channel.
	in := make(chan int, 3)
	in <- 10
	in <- 20
	in <- 30
	close(in)
	// Call the Split function with 1 output channel.
	outChannels := Split(in, 1)
	// Collect values from the single output channel.
	var result []int
	for val := range outChannels[0] {
		result = append(result, val)
	}
	// Verify all values from the input channel are present in the output channel.
	expected := []int{10, 20, 30}
	if len(result) != len(expected) {
		t.Fatalf("expected %v elements, got %v", len(expected), len(result))
	}
	for i, val := range expected {
		if result[i] != val {
			t.Fatalf("expected %v at index %v, got %v", val, i, result[i])
		}
	}
}

func TestSplit_ConcurrentProcessing(t *testing.T) {
	// Create an input channel.
	in := make(chan int, 10)
	// Populate the input channel.
	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
			time.Sleep(10 * time.Millisecond) // Simulate delay in input
		}
		close(in)
	}()
	// Call the Split function with 5 output channels.
	outChannels := Split(in, 5)
	// Collect values from output channels concurrently.
	results := make([][]int, 5)
	for i, ch := range outChannels {
		go func(i int, ch <-chan int) {
			for val := range ch {
				results[i] = append(results[i], val)
			}
		}(i, ch)
	}
	// Allow time for processing.
	time.Sleep(500 * time.Millisecond)
	// Verify all values from the input channel are distributed.
	collected := make(map[int]bool)
	for _, res := range results {
		for _, val := range res {
			collected[val] = true
		}
	}
	for i := 1; i <= 10; i++ {
		if !collected[i] {
			t.Fatalf("value %d from input channel is missing in output channels", i)
		}
	}
}
