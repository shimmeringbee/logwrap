package logwrap

const sourceField = "source"

// SourceAsField is an option which copies the message source to the fields, useful for log implementations that
// do not natively support source values.
func SourceAsField(message *Message) {
	message.Data[sourceField] = message.Source
}
