package mockingbird

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync"
)

func NewMockingbird() *Mockingbird {
	return &Mockingbird{serviceMap: make(map[string]*service)}
}

// DefaultServer is the default instance of *Server.
var DefaultMockingbird = NewMockingbird()

type response struct {
	conn *net.Conn
	req  *http.Request // request for this response

	w *bufio.Writer // buffers output in chunks to chunkWriter
}

func (w *response) Header() http.Header {
	header := make(http.Header)
	return header
}

func (w *response) WriteHeader(code int) {
}

func (w *response) Write(data []byte) (n int, err error) {
	return w.w.Write(data)
}

func (w *response) WriteString(data string) (n int, err error) {
	return w.w.Write([]byte(data))
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

func (mux *ServeMux) Handle(pattern string, handler HandlerFunc) muxEntry {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("http: nil handler")
	}

	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}
	return mux.m[pattern]
}

type ResponseWriter interface {
	// Header returns the header map that will be sent by WriteHeader.
	// Changing the header after a call to WriteHeader (or Write) has
	// no effect.
	Header() http.Header

	// Write writes the data to the connection as part of an HTTP reply.
	// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
	// before writing the data.  If the Header does not contain a
	// Content-Type line, Write adds a Content-Type set to the result of passing
	// the initial 512 bytes of written data to DetectContentType.
	Write([]byte) (int, error)

	WriteString(string) (int, error)

	// WriteHeader sends an HTTP response header with status code.
	// If WriteHeader is not called explicitly, the first call to Write
	// will trigger an implicit WriteHeader(http.StatusOK).
	// Thus explicit calls to WriteHeader are mainly used to
	// send error codes.
	WriteHeader(int)
}

type Request struct {
}

type HandlerFunc func(ResponseWriter, *http.Request)

type muxEntry struct {
	h       HandlerFunc
	pattern string
}

type ServeMux struct {
	m     map[string]muxEntry
	mu    sync.RWMutex
	hosts bool // whether any patterns contain hostnames
}

type methodType struct {
}

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}

type Mockingbird struct {
	serviceMap map[string]*service
}

func Handle(pattern string, handler HandlerFunc) muxEntry {
	return DefaultServeMux.Handle(pattern, handler)
}

func NewServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(mockServiceHandler))

	proxyUrl, _ := url.Parse(ts.URL)
	http.DefaultTransport = &http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return ts
}

func mockServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		panic("not a proxy connection")
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		panic("could not hijack connection")
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		panic(err)
	}

	// because we are proxying the connection
	conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

	// how to detect https conn?

	if r.URL.Scheme == "https" || true {
		fmt.Println("Using HTTPS")
		cert, err := tls.X509KeyPair(localhostCert, localhostKey)
		if err != nil {
			panic(fmt.Sprintf("mockingbird: X509KeyPair: %v", err))
		}

		config := new(tls.Config)
		config.NextProtos = []string{"http/1.1"}
		config.Certificates = []tls.Certificate{cert}
		tlsconn := tls.Server(conn, config)
		defer tlsconn.Close()

		if err := tlsconn.Handshake(); err != nil {
			panic(err)
		}

		conn = tlsconn
	}

	reader := bufio.NewReader(conn)

	req, err := http.ReadRequest(reader)
	if err == io.EOF {
	} else if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(conn)
	w2 := response{conn: &conn, w: writer}
	defer writer.Flush()

	if entry, ok := DefaultServeMux.m[r.URL.String()+req.URL.String()]; ok {
		entry.h(&w2, req)
	}
}

var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBdzCCASOgAwIBAgIBADALBgkqhkiG9w0BAQUwEjEQMA4GA1UEChMHQWNtZSBD
bzAeFw03MDAxMDEwMDAwMDBaFw00OTEyMzEyMzU5NTlaMBIxEDAOBgNVBAoTB0Fj
bWUgQ28wWjALBgkqhkiG9w0BAQEDSwAwSAJBAN55NcYKZeInyTuhcCwFMhDHCmwa
IUSdtXdcbItRB/yfXGBhiex00IaLXQnSU+QZPRZWYqeTEbFSgihqi1PUDy8CAwEA
AaNoMGYwDgYDVR0PAQH/BAQDAgCkMBMGA1UdJQQMMAoGCCsGAQUFBwMBMA8GA1Ud
EwEB/wQFMAMBAf8wLgYDVR0RBCcwJYILZXhhbXBsZS5jb22HBH8AAAGHEAAAAAAA
AAAAAAAAAAAAAAEwCwYJKoZIhvcNAQEFA0EAAoQn/ytgqpiLcZu9XKbCJsJcvkgk
Se6AbGXgSlq+ZCEVo0qIwSgeBqmsJxUu7NCSOwVJLYNEBO2DtIxoYVk+MA==
-----END CERTIFICATE-----`)

// localhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBAN55NcYKZeInyTuhcCwFMhDHCmwaIUSdtXdcbItRB/yfXGBhiex0
0IaLXQnSU+QZPRZWYqeTEbFSgihqi1PUDy8CAwEAAQJBAQdUx66rfh8sYsgfdcvV
NoafYpnEcB5s4m/vSVe6SU7dCK6eYec9f9wpT353ljhDUHq3EbmE4foNzJngh35d
AekCIQDhRQG5Li0Wj8TM4obOnnXUXf1jRv0UkzE9AHWLG5q3AwIhAPzSjpYUDjVW
MCUXgckTpKCuGwbJk7424Nb8bLzf3kllAiA5mUBgjfr/WtFSJdWcPQ4Zt9KTMNKD
EUO0ukpTwEIl6wIhAMbGqZK3zAAFdq8DD2jPx+UJXnh0rnOkZBzDtJ6/iN69AiEA
1Aq8MJgTaYsDQWyU/hDq5YkDJc9e9DSCvUIzqxQWMQE=
-----END RSA PRIVATE KEY-----`)
