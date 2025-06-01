package concurrency_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/fmotalleb/the-one/concurrency"
)

func TestDrain_NonEmptyChannel(t *testing.T) {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3

	result := concurrency.Drain(ch)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDrain_EmptyChannel(t *testing.T) {
	ch := make(chan int, 1)
	result := concurrency.Drain(ch)
	assert.Equal(t, 0, len(result))
}

func TestDrain_ClosedChannelWithData(t *testing.T) {
	ch := make(chan string, 2)
	ch <- "foo"
	ch <- "bar"
	close(ch)

	result := concurrency.Drain(ch)
	assert.Equal(t, []string{"foo", "bar"}, result)
}

func TestDrain_ClosedEmptyChannel(t *testing.T) {
	ch := make(chan float64)
	close(ch)
	result := concurrency.Drain(ch)
	assert.Equal(t, 0, len(result))
}

func TestDrain_PartialBuffered(t *testing.T) {
	ch := make(chan int, 5)
	ch <- 10
	ch <- 20

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- 30 // too late
	}()

	result := concurrency.Drain(ch)
	assert.Equal(t, []int{10, 20}, result)
}
