package logwrap

import "fmt"

const trailField = "trail"

// Trail builds a period delimited path, useful for assigning hierarchical identifiers to log messages.
func Trail(s string) Option {
	return func(message *Message) {
		if v, found := message.Data[trailField]; found {
			message.Data[trailField] = fmt.Sprintf("%s.%s", v, s)
		} else {
			message.Data[trailField] = s
		}
	}
}
