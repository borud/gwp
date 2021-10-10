package gwp

import (
	"errors"
	"log"
)

type Server struct {
	Listeners         []Listener
	Handler           Handler
	doneChan          chan struct{}
	shutdown          AtomicBool
	shutdownCallbacks []func()
}

var (
	// ErrServerAlreadyStarted means that the server was already started.  It is an error to
	// start a server more than once.
	ErrServerAlreadyStarted = errors.New("server already started")

	// ErrNoListenersDefined means that the server was started without any listeners.  Since
	// this is pointless we define it as an error.
	ErrNoListenersDefined = errors.New("no listeners defined")

	// ErrNoHandlerDefined means that we did not specify a handler.
	ErrNoHandlerDefined = errors.New("no handler defined")

	// ErrAlreadyShutDown means we have already shut down the server.  It is an error
	// to shut it down twice.
	ErrAlreadyShutDown = errors.New("server has already been shut down")
)

// Start the server.
func (s *Server) Start() error {
	if s.doneChan != nil {
		return ErrServerAlreadyStarted
	}

	if len(s.Listeners) == 0 {
		return ErrNoListenersDefined
	}

	if s.Handler == nil {
		return ErrNoHandlerDefined
	}

	s.doneChan = make(chan struct{})
	s.shutdown.SetFalse()

	// Handle incoming requests, each in their own goroutine
	for _, lis := range s.Listeners {
		requestChan := lis.Requests()
		go func() {
			for req := range requestChan {
				go s.Handler(req)

				select {
				case <-s.doneChan:
					return
				default:
				}
			}
		}()
	}

	return nil
}

// Shutdown server.
func (s *Server) Shutdown() error {
	if s.shutdown.IsTrue() {
		return ErrAlreadyShutDown
	}
	s.shutdown.SetTrue()

	// Shut down listeners
	for _, lis := range s.Listeners {
		err := lis.Close()
		if err != nil {
			log.Printf("listener %+v shutdown error: %v", lis, err)
		}
	}

	close(s.doneChan)

	// Call
	for _, cb := range s.shutdownCallbacks {
		cb()
	}

	return nil
}
