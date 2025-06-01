package broadcast_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/fmotalleb/the-one/broadcast"
	"github.com/fmotalleb/the-one/logging"
)

type payload = int

func broadCastGenerator() (*broadcast.Broadcaster[payload], chan payload, <-chan payload, <-chan payload, <-chan payload) {
	b := broadcast.NewBroadcaster[payload]()
	input := make(chan payload)
	_, sub1 := b.Subscribe()
	_, sub2 := b.Subscribe(5)
	_, sub3 := b.Subscribe(0)
	return b, input, sub1, sub2, sub3
}

func TestBroadcaster(t *testing.T) {
	e := logging.BootLogger(logging.LogConfig{
		ShowCaller:  false,
		Development: false,
	})
	assert.NoError(t, e)

	assertReceived := func(t *testing.T, ch <-chan payload, expected int, label string) {
		t.Helper()
		select {
		case val := <-ch:
			assert.Equal(t, expected, val, label)
		default:
			t.Errorf("%s did not receive broadcast", label)
		}
	}

	t.Run("MultipleSubscribersReceiveBroadcast", func(t *testing.T) {
		b, input, sub1, sub2, sub3 := broadCastGenerator()
		go b.Broadcast(input)

		input <- 42
		time.Sleep(50 * time.Millisecond)

		assertReceived(t, sub1, 42, "sub1")
		assertReceived(t, sub2, 42, "sub2")
		assertReceived(t, sub3, 42, "sub3")

		close(input)
	})

	t.Run("MultipleSubscribersReceiveBindClose", func(t *testing.T) {
		b, input, sub1, sub2, sub3 := broadCastGenerator()
		go b.BindTo(input)

		input <- 42
		time.Sleep(50 * time.Millisecond)

		assertReceived(t, sub1, 42, "sub1")
		assertReceived(t, sub2, 42, "sub2")
		assertReceived(t, sub3, 42, "sub3")

		close(input)
		<-sub1 // ensure sub1 closes
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		b := broadcast.NewBroadcaster[payload]()
		input := make(chan payload)
		indx, sub := b.Subscribe()
		go b.BindTo(input)

		input <- 42
		time.Sleep(50 * time.Millisecond)

		assertReceived(t, sub, 42, "sub")

		assert.True(t, b.Unsubscribe(indx), "first unsubscribe should succeed")
		assert.False(t, b.Unsubscribe(indx), "second unsubscribe should fail")

		<-sub // ensure closed
	})

	t.Run("helpers.Subscribe", func(t *testing.T) {
		b := broadcast.NewBroadcaster[payload]()
		input := make(chan payload)
		go b.BindTo(input)
		var received []int
		done := make(chan struct{})
		go broadcast.Subscribe(b, func(ch <-chan payload) {
			select {
			case val := <-ch:
				received = append(received, val)
			case <-time.After(1 * time.Second):
				t.Error("timed out waiting for value")
			}
			close(done)
		})

		input <- 42

		<-done

		assert.Equal(t, []int{42}, received)

		// No subscribers should remain
		assert.Equal(t, 0, b.SubscriberCount(), "should have 0 subscribers after Subscribe exits")

		close(input)
	})
}
