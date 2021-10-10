package transport

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/borud/gwp/pkg/gwp"
	"github.com/pion/dtls/v2"
)

type dtlsListener struct {
	addr           string
	listener       net.Listener
	dtlsConfig     *dtls.Config
	requestChannel chan gwp.Request
	ctx            context.Context
	cancel         context.CancelFunc
	rwmu           sync.RWMutex
	children       map[net.Conn]*dtlsConnection
}

// NewDTLSListener creates a new DTLS listener
func NewDTLSListener(addr string, dtlsConfig *dtls.Config, requestChanLen int) (gwp.Listener, error) {
	localAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := dtls.Listen("udp", localAddr, dtlsConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	listener := &dtlsListener{
		addr:           addr,
		listener:       conn,
		dtlsConfig:     dtlsConfig,
		requestChannel: make(chan gwp.Request),
		ctx:            ctx,
		cancel:         cancel,
		rwmu:           sync.RWMutex{},
		children:       make(map[net.Conn]*dtlsConnection),
	}

	go listener.acceptLoop()

	return listener, nil
}

func (l *dtlsListener) Close() error {
	l.cancel()
	l.listener.Close()

	for _, child := range l.children {
		child.Close()
	}

	return nil
}

func (l *dtlsListener) Requests() <-chan gwp.Request {
	return l.requestChannel
}

func (l *dtlsListener) acceptLoop() {
	defer close(l.requestChannel)

	for {
		c, err := l.listener.Accept()
		if err != nil {
			select {
			case <-l.ctx.Done():
				return
			default:
				log.Printf("dtls accept failed: %v", err)
				return
			}
		}

		ctx, cancel := context.WithCancel(context.Background())

		conn := dtlsConnection{
			addr:           c.RemoteAddr().String(),
			conn:           c,
			requestChannel: make(chan gwp.Request),
			ctx:            ctx,
			cancel:         cancel,
		}

		l.rwmu.Lock()
		l.children[c] = &conn
		l.rwmu.Unlock()

		conn.stopped.Add(1)
		go conn.readLoop()

		go func() {
			for req := range conn.Requests() {
				l.requestChannel <- req
			}
			l.rwmu.Lock()
			defer l.rwmu.Unlock()
			delete(l.children, c)
		}()
	}
}
