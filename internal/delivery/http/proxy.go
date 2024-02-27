package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// listen and serve for proxy
// func proxyHTTPS or tunnelHTTP

func (h *HandlerHTTP) ProxyHTTP(w http.ResponseWriter, req *http.Request) {
	modified, err := modifiedRequest(req)
	if err != nil {
		h.responseErr(w, req, err)
		return
	}

	resp, err := http.DefaultClient.Do(modified)
	if err != nil {
		h.responseErr(w, req, err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func modifiedRequest(r *http.Request) (*http.Request, error) {
	r.Header = changeHeaders(r.Header)

	var err error
	if r.URL, err = modifyPath(r.URL); err != nil {
		return nil, fmt.Errorf("modifying request: %w", err)
	}

	return r, nil
}

func changeHeaders(h http.Header) http.Header {
	h.Del("Proxy-Connection")
	return h
}

func modifyPath(path *url.URL) (*url.URL, error) {
	relative := path.Path

	newUrl, err := url.Parse(relative)
	if err != nil {
		return nil, fmt.Errorf("change absolute URL with relative path: %w", err)
	}

	return newUrl, nil
}

// change all
// make new req
// make req
// write resp to the client writer

// hijacker, ok := w.(http.Hijacker)
// if !ok {
// 	http.Error(w, "Hijacker is not supported", http.StatusInternalServerError)
// }

// conn, buf, err := hijacker.Hijack()
// buf.

// resp, err := http.DefaultTransport.RoundTrip(req)
// if err != nil {
// 	responseError(w, err, p.logger)
// }

// defer resp.Body.Close()
