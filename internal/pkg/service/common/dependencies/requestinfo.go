package dependencies

import (
	"net"
	"net/http"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpserver/middleware"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/ip"
)

type requestInfo struct {
	request   *http.Request
	requestID string
}

func NewRequestInfo(req *http.Request) RequestInfo {
	return newRequestInfo(req)
}

func newRequestInfo(req *http.Request) *requestInfo {
	requestID, _ := req.Context().Value(middleware.RequestIDCtxKey).(string)
	return &requestInfo{request: req, requestID: requestID}
}

func (v *requestInfo) check() {
	if v == nil {
		panic(errors.New("dependencies request info scope is not initialized"))
	}
}

func (v *requestInfo) RequestID() string {
	v.check()
	return v.requestID
}

func (v *requestInfo) RequestHeader() http.Header {
	v.check()
	return v.request.Header.Clone()
}

func (v *requestInfo) RequestClientIP() net.IP {
	v.check()
	return ip.From(v.request)
}
