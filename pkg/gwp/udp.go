package gwp

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"google.golang.org/protobuf/proto"
)

type udpConnection struct {
	addr           string
	conn           *net.UDPConn
	requestChannel chan Request
	stopped        sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

const (
	// readBufferSize is the size of the readbuffer.  Note that this will always be
	// larger than the maximum packet size we want
	readBufferSize = 1024

	// udpReadTimeout is how long to wait during read.  The practical value of this timeout
	// right now is to provide a means for checking if the listener should be shut down.
	udpReadTimeout = 500 * time.Millisecond
)

// NewUDPListener creates a new UDP listener.
func NewUDPListener(addr string, requestChanLen int) (Listener, error) {
	localAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	listener := &udpConnection{
		addr:           addr,
		conn:           conn,
		requestChannel: make(chan Request, requestChanLen),
		ctx:            ctx,
		cancel:         cancel,
	}

	listener.stopped.Add(1)

	go listener.readLoop()

	return listener, nil
}

// NewUDPConnection creates a connection to a remote server.
func NewUDPConnection(addr string, requestChanLen int) (Connection, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	connection := &udpConnection{
		addr:           addr,
		conn:           conn,
		requestChannel: make(chan Request, requestChanLen),
		ctx:            ctx,
		cancel:         cancel,
	}

	connection.stopped.Add(1)

	go connection.readLoop()

	return connection, nil

}

func (u *udpConnection) Send(packet *gwpb.Packet) error {
	buffer, err := proto.Marshal(packet)
	if err != nil {
		return err
	}

	_, err = u.conn.Write(buffer)

	return err
}

func (u *udpConnection) Requests() <-chan Request {
	return u.requestChannel
}

func (u *udpConnection) Close() error {
	u.cancel()
	u.stopped.Wait()
	return nil
}

func (u *udpConnection) readLoop() {
	defer u.stopped.Done()
	defer close(u.requestChannel)
	defer u.conn.Close()

	buffer := make([]byte, readBufferSize)

	for {
		select {
		case <-u.ctx.Done():
			return
		default:
			// do nothing
		}

		u.conn.SetReadDeadline(time.Now().Add(udpReadTimeout))
		n, remoteAddr, err := u.conn.ReadFrom(buffer)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			continue
		}

		if err != nil {
			log.Printf("UDP Listener, ReadFrom error: %v", err)
			continue
		}

		if n > MaxPacketSize {
			log.Printf("UDP listener, oversize packet: remoteAddr=%s, size=%d maxPacketSize=%d", remoteAddr, n, MaxPacketSize)
		}

		packet := gwpb.Packet{}

		err = proto.Unmarshal(buffer[:n], &packet)
		if err != nil {
			log.Printf("UDP listener, error unmarshalling protobuffer: remoteAddr=%s: %v", remoteAddr.String(), err)
			continue
		}

		// no timeouts for now
		u.requestChannel <- Request{
			Peer:       u,
			RemoteAddr: remoteAddr,
			Packet:     &packet,
			Timestamp:  time.Now(),
		}
	}
}
