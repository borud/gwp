package experiment

// This file contains a collection of unit tests to look at how big the packets end up
// being and how complicated or simple the instantiation of these packets are.  A helpful
// tip when working with protobuffer structures is to use the "fill"-functinality in
// VS Code since it does a good job of filling in the data structures at the various
// levels.

import (
	"log"
	"testing"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"google.golang.org/protobuf/proto"
)

// TestDataPacket tests what a 50 byte payload with a 32 bit recipient address
// ends up looking like.
func TestDataPacket(t *testing.T) {
	p := &gwpb.Packet{
		// Simple 32 bit device address
		To: &gwpb.Address{
			Addr: &gwpb.Address_B32{
				B32: 1234,
			},
		},

		Payload: &gwpb.Packet_Data{
			Data: &gwpb.Data{
				Type: 1,
				Id:   1,
				Data: make([]byte, 50),
			},
		},
	}
	log.Printf("Data packet: %d", wireSize(p))
	log.Printf("%s", p)
}

// TestConfigPacket tests a config packet with 10 fields of varying types and sizes
// sent to a 32 bit address.
func TestConfigPacket(t *testing.T) {
	p := &gwpb.Packet{
		// Simple 32 bit address
		To: &gwpb.Address{Addr: &gwpb.Address_B32{B32: 1234}},

		// Config payload
		Payload: &gwpb.Packet_Config{
			Config: &gwpb.Config{
				Config: map[string]*gwpb.Value{
					"foo": {Value: &gwpb.Value_Int32Val{Int32Val: 1234}},
					"bar": {Value: &gwpb.Value_Int32Val{Int32Val: 4567}},
					"f1":  {Value: &gwpb.Value_FloatVal{FloatVal: 1.1}},
					"f2":  {Value: &gwpb.Value_FloatVal{FloatVal: 10.2}},
					"f3":  {Value: &gwpb.Value_FloatVal{FloatVal: 100.3}},
					"f4":  {Value: &gwpb.Value_FloatVal{FloatVal: 1000.4}},
					"f5":  {Value: &gwpb.Value_FloatVal{FloatVal: 10000.5}},
					"f6":  {Value: &gwpb.Value_FloatVal{FloatVal: 100000.6}},
					"baz": {Value: &gwpb.Value_StringVal{StringVal: "some string"}},
					"xxx": {Value: &gwpb.Value_StringVal{StringVal: "some other string"}},
				},
			},
		},
	}

	log.Printf("Config packet: %d", wireSize(p))
	log.Printf("%s", p)
}

// TestPollConfig is a PollConfig payload that uses a NodeID instead of an
// actual address and asks for 10 named fields.
func TestPollConfig(t *testing.T) {
	p := &gwpb.Packet{
		To: &gwpb.Address{Addr: &gwpb.Address_B32{B32: 123}},
		Payload: &gwpb.Packet_PollConfig{
			PollConfig: &gwpb.PollConfig{
				Fields: []string{"foo", "bar", "baz", "f1", "f2", "f3", "f4", "f5", "f6", "xxx"},
			},
		},
	}

	log.Printf("ConfigPoll packet: %d", wireSize(p))
	log.Printf("%s", p)
}

// TestSample is a Sample payload from a 32 bit address with a single float value.
func TestSample(t *testing.T) {
	p := &gwpb.Packet{
		Payload: &gwpb.Packet_Sample{
			Sample: &gwpb.Sample{
				From:      &gwpb.Address{Addr: &gwpb.Address_B32{B32: 123}},
				Timestamp: uint64(time.Now().UnixMilli()),
				Type:      1,
				Value:     &gwpb.Value{Value: &gwpb.Value_FloatVal{FloatVal: 25.75}},
			},
		},
	}

	log.Printf("Sample packet: %d", wireSize(p))
	log.Printf("%s", p)
}

func wireSize(m proto.Message) int {
	b, _ := proto.Marshal(m)
	return len(b)
}
