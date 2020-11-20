// Package app encapsulates the client, server and other interfaces to provide a complete dapp
package app

// Program is an interface for distributed application programming
type Program interface {
	// Set the current application name
	Name(string)
	// Request an application by name and endpoint
	Request(name, ep string, req, rsp interface{}) error
	// Broadcast a message to all subscribers
	Broadcast(topic string, msg interface{}) error
	// Register a handler e.g a public Go struct/method with signature func(context.Context, *Request, *Response) error
	Handle(v interface{}) error
	// Subscribe to broadcast messages. Signature is public Go func or struct with signature func(context.Context, *Message) error
	Subscribe(topic string, v interface{}) error
	// Run the application
	Run() error
}
