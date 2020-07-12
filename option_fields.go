package logwrap

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
