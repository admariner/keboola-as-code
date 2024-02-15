package event_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/jarcoal/httpmock"
	"github.com/keboola/go-client/pkg/keboola"
	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/stretchr/testify/require"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/utctime"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/config"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/dependencies"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/event"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/statistics"
	"github.com/keboola/keboola-as-code/internal/pkg/service/stream/storage/test"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

func TestSender_SendSliceUploadEvent_OkEvent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	var body string
	transport := mock.MockedHTTPTransport()
	registerOkResponder(t, transport, &body)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := error(nil)
	slice := test.NewSlice()
	sender.SendSliceUploadEvent(ctx, api, duration, &err, *slice, testStatsForSlice(slice.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"debug","message":"Sent \"slice-upload\" event id: \"12345\""}`)
	wildcards.Assert(t, `
{
  "component": "keboola.keboola-stream",
  "duration": 3,
  "message": "Slice upload done.",
  "params": "{\"eventName\":\"slice-upload\"}",
  "results": "{\"projectId\":123,\"sinkId\":\"my-sink\",\"sourceId\":\"my-source\",\"statistics\":{\"compressedSize\":52428800,\"firstRecordAt\":\"2000-01-01T20:00:00.000Z\",\"lastRecordAt\":\"2000-01-02T01:00:00.000Z\",\"recordsCount\":123,\"slicesCount\":1,\"stagingSize\":26214400,\"uncompressedSize\":104857600}}",
  "type": "info"
}`, body)
}

func TestSender_SendSliceUploadEvent_ErrorEvent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	var body string
	transport := mock.MockedHTTPTransport()
	registerOkResponder(t, transport, &body)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := errors.New("some error")
	slice := test.NewSlice()
	sender.SendSliceUploadEvent(ctx, api, duration, &err, *slice, testStatsForSlice(slice.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"debug","message":"Sent \"slice-upload\" event id: \"12345\""}`)
	wildcards.Assert(t, `
{
  "component": "keboola.keboola-stream",
  "duration": 3,
  "message": "Slice upload failed.",
  "params": "{\"eventName\":\"slice-upload\"}",
  "results": "{\"error\":\"some error\",\"projectId\":123,\"sinkId\":\"my-sink\",\"sourceId\":\"my-source\"}",
  "type": "error"
}`, body)
}

func TestSender_SendSliceUploadEvent_HTTPError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	transport := mock.MockedHTTPTransport()
	registerErrorResponder(t, transport)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := error(nil)
	slice := test.NewSlice()
	sender.SendSliceUploadEvent(ctx, api, duration, &err, *slice, testStatsForSlice(slice.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"warn","message":"Cannot send \"slice-upload\" event: some error, method: \"POST\", url: \"%s/v2/storage/events\", httpCode: \"403\""}`)
}

func TestSender_SendFileImportEvent_OkEvent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	var body string
	transport := mock.MockedHTTPTransport()
	registerOkResponder(t, transport, &body)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := error(nil)
	file := test.NewFile()
	sender.SendFileImportEvent(ctx, api, duration, &err, file, testStatsForFile(file.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"debug","message":"Sent \"file-import\" event id: \"12345\""}`)
	wildcards.Assert(t, `
{
  "component": "keboola.keboola-stream",
  "duration": 3,
  "message": "File import done.",
  "params": "{\"eventName\":\"file-import\"}",
  "results": "{\"projectId\":123,\"sinkId\":\"my-sink\",\"sourceId\":\"my-source\",\"statistics\":{\"compressedSize\":52428800,\"firstRecordAt\":\"2000-01-01T01:00:00.000Z\",\"lastRecordAt\":\"2000-01-02T01:00:00.000Z\",\"recordsCount\":123,\"slicesCount\":10,\"stagingSize\":26214400,\"uncompressedSize\":104857600}}",
  "type": "info"
}`, body)
}

func TestSender_SendFileImportEvent_ErrorEvent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	var body string
	transport := mock.MockedHTTPTransport()
	registerOkResponder(t, transport, &body)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := errors.New("some error")
	file := test.NewFile()
	sender.SendFileImportEvent(ctx, api, duration, &err, file, testStatsForFile(file.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"debug","message":"Sent \"file-import\" event id: \"12345\""}`)
	wildcards.Assert(t, `
{
  "component": "keboola.keboola-stream",
  "duration": 3,
  "message": "File import failed.",
  "params": "{\"eventName\":\"file-import\"}",
  "results": "{\"error\":\"some error\",\"projectId\":123,\"sinkId\":\"my-sink\",\"sourceId\":\"my-source\"}",
  "type": "error"
}`, body)
}

func TestSender_SendFileImportEvent_HTTPError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	d, mock := dependencies.NewMockedServiceScope(t, config.New())
	api := d.KeboolaPublicAPI().WithToken("my-token")

	transport := mock.MockedHTTPTransport()
	registerErrorResponder(t, transport)

	// Send event
	sender := event.NewSender(d)
	now := utctime.MustParse("2000-01-02T01:00:00.000Z")
	duration := 3 * time.Second
	err := error(nil)
	file := test.NewFile()
	sender.SendFileImportEvent(ctx, api, duration, &err, file, testStatsForFile(file.OpenedAt(), now))

	// Assert
	require.Equal(t, 1, transport.GetCallCountInfo()["POST /v2/storage/events"])
	mock.DebugLogger().AssertJSONMessages(t, `{"level":"warn","message":"Cannot send \"file-import\" event: some error, method: \"POST\", url: \"%s/v2/storage/events\", httpCode: \"403\""}`)
}

func testStatsForSlice(firstAt, lastAt utctime.UTCTime) statistics.Value {
	return statistics.Value{
		SlicesCount:      1,
		FirstRecordAt:    firstAt,
		LastRecordAt:     lastAt,
		RecordsCount:     123,
		UncompressedSize: 100 * datasize.MB,
		CompressedSize:   50 * datasize.MB,
		StagingSize:      25 * datasize.MB,
	}
}

func testStatsForFile(firstAt, lastAt utctime.UTCTime) statistics.Value {
	return statistics.Value{
		SlicesCount:      10,
		FirstRecordAt:    firstAt,
		LastRecordAt:     lastAt,
		RecordsCount:     123,
		UncompressedSize: 100 * datasize.MB,
		CompressedSize:   50 * datasize.MB,
		StagingSize:      25 * datasize.MB,
	}
}

func registerOkResponder(t *testing.T, transport *httpmock.MockTransport, capturedBody *string) {
	t.Helper()
	transport.RegisterResponder(http.MethodPost, "/v2/storage/events", func(req *http.Request) (*http.Response, error) {
		reqBytes, err := httputil.DumpRequest(req, true)
		_, rawBody, _ := bytes.Cut(reqBytes, []byte("\r\n\r\n")) // headers and body are separated by an empty line
		require.NoError(t, err)

		var prettyBody bytes.Buffer
		require.NoError(t, json.Indent(&prettyBody, rawBody, "", "  "))

		*capturedBody = prettyBody.String()

		return httpmock.NewJsonResponderOrPanic(http.StatusCreated, map[string]any{"id": "12345"})(req)
	})
}

func registerErrorResponder(t *testing.T, transport *httpmock.MockTransport) {
	t.Helper()
	errResponse := httpmock.NewJsonResponderOrPanic(http.StatusForbidden, &keboola.StorageError{Message: "some error"})
	transport.RegisterResponder(http.MethodPost, "/v2/storage/events", errResponse)
}
