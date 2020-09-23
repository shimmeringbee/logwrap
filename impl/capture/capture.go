package capture

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"sync"
)

// NewCapture initialises a new Capture implementation, the handle to be provided to logwrap should be obtained by
// calling Impl().
func NewCapture() *Capture {
	return &Capture{
		mutex:    &sync.Mutex{},
		messages: []logwrap.Message{},
	}
}

// Capture is a structure which provides a log implementation that captures messages sent to it's implementation.
type Capture struct {
	mutex    *sync.Mutex
	messages []logwrap.Message
}

// Impl returns an implementation that can be passed to logwrap.
func (c *Capture) Impl() logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		c.messages = append(c.messages, message)
	}
}

// Messages returns any messages that have been sent to the Capture so far.
func (c *Capture) Messages() []logwrap.Message {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.messages
}

// Clear resets the structure and wipes any received messages.
func (c *Capture) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.messages = []logwrap.Message{}
}
