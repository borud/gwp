package gwp

import (
	"errors"
	"log"

	"github.com/borud/gwp/pkg/gwpb"
	"google.golang.org/protobuf/proto"
)

// Mux is a multiplexing Handler.
type Mux struct {
	sealed   AtomicBool
	handlers map[string]Handler
}

var (
	// ErrHandlerAlreadyDefined indicates that you tried to set a handler
	// that has already been set.
	ErrHandlerAlreadyDefined = errors.New("handler already added")

	// ErrCannotResolveName means that the payload type you tried to add a handler
	// for is not a known payload type.
	ErrCannotResolveName = errors.New("cannot resolve name of message type")

	// ErrMuxSealed is returned when you try to add a handler after the
	// mux has been sealed.  The mux is sealed when you explicitly call Seal()
	// or when you.
	ErrMuxSealed = errors.New("you cannot change a sealed mux")
)

// NewMux creates a new Mux instance.  You can add handlers to it
// as long as it is not Seal()'ed.  After it is Sealed you cannot
// alter its state.
func NewMux() *Mux {
	return &Mux{
		handlers: map[string]Handler{},
	}
}

// Handle a request
func (m *Mux) Handle(r Request) {
	name := payloadTypeName(r.Packet)
	if name == "" {
		return
	}

	h, ok := m.handlers[name]
	if ok {
		h(r)
		return
	}
	// TODO(borud): remove this when we have a better way of handling requests.
	log.Printf("Mux: unhandled request: %+v", r)
}

// Seal ensures that the Mux cannot be changed.
func (m *Mux) Seal() {
	m.sealed.SetTrue()
}

// AddHandler adds a handler for a given payload type.
func (m *Mux) AddHandler(payloadType proto.Message, handler Handler) error {
	if m.sealed.IsTrue() {
		return ErrMuxSealed
	}

	name := typeName(payloadType)
	if name == "" {
		return ErrCannotResolveName
	}

	_, ok := m.handlers[name]
	if ok {
		return ErrHandlerAlreadyDefined
	}

	m.handlers[name] = handler
	return nil
}

// typeName returns the name of the message type
func typeName(p proto.Message) string {
	return string(p.ProtoReflect().Descriptor().FullName())
}

// payloadTypeName returns the full name of the payload message type or
// an empty string if the payload is empty.
// While this looks messy, benchmarking shows that this is typically in
// the 100ns range on a Intel(R) Core(TM) i5-7500 CPU @ 3.40GHz.
func payloadTypeName(p *gwpb.Packet) string {
	oneof := p.ProtoReflect().Descriptor().Oneofs().ByName("payload")
	if oneof == nil {
		return ""
	}

	which := p.ProtoReflect().WhichOneof(oneof)
	if which == nil {
		return ""
	}

	return string(which.Message().FullName())
}
