package logwrap

const sequenceField = "sequence"

// SequenceAsField is an option which copies the message sequence to the fields.
func SequenceAsField(message *Message) {
	message.Fields[sequenceField] = message.Sequence
}
