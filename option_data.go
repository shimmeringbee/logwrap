package logwrap

// Datum is an option which adds a single key/value to the data of the message.
func Datum(key string, value interface{}) Option {
	return func(message *Message) {
		message.Data[key] = value
	}
}

// Data is an option which takes a list of options and adds them to a message.
func Data(list List) Option {
	return func(message *Message) {
		for key, value := range list {
			message.Data[key] = value
		}
	}
}

// List is syntactic sugar to allow users to `Data(List{"key": "value"})`.
type List map[string]interface{}
