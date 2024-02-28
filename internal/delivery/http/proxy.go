package http

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	headersToSkip = []string{"Proxy-Connection"}
	dialTimeout   = 5 * time.Second
)

func (h *HandlerHTTP) ProxyHTTP(w http.ResponseWriter, req *http.Request) {
	new, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		h.responseErr(w, req, err)
		return
	}
	cpyHeaders(new.Header, req.Header)

	resp, err := h.proxyClient.Do(new)
	if err != nil {
		h.responseErr(w, req, err)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (h *HandlerHTTP) tunnelHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.RequestURI)
	dstConn, err := net.DialTimeout("tcp", req.Host, dialTimeout)
	if err != nil {
		h.responseErr(w, req, err)
		return
	}
	w.WriteHeader(http.StatusOK)

	// b := &bytes.Buffer{}
	// req.Write(b)
	// fmt.Println(b.String())
	// fmt.Printf("%+v\n", req)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		h.responseErr(w, req, errors.New("hijacker isn't supported"))
		return
	}

	srcConn, _, err := hijacker.Hijack()
	if err != nil {
		h.responseErr(w, req, err)
		return
	}

	go streamTCP(dstConn, srcConn)
	go streamTCP(srcConn, dstConn)
}

func streamTCP(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func cpyHeaders(dst, src http.Header) {
	for key, values := range src {
		if !shouldBeSkipped(key) {
			dst[key] = values
		}
	}
}

func shouldBeSkipped(key string) bool {
	for _, h := range headersToSkip {
		if key == h {
			return true
		}
	}
	return false
}
