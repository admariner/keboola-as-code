package etcdop

import (
	"context"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/keboola/go-utils/pkg/wildcards"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/etcdhelper"
)

// nolint:paralleltest // the test run the "compact" operation and breaks the other tests running in parallel
func TestWatchConsumer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	logger := log.NewDebugLogger()

	// Both clients must use same namespace
	etcdNamespace := "unit-" + t.Name() + "-" + gonanoid.Must(8)

	// Create client for the test
	testClient := etcdhelper.ClientForTestWithNamespace(t, etcdNamespace)

	// Create watcher client with custom dialer
	var conn net.Conn
	dialerLock := &sync.Mutex{}
	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		dialerLock.Lock()
		defer dialerLock.Unlock()
		var err error
		conn, err = (&net.Dialer{}).DialContext(ctx, "tcp", s)
		return conn, err
	}
	c := etcdhelper.ClientForTestWithNamespace(t, etcdNamespace, grpc.WithContextDialer(dialer))

	// Create consumer
	pfx := prefixForTest()
	init := pfx.
		GetAllAndWatch(ctx, c).
		SetupConsumer(logger).
		WithOnCreated(func(header *Header) {
			logger.Infof(`OnCreated: created (rev %v)`, header.Revision)
		}).
		WithOnRestarted(func(reason string, delay time.Duration) {
			logger.Infof(`OnRestarted: %s`, reason)
		}).
		WithForEach(func(events []WatchEvent, header *Header, restart bool) {
			var str strings.Builder
			for _, e := range events {
				str.WriteString(e.Type.String())
				str.WriteString(` "`)
				str.Write(e.Kv.Key)
				str.WriteString(`", `)
			}
			logger.Infof(`ForEach: restart=%t, events(%d): %s`, restart, len(events), strings.TrimSuffix(str.String(), ", "))
		}).
		StartConsumer(wg)

	// Wait for initialization
	assert.NoError(t, <-init)

	// Expect created event
	wildcards.Assert(t, "INFO  OnCreated: created (rev %d)", logger.AllMessages())
	logger.Truncate()

	// Put some key
	assert.NoError(t, pfx.Key("key1").Put("value1").Do(ctx, c))

	// Expect forEach event
	assert.Eventually(t, func() bool {
		return strings.Count(logger.AllMessages(), "ForEach:") == 1
	}, time.Second, 10*time.Millisecond)
	wildcards.Assert(t, `
INFO  ForEach: restart=false, events(1): create "my/prefix/key1"
`, logger.AllMessages())
	logger.Truncate()

	// Close watcher connection and block a new one
	dialerLock.Lock()
	assert.NoError(t, conn.Close())

	// Add some other keys, during the watcher is disconnected
	assert.NoError(t, pfx.Key("key2").Put("value2").Do(ctx, testClient))
	assert.NoError(t, pfx.Key("key3").Put("value3").Do(ctx, testClient))

	// Compact, during the watcher is disconnected
	status, err := testClient.Status(ctx, c.Endpoints()[0])
	assert.NoError(t, err)
	_, err = testClient.Compact(ctx, status.Header.Revision)
	assert.NoError(t, err)

	// Unblock dialer, watcher will be reconnected
	dialerLock.Unlock()

	// Expect restart event, followed with all 3 keys.
	// The restart flag is true.
	assert.Eventually(t, func() bool {
		return strings.Count(logger.AllMessages(), "my/prefix/key") == 3
	}, time.Second, 10*time.Millisecond)
	wildcards.Assert(t, `
WARN  etcdserver: mvcc: required revision has been compacted
WARN  restarted after %s, reason: etcdserver: mvcc: required revision has been compacted
INFO  OnRestarted: restarted after %s, reason: etcdserver: mvcc: required revision has been compacted
INFO  ForEach: restart=true, events(3): create "my/prefix/key1", create "my/prefix/key2", create "my/prefix/key3"
`, logger.AllMessages())
	logger.Truncate()

	// The restart flag is false in further events.
	assert.NoError(t, pfx.Key("key4").Put("value4").Do(ctx, testClient))
	assert.Eventually(t, func() bool {
		return strings.Count(logger.AllMessages(), "ForEach:") == 1
	}, time.Second, 10*time.Millisecond)
	wildcards.Assert(t, `
INFO  ForEach: restart=false, events(1): create "my/prefix/key4"
`, logger.AllMessages())
	logger.Truncate()

	// Stop
	cancel()
	wg.Wait()
}