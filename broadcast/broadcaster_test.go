package broadcast_test

import (
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/fmotalleb/the-one/broadcast"
)

func TestBroadcaster_MultipleSubscribersReceiveBroadcast(t *testing.T) {
	type payload = int

	b := broadcast.NewBroadcaster[payload]()
	input := make(chan payload)

	sub1 := b.Subscribe()
	sub2 := b.Subscribe()
	sub3 := b.Subscribe()

	go b.Broadcast(input)

	// Send a value to the broadcaster
	input <- 42
	time.Sleep(50 * time.Millisecond) // allow goroutine to fan out

	// Validate all subscribers received the value
	select {
	case val := <-sub1:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub1 did not receive broadcast")
	}

	select {
	case val := <-sub2:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub2 did not receive broadcast")
	}

	select {
	case val := <-sub3:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub3 did not receive broadcast")
	}

	close(input)
}

func TestBroadcaster_MultipleSubscribersReceiveBindClose(t *testing.T) {
	type payload = int

	b := broadcast.NewBroadcaster[payload]()
	input := make(chan payload)

	sub1 := b.Subscribe()
	sub2 := b.Subscribe()
	sub3 := b.Subscribe()

	go b.BindTo(input)

	// Send a value to the broadcaster
	input <- 42
	time.Sleep(50 * time.Millisecond) // allow goroutine to fan out

	// Validate all subscribers received the value
	select {
	case val := <-sub1:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub1 did not receive broadcast")
	}

	select {
	case val := <-sub2:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub2 did not receive broadcast")
	}

	select {
	case val := <-sub3:
		assert.Equal(t, 42, val)
	default:
		t.Error("sub3 did not receive broadcast")
	}

	close(input)
	<-sub1
}
