package distribution

import (
	"context"
	"sync"
	"time"

	"github.com/keboola/keboola-as-code/internal/pkg/idgenerator"
)

// Listener listens for distribution changes, when a node is added or removed.
// It contains the C channel with distribution change Events.
type Listener struct {
	C      chan Events
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	all    *listeners
	id     listenerID
}

type listeners struct {
	config         nodeConfig
	lock           *sync.Mutex
	bufferedEvents Events
	listeners      map[listenerID]*Listener
}

type listenerID string

func newListeners(n *Node) *listeners {
	logger := n.logger.AddPrefix("[listeners]")

	v := &listeners{
		config:    n.config,
		lock:      &sync.Mutex{},
		listeners: make(map[listenerID]*Listener),
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	n.proc.OnShutdown(func() {
		logger.Info("received shutdown request")
		cancel()
		wg.Wait()
		logger.Info("shutdown done")
	})

	wg.Add(1)
	go func() {
		defer wg.Done()

		// If the interval > 0, then listeners are not triggered immediately on change,
		// but all events within the groupInterval are processed at once.
		// Otherwise, trigger is called immediately, see Notify method.
		var tickerC <-chan time.Time
		if v.config.eventsGroupInterval > 0 {
			ticker := n.clock.Ticker(v.config.eventsGroupInterval)
			defer ticker.Stop()
			tickerC = ticker.C
		}

		for {
			select {
			case <-ctx.Done():
				// Handle shutdown
				v.lock.Lock()

				// Log info
				count := len(v.listeners)
				if count > 0 {
					logger.Infof(`waiting for "%d" listeners`, count)
				}

				// Process remaining events
				v.trigger()

				// Stop all listeners
				for _, l := range v.listeners {
					l.cancel()
					l.wg.Wait()
				}

				v.listeners = nil
				v.lock.Unlock()
				return
			case <-tickerC:
				// Trigger listeners at most once per "group interval"
				v.lock.Lock()
				v.trigger()
				v.lock.Unlock()
			}
		}
	}()

	return v
}

func (v *listeners) Reset() {
	v.lock.Lock()
	v.bufferedEvents = nil
	v.lock.Unlock()
}

// Notify listeners about a new event. The event is not processed immediately.
// All events within the "group interval" are processed at once, see trigger method.
func (v *listeners) Notify(events Events) {
	v.lock.Lock()
	defer v.lock.Unlock()

	// All events within the "group interval" are processed at once.
	v.bufferedEvents = append(v.bufferedEvents, events...)

	// Trigger listeners immediately, if there is no grouping interval
	if v.config.eventsGroupInterval == 0 {
		v.trigger()
	}
}

// add a new listener, it contains channel C with streamed distribution change Events.
func (v *listeners) add() *Listener {
	ctx, cancel := context.WithCancel(context.Background())
	out := &Listener{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		all:    v,
		id:     listenerID(idgenerator.Random(10)),
		C:      make(chan Events),
	}
	v.lock.Lock()
	v.listeners[out.id] = out
	v.lock.Unlock()
	return out
}

func (v *listeners) trigger() {
	for _, l := range v.listeners {
		l.trigger(v.bufferedEvents)
	}
	v.bufferedEvents = nil
}

func (l *Listener) Wait(ctx context.Context) (Events, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case events := <-l.C:
		return events, nil
	}
}

func (l *Listener) Stop() {
	l.all.lock.Lock()
	defer l.all.lock.Unlock()

	l.cancel()
	l.wg.Wait()
	delete(l.all.listeners, l.id)
}

func (l *Listener) trigger(events Events) {
	if len(events) == 0 {
		return
	}

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		select {
		case <-l.ctx.Done():
			// stop goroutine on stop/shutdown
		case l.C <- events:
			// propagate events, wait for receiver side
		}
	}()
}
