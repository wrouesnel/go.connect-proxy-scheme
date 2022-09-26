package connect_proxy_scheme_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/elazarl/goproxy"
	"github.com/samber/lo"
	connect_proxy_scheme "github.com/wrouesnel/go.connect-proxy-scheme"
	"golang.org/x/net/proxy"

	. "gopkg.in/check.v1"
)

func GetHttpClient(proxy proxy.ContextDialer) http.Client {
	tr := &http.Transport{
		DialContext: proxy.DialContext,
	}
	client := http.Client{Transport: tr}
	return client
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type DialerSuite struct {
	server *httptest.Server
	proxy  *httptest.Server
}

var _ = Suite(&DialerSuite{})

func (s *DialerSuite) SetUpSuite(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello")
	})

	s.server = httptest.NewServer(mux)
	s.proxy = httptest.NewServer(goproxy.NewProxyHttpServer())
}

func (s *DialerSuite) TearDownSuite(c *C) {
	s.proxy.Close()
	s.server.Close()
}

func (s *DialerSuite) TestServer(c *C) {
	resp := lo.Must(http.Get(s.server.URL))
	c.Log(string(lo.Must(ioutil.ReadAll(resp.Body))))
}

func (s *DialerSuite) TestProxyDialer(c *C) {
	proxyDialer, err := connect_proxy_scheme.New(lo.Must(url.Parse(s.proxy.URL)), proxy.Direct)
	c.Check(err, IsNil)

	client := GetHttpClient(proxyDialer)
	resp := lo.Must(client.Get(s.server.URL))
	c.Log(string(lo.Must(ioutil.ReadAll(resp.Body))))
}

func (s *DialerSuite) TestRegisteredDialer(c *C) {
	proxy.RegisterDialerType("http", connect_proxy_scheme.ConnectProxy)

	dialer := lo.Must(proxy.FromURL(lo.Must(url.Parse(s.proxy.URL)), proxy.Direct))

	client := GetHttpClient(dialer.(proxy.ContextDialer))
	resp, err := client.Get(s.server.URL)
	c.Check(err, IsNil)
	c.Log(string(lo.Must(ioutil.ReadAll(resp.Body))))
}
