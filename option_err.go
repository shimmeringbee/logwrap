package logwrap

const errField = "err"

// Err is syntactic sugar to place errors in a messages fields.
func Err(err error) Option {
	return func(message *Message) {
		message.Data[errField] = err.Error()
	}
}
