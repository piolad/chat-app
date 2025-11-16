package data

type Message struct {
	Message        string
	Timestamp      string
	ConversationID string
	Sender         string
}

type Conversation struct {
	ID            string
	Sender        string
	Receiver      string
	LastTimestamp string
	IVVector      string
}
