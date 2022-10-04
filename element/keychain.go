package element

import (
	"sync"
)

/*

  File:    keychain.go
  Author:  Bob Shofner

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: Handle key events for the application.
*/
type keyChain struct {
	lock        sync.Mutex
	input       chan KeyPressed
	subscribe   chan chan<- KeyPressed
	unsubscribe chan chan<- KeyPressed
	outputs     map[chan<- KeyPressed]bool
}

// The KeyChain interface describes the methods of KeyChain(s).
type KeyChain interface {
	// Register a new channel to receive broadcasts
	Register(chan<- KeyPressed)
	// Unregister a channel so that it no longer receives broadcasts.
	Unregister(chan<- KeyPressed)
	// Close this key chain
	Close() error
	// Submit a new object to all subscribers return false if input chan is full
	Submit(KeyPressed) bool
}

// run - forever handling requests.
func (k *keyChain) run() {
	for {
		select {
		case m := <-k.input: // publish (from k.Submit(m interface{}))
			for ch := range k.outputs {
				ch <- m
			}
		case ch, ok := <-k.subscribe: // new subscriber
			if ok {
				k.outputs[ch] = true
			} else {
				return
			}
		case ch := <-k.unsubscribe: // quit subscribing
			delete(k.outputs, ch)
		}
	}
}

// NewKeyChain creates a publisher with the given channel buffer length.
//goland:noinspection GoUnusedExportedFunction
func NewKeyChain(buflen int) KeyChain {
	k := &keyChain{
		input:       make(chan KeyPressed, buflen),    // bi-directional
		subscribe:   make(chan chan<- KeyPressed),     // send only
		unsubscribe: make(chan chan<- KeyPressed),     // send only
		outputs:     make(map[chan<- KeyPressed]bool), // send only
	}
	go k.run()
	return k
}
func (k *keyChain) Close() error {
	defer k.lock.Unlock()
	k.lock.Lock()
	close(k.subscribe)
	close(k.unsubscribe)
	return nil
}
func (k *keyChain) Register(ch chan<- KeyPressed) {
	defer k.lock.Unlock()
	k.lock.Lock()
	k.subscribe <- ch
}
func (k *keyChain) Unregister(ch chan<- KeyPressed) {
	defer k.lock.Unlock()
	k.lock.Lock()
	k.unsubscribe <- ch
}

// Submit attempts to submit an item to be key chained, returning
// true if successful, else false.
func (k *keyChain) Submit(m KeyPressed) bool {
	select { // block until a case can run and buffer has room
	case k.input <- m:
		return true
	default: // input channel is not ready. ignore
		return false
	}
}
