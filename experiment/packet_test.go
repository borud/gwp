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
		Payload: &gwpb.Packet_ConfigPoll{
			ConfigPoll: &gwpb.ConfigPoll{
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

// TestSamples exemplifies an aggregate of measurement values from different
// devices behind the gateway, logging different sensor types, at different
// times.  This is typical if you have devices that log values at different
// rates.
//
// Note that the decision to aggregate values is a function of the gateway and
// does not depend on any assumptions in the protocol.  If you have a device
// that sends a priority message of some sort (fire, door opened etc), you
// have to make a design choice if you send this in a separate packet with
// a single Sample payload or if you combine it into a Samples payload which
// is sent immediately (ie:  send the urgent sample and any packets you have
// already queued up).
//
// One thing you also have to give some thought is the MTU of the transport
// you are using.  You may need to have some logic in the gateway to gauge
// how large the packet is going to be and then decide if you are going to
// send one Samples packet, or multiple Samples packets.
func TestSamples(t *testing.T) {
	p := &gwpb.Packet{
		// From or To not needed since this packet comes from the GW and goes to the server end.
		Payload: &gwpb.Packet_Samples{
			Samples: &gwpb.Samples{
				Samples: []*gwpb.Sample{
					// Sample from device 1
					{
						From:      &gwpb.Address{Addr: &gwpb.Address_B32{B32: 1}},
						Timestamp: uint64(time.Now().UnixMilli()) - 1000,
						Type:      1,
						Value:     &gwpb.Value{Value: &gwpb.Value_FloatVal{FloatVal: 25.5}},
					},

					// Sample from device 2
					{
						From:      &gwpb.Address{Addr: &gwpb.Address_B32{B32: 2}},
						Timestamp: uint64(time.Now().UnixMilli()) - 1000,
						Type:      10,
						Value:     &gwpb.Value{Value: &gwpb.Value_FloatVal{FloatVal: 42.0}},
					},

					// Sample from device 3
					{
						From:      &gwpb.Address{Addr: &gwpb.Address_B32{B32: 3}},
						Timestamp: uint64(time.Now().UnixMilli()) - 1000,
						Type:      100,
						Value:     &gwpb.Value{Value: &gwpb.Value_FloatVal{FloatVal: 1000.1}},
					},
				},
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
