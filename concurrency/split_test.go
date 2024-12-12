package concurrency_test

import (
	"cloud-native/concurrency"
	"testing"
	"time"
)

func TestSplit_One_Consumer(t *testing.T) {
	in := []int{1, 2, 3}
	producer := concurrency.Generate[int](in...)
	consumer := concurrency.Split(producer, 1)
	sum := 0
	for range 3 {
		val := <-consumer[0]
		sum += val
	}
	if sum != 6 {
		t.Fatalf("sum must be 6, but got %d", sum)
	}
}

func TestSplit_Two_Consumers(t *testing.T) {
	in := []int{1, 2, 3, 5}
	producer := concurrency.Generate[int](in...)
	consumer := concurrency.Split(producer, 2)
	sum := 0
	go func() {
		for val := range consumer[0] {
			sum += val
		}
	}()
	go func() {
		for val := range consumer[1] {
			sum += val
		}
	}()
	time.Sleep(100 * time.Millisecond)
	if sum != 11 {
		t.Fatalf("sum must be 11, but got %d", sum)
	}
}
