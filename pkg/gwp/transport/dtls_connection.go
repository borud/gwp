package transport

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/borud/gwp/pkg/gwp"
	"github.com/borud/gwp/pkg/gwpb"
	"github.com/pion/dtls/v2"
	"google.golang.org/protobuf/proto"
)

type dtlsConnection struct {
	addr           string
	conn           net.Conn
	dtlsConfig     *dtls.Config
	requestChannel chan gwp.Request
	stopped        sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

const (
	dtlsReadTimeout = 500 * time.Millisecond
)

// NewDTLSConnection creates a connection to a remote server.
func NewDTLSConnection(addr string, dtlsConfig *dtls.Config, requestChanLen int) (gwp.Connection, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := dtls.Dial("udp", remoteAddr, dtlsConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	connection := &dtlsConnection{
		addr:           addr,
		conn:           conn,
		dtlsConfig:     dtlsConfig,
		requestChannel: make(chan gwp.Request, requestChanLen),
		stopped:        sync.WaitGroup{},
		ctx:            ctx,
		cancel:         cancel,
	}

	connection.stopped.Add(1)
	go connection.readLoop()

	return connection, nil
}

func (d *dtlsConnection) Requests() <-chan gwp.Request {
	return d.requestChannel
}

func (d *dtlsConnection) Send(packet *gwpb.Packet) error {
	buffer, err := proto.Marshal(packet)
	if err != nil {
		return err
	}

	_, err = d.conn.Write(buffer)
	return err
}

func (d *dtlsConnection) Close() error {
	d.cancel()
	d.stopped.Wait()
	return nil
}

func (d *dtlsConnection) readLoop() {
	defer d.stopped.Done()
	defer close(d.requestChannel)

	buffer := make([]byte, readBufferSize)
	for {
		select {
		case <-d.ctx.Done():
			d.conn.Close()
			return
		default:
			// do nothing
		}

		d.conn.SetReadDeadline(time.Now().Add(dtlsReadTimeout))
		n, err := d.conn.Read(buffer)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return
		}

		if err != nil && err == io.EOF {
			return
		}

		if err != nil {
			log.Printf("ReadFrom error, exiting: %v", err)
			return
		}

		packet := gwpb.Packet{}
		err = proto.Unmarshal(buffer[:n], &packet)
		if err != nil {
			log.Printf("DTLS listener, error unmarshalling protobuffer remoteAddr=%s: %v", d.conn.RemoteAddr().String(), err)
			continue
		}

		d.requestChannel <- gwp.Request{
			Peer:       d,
			RemoteAddr: d.conn.RemoteAddr(),
			Packet:     &packet,
			Timestamp:  time.Now(),
		}
	}
}
