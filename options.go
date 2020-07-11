package logwrap

// Level is an option which sets the messages level.
func Level(l LogLevel) Option {
	return func(message *Message) {
		message.Level = l
	}
}

// Field is an option which adds a single key/value to the fields of the message.
func Field(key string, value interface{}) Option {
	return func(message *Message) {
		message.Fields[key] = value
	}
}

// Fields is an option which takes a list of options and adds them to a message.
func Fields(list List) Option {
	return func(message *Message) {
		for key, value := range list {
			message.Fields[key] = value
		}
	}
}

// List is syntactic sugar to allow users to `Fields(List{"key": "value"})`.
type List map[string]interface{}

const sequenceField = "sequence"

// SequenceAsField is an option which copies the message sequence to the fields.
func SequenceAsField(message *Message) {
	message.Fields[sequenceField] = message.Sequence
}

const errField = "err"

// Err is syntactic sugar to place errors in a messages fields.
func Err(err error) Option {
	return func(message *Message) {
		message.Fields[errField] = err
	}
}
