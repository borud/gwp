package gwp

import (
	"sync"
	"testing"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/stretchr/testify/assert"
)

func TestMux(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	mux := NewMux()
	mux.AddHandler(&gwpb.Data{}, func(r Request) error {
		wg.Done()
		return nil
	})
	mux.AddHandler(&gwpb.Sample{}, func(r Request) error {
		wg.Done()
		return nil
	})

	// whitebox tests
	assert.Equal(t, 2, len(mux.handlers))
	assert.Contains(t, mux.handlers, "gwpb.Data")
	assert.Contains(t, mux.handlers, "gwpb.Sample")

	// request with data payload
	mux.Handle(Request{
		Peer:       nil,
		RemoteAddr: nil,
		Packet: &gwpb.Packet{
			Payload: &gwpb.Packet_Data{
				Data: &gwpb.Data{
					Type: 1,
					Id:   2,
					Data: []byte{3, 4, 5},
				},
			},
		},
	})

	// request with samples payload
	mux.Handle(Request{
		Peer:       nil,
		RemoteAddr: nil,
		Packet: &gwpb.Packet{
			Payload: &gwpb.Packet_Sample{
				Sample: &gwpb.Sample{
					Type: 1,
					Value: &gwpb.Value{
						Value: &gwpb.Value_FloatVal{
							FloatVal: 3.14,
						},
					},
				},
			},
		},
	})

	mux.Handle(Request{
		Peer:       nil,
		RemoteAddr: nil,
		Packet:     &gwpb.Packet{},
	})
}

func BenchmarkPayloadTypeName(b *testing.B) {
	packet := &gwpb.Packet{
		To: &gwpb.Address{
			Addr: &gwpb.Address_B32{
				B32: 1234,
			},
		},
		Id:         2,
		ResponseTo: 1,
		RequireAck: true,
		Payload:    &gwpb.Packet_Data{Data: &gwpb.Data{Type: 1, Id: 2, Data: []byte{3, 4, 5}}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = payloadTypeName(packet)
	}
}

func BenchmarkPayloadTypeNameNoPayload(b *testing.B) {
	packet := &gwpb.Packet{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = payloadTypeName(packet)
	}
}
