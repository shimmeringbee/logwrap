package logwrap

// Source is an option to populate the source field of a message.
func Source(source string) Option {
	return func(message *Message) {
		message.Source = source
	}
}
