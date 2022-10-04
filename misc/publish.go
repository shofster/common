package misc

import "sync"

/*

  File:    publish.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Publish / Subscribe.
*/

type publisher struct {
	lock        sync.Mutex
	input       chan interface{}
	subscribe   chan chan<- interface{}
	unsubscribe chan chan<- interface{}
	outputs     map[chan<- interface{}]bool
}

// The Publisher interface describes the main entry points to Publisher(s).
type Publisher interface {
	// Register a new channel to receive broadcasts
	Register(chan<- interface{})
	// Unregister a channel so that it no longer receives broadcasts.
	Unregister(chan<- interface{})
	// Close this publisher
	Close() error
	// Submit a new object to all subscribers return false if input chan is full
	Submit(interface{}) bool
}

// run - forever handling requests.
func (p *publisher) run() {
	for {
		select {
		case m := <-p.input: // publish (from p.Submit(m interface{}))
			for ch := range p.outputs {
				ch <- m
			}
		case ch, ok := <-p.subscribe: // new subscriber
			if ok {
				p.outputs[ch] = true
			} else {
				return
			}
		case ch := <-p.unsubscribe: // quit subscribing
			delete(p.outputs, ch)
		}
	}
}

// NewPublisher creates a publisher with the given channel buffer length.
//goland:noinspection GoUnusedExportedFunction
func NewPublisher(buflen int) Publisher {
	p := &publisher{
		input:       make(chan interface{}, buflen),    // bi-directional
		subscribe:   make(chan chan<- interface{}),     // send only
		unsubscribe: make(chan chan<- interface{}),     // send only
		outputs:     make(map[chan<- interface{}]bool), // send only
	}
	go p.run()
	return p
}
func (p *publisher) Close() error {
	defer p.lock.Unlock()
	p.lock.Lock()
	close(p.subscribe)
	close(p.unsubscribe)
	return nil
}
func (p *publisher) Register(ch chan<- interface{}) {
	defer p.lock.Unlock()
	p.lock.Lock()
	p.subscribe <- ch
}
func (p *publisher) Unregister(ch chan<- interface{}) {
	defer p.lock.Unlock()
	p.lock.Lock()
	p.unsubscribe <- ch
}

// Submit attempts to submit an item to be published, returning
// true if successful, else false.
func (p *publisher) Submit(m interface{}) bool {
	select { // block until a case can run and buffer has room
	case p.input <- m:
		return true
	default: // input channel is not ready. ignore
		return false
	}
}
