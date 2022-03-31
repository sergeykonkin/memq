package memq

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {
	x := 0
	handler1 := func(msg interface{}) {
		x += 1
	}
	b := NewBroker()

	b.Subscribe("topic1", handler1)
	b.Publish("topic1", struct{}{})

	if x != 1 {
		t.Errorf("Expected x to be 1, got %d", x)
	}
}

func TestUnsubscribe(t *testing.T) {
	x := 0
	handler1 := func(msg interface{}) {
		x += 1
	}
	b := NewBroker()

	sub := b.Subscribe("topic1", handler1)
	b.Publish("topic1", struct{}{})
	sub.Unsubscribe()
	b.Publish("topic1", struct{}{})

	if x != 1 {
		t.Errorf("Expected x to be 1, got %d", x)
	}
}

func TestBrokerThreadSafety(t *testing.T) {
	count := 100000
	var handled int32
	handler1 := func(msg interface{}) {
		atomic.AddInt32(&handled, 1)
	}
	handler2 := func(msg interface{}) {
		atomic.AddInt32(&handled, 1)
	}
	handler3 := func(msg interface{}) {
		atomic.AddInt32(&handled, 1)
	}

	b := NewBroker()
	sub1 := b.Subscribe("topic1", handler1)
	sub2 := b.Subscribe("topic1", handler2)
	sub3 := b.Subscribe("topic2", handler3)
	defer sub1.Unsubscribe()
	defer sub2.Unsubscribe()
	defer sub3.Unsubscribe()

	for i := 0; i < count; i++ {
		b.Publish("topic1", i)
		b.Publish("topic2", i)
	}

	time.Sleep(time.Millisecond * 100)

	if int(handled) != count*3 {
		t.Errorf("Expected %d results, got %d", count*3, handled)
	}
}

func TestSubscribeThreadSafety(t *testing.T) {
	count := 100000
	handler1 := func(msg interface{}) {}
	b := NewBroker()

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			b.Subscribe("topic1", handler1)
			wg.Done()
		}()
	}

	wg.Wait()

	if len(b.(*broker).subscribers["topic1"]) != count {
		t.Errorf("Expected %d subscribers, got %d", count, len(b.(*broker).subscribers["topic1"]))
	}
}

func TestUnsubscribeThreadSafety(t *testing.T) {
	count := 100000
	handler1 := func(msg interface{}) {}
	b := NewBroker()

	subs := make([]Subscription, 0, count)
	for i := 0; i < count; i++ {
		subs = append(subs, b.Subscribe("topic1", handler1))
	}

	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		i := i
		go func() {
			subs[i].Unsubscribe()
			wg.Done()
		}()
	}

	wg.Wait()

	if len(b.(*broker).subscribers["topic1"]) != 0 {
		t.Errorf("Expected %d subscribers, got %d", 0, len(b.(*broker).subscribers["topic1"]))
	}
}
