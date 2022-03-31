package memq

import (
	"sync"
)

type Broker interface {
	Publish(topic string, msg interface{})
	Subscribe(topic string, handler func(interface{})) Subscription
}

type Subscription interface {
	Unsubscribe()
}

type broker struct {
	identity    int
	subscribers map[string]map[int]subscription
	m           *sync.Mutex
}

type subscription struct {
	topic string
	id    int
	b     *broker
	ch    chan interface{}
	done  chan struct{}
}

func NewBroker() Broker {
	return &broker{
		subscribers: make(map[string]map[int]subscription),
		m:           &sync.Mutex{},
	}
}

func (b *broker) Subscribe(topic string, handler func(interface{})) Subscription {
	b.m.Lock()
	defer b.m.Unlock()
	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make(map[int]subscription)
	}
	topicSubs := b.subscribers[topic]
	subID := b.identity
	b.identity++
	ch := make(chan interface{})
	done := make(chan struct{})
	sub := subscription{
		topic: topic,
		id:    subID,
		b:     b,
		ch:    ch,
		done:  done,
	}
	topicSubs[subID] = sub
	go func() {
		for {
			select {
			case <-done:
				return
			case msg := <-ch:
				handler(msg)
			}
		}
	}()
	return &sub
}

func (s *subscription) Unsubscribe() {
	s.b.m.Lock()
	defer s.b.m.Unlock()
	delete(s.b.subscribers[s.topic], s.id)
	close(s.done)
	close(s.ch)
}

func (b *broker) Publish(topic string, msg interface{}) {
	if b.subscribers == nil {
		return
	}

	topicSubs, ok := b.subscribers[topic]
	if !ok {
		return
	}
	for _, sub := range topicSubs {
		select {
		case <-sub.done:
			continue
		case sub.ch <- msg:
		}
	}
}
