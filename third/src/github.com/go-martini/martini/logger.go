package martini

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		buf, _ := ioutil.ReadAll(req.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		req.Body = rdr1

		// likun add [s] and add reqId to header for xlog
		reqId := genReqId()
		req.Header.Set("X-Reqid", reqId)
		log.Printf("[%s] Started %s %s %s for %s", reqId, req.Method, req.URL.String(), string(buf), addr)

		rw := res.(ResponseWriter)
		c.Next()

		// likun add [s]
		log.Printf("[%s] Completed %v %s in %v\n", reqId, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}

// copy from qiniu log
func genReqId() string {
	var pid = uint32(os.Getpid())
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}
