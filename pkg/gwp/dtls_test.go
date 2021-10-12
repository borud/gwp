package gwp

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/borud/gwp/pkg/gwpb"
	"github.com/borud/gwp/pkg/testutil"
	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
	"github.com/stretchr/testify/assert"
)

func TestDTLSListener(t *testing.T) {
	addrStr := fmt.Sprintf(":%d", testutil.FreeUDPPort(t))

	listenerCertificate, err := selfsign.GenerateSelfSigned()
	assert.Nil(t, err)

	dtlsConfig := &dtls.Config{
		Certificates:         []tls.Certificate{listenerCertificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	listener, err := NewDTLSListener(addrStr, dtlsConfig, 1)
	assert.Nil(t, err)

	assert.Nil(t, listener.Close())
}

func TestDTLSConnect(t *testing.T) {
	addrStr := fmt.Sprintf(":%d", testutil.FreeUDPPort(t))

	listenerCertificate, err := selfsign.GenerateSelfSigned()
	assert.Nil(t, err)

	dtlsConfig := &dtls.Config{
		Certificates:         []tls.Certificate{listenerCertificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	listener, err := NewDTLSListener(addrStr, dtlsConfig, 1)
	assert.Nil(t, err)

	// fire up client and send messages
	clientCertificate, err := selfsign.GenerateSelfSigned()
	assert.Nil(t, err)

	clientDtlsConfig := &dtls.Config{
		Certificates:         []tls.Certificate{clientCertificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	conn, err := NewDTLSConnection(addrStr, clientDtlsConfig, 1)
	assert.Nil(t, err)
	assert.NotNil(t, conn)

	assert.Nil(t, conn.Close())

	// close listener
	assert.Nil(t, listener.Close())
}

func TestDTLS(t *testing.T) {
	addrStr := fmt.Sprintf(":%d", testutil.FreeUDPPort(t))

	listenerCertificate, err := selfsign.GenerateSelfSigned()
	assert.Nil(t, err)

	dtlsConfig := &dtls.Config{
		Certificates:         []tls.Certificate{listenerCertificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	listener, err := NewDTLSListener(addrStr, dtlsConfig, 1)
	assert.Nil(t, err)

	// receive and count messages
	numMessages := 20
	waitCh := make(chan struct{})
	go func() {
		defer close(waitCh)

		count := 0
		for range listener.Requests() {
			count++
			if count == numMessages {
				return
			}
		}

	}()

	// fire up client and send messages
	clientCertificate, err := selfsign.GenerateSelfSigned()
	assert.Nil(t, err)

	clientDtlsConfig := &dtls.Config{
		Certificates:         []tls.Certificate{clientCertificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	conn, err := NewDTLSConnection(addrStr, clientDtlsConfig, 1)
	assert.Nil(t, err)
	assert.NotNil(t, conn)

	go func() {
		for i := 0; i < numMessages; i++ {
			err := conn.Send(&gwpb.Packet{
				To: &gwpb.Address{
					Addr: &gwpb.Address_Name{
						Name: "some name",
					},
				},
				Payload: &gwpb.Packet_Data{
					Data: &gwpb.Data{
						Type: 1,
						Data: []byte{1, 2, 3, 4},
					},
				},
			})
			assert.Nil(t, err)
		}
		assert.Nil(t, conn.Close())
	}()

	select {
	case <-waitCh:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "test timed out")
	}

	assert.Nil(t, listener.Close())
}
