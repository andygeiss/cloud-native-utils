package concurrency

import (
	"testing"
	"time"
)

func TestMerge_MultipleChannels(t *testing.T) {
	// Create input channels.
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 2)
	ch3 := make(chan int, 1)
	// Populate input channels.
	ch1 <- 1
	ch1 <- 2
	ch1 <- 3
	close(ch1)
	ch2 <- 4
	ch2 <- 5
	close(ch2)
	ch3 <- 6
	close(ch3)
	// Call the Funnel function.
	out := Merge(ch1, ch2, ch3)
	// Collect output values.
	var result []int
	for val := range out {
		result = append(result, val)
	}
	// Expected values (order may vary due to concurrency).
	expected := []int{1, 2, 3, 4, 5, 6}
	if len(result) != len(expected) {
		t.Fatalf("expected %v elements, got %v", len(expected), len(result))
	}
	// Verify all expected values are present in the result.
	for _, val := range expected {
		found := false
		for _, res := range result {
			if res == val {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("value %v is missing in result %v", val, result)
		}
	}
}

func TestMerge_EmptyChannels(t *testing.T) {
	// Create empty input channels.
	ch1 := make(chan int)
	ch2 := make(chan int)
	// Close input channels immediately.
	close(ch1)
	close(ch2)
	// Call the Funnel function.
	out := Merge(ch1, ch2)
	// Verify the output channel is empty and closed.
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected output channel to be closed, but received a value")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("output channel did not close in time")
	}
}

func TestMerge_SingleChannel(t *testing.T) {
	// Create a single input channel.
	ch := make(chan int, 3)
	// Populate input channel.
	ch <- 10
	ch <- 20
	ch <- 30
	close(ch)
	// Call the Funnel function.
	out := Merge(ch)
	// Collect output values.
	var result []int
	for val := range out {
		result = append(result, val)
	}
	// Check the expected values.
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

func TestMerge_NoChannels(t *testing.T) {
	// Call the Funnel function with no input channels.
	out := Merge[int]()
	// Verify the output channel is closed immediately.
	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected output channel to be closed, but received a value")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("output channel did not close in time")
	}
}

func TestMerge_ConcurrentOutput(t *testing.T) {
	// Create multiple input channels.
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	// Populate input channels.
	go func() {
		for i := 1; i <= 5; i++ {
			ch1 <- i
			time.Sleep(10 * time.Millisecond) // Simulate some delay
		}
		close(ch1)
	}()
	go func() {
		for i := 6; i <= 10; i++ {
			ch2 <- i
			time.Sleep(5 * time.Millisecond) // Simulate a different delay
		}
		close(ch2)
	}()
	// Call the Funnel function.
	out := Merge(ch1, ch2)
	// Collect output values.
	var result []int
	for val := range out {
		result = append(result, val)
	}
	// Verify all expected values are present.
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(result) != len(expected) {
		t.Fatalf("expected %v elements, got %v", len(expected), len(result))
	}
	// Verify all values are within the range 1-10.
	for _, val := range result {
		if val < 1 || val > 10 {
			t.Fatalf("unexpected value %v in result %v", val, result)
		}
	}
}
