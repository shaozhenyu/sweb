package xlog

import (
	"bytes"
	"net/http"
	"regexp"
	"testing"

	"github.com/qiniu/log"
	"github.com/stretchr/testify/assert"
)

func TestXlog_Info(t *testing.T) {
	std := log.Std

	out := bytes.Buffer{}
	log.Std = log.New(&out, log.Std.Prefix(), log.Std.Flags())
	NewWith("RhQAAIfWo0-SNUwT").Info("test")
	outStr := out.String()
	assert.True(t, regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}.\d{6}$`).MatchString(outStr[:26]))
	assert.Equal(t, outStr[26:], " [RhQAAIfWo0-SNUwT][INFO] github.com/qiniu/xlog/xlog_test.go:18: test\n")

	log.Std = std
}

// -----------------------------------------------------------------------------

type httpHeader http.Header

func (p httpHeader) ReqId() string {

	return p[reqidKey][0]
}

func (p httpHeader) Header() http.Header {

	return http.Header(p)
}

func TestNewWithHeader(t *testing.T) {

	reqid := "testnewwithheader"

	h := httpHeader(make(http.Header))
	h[logKey] = []string{"origin"}
	h[reqidKey] = []string{reqid}

	xlog := NewWith(h)
	xlog.Xput([]string{"append"})

	assert.Equal(t, h.ReqId(), reqid)
	assert.Equal(t, xlog.ReqId(), reqid)

	log := []string{"origin", "append"}
	assert.Equal(t, h[logKey], log)
	assert.Equal(t, xlog.Xget(), log)
}
