package logwrap

// Level is an option which sets the messages level.
func Level(l LogLevel) Option {
	return func(message *Message) {
		message.Level = l
	}
}
