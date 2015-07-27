package hostroute

import (
	. "github.com/FoxComm/vulcan/location"
	. "github.com/FoxComm/vulcan/netutils"
	. "github.com/FoxComm/vulcan/request"
	. "github.com/FoxComm/vulcan/route"
	. "gopkg.in/check.v1"
	"net/http"
	"testing"
)

func TestPathRoute(t *testing.T) { TestingT(t) }

type HostSuite struct {
}

var _ = Suite(&HostSuite{})

func (s *HostSuite) SetUpSuite(c *C) {
}

func (s *HostSuite) TestRouteEmpty(c *C) {
	m := NewHostRouter()

	out, err := m.Route(request("google.com", "http://google.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, nil)
}

func (s *HostSuite) TestSetNil(c *C) {
	m := NewHostRouter()
	c.Assert(m.SetRouter("google.com", nil), Not(Equals), nil)
}

func (s *HostSuite) TestRouteMatching(c *C) {
	m := NewHostRouter()
	r := &ConstRouter{Location: &Loc{Name: "a"}}
	m.SetRouter("google.com", r)

	out, err := m.Route(request("google.com", "http://google.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, r.Location)
}

func (s *HostSuite) TestRouteMatchingMultiple(c *C) {
	m := NewHostRouter()
	rA := &ConstRouter{Location: &Loc{Name: "a"}}
	rB := &ConstRouter{Location: &Loc{Name: "b"}}
	m.SetRouter("google.com", rA)
	m.SetRouter("yahoo.com", rB)

	out, err := m.Route(request("google.com", "http://google.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, rA.Location)

	out, err = m.Route(request("yahoo.com", "http://yahoo.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, rB.Location)
}

func (s *HostSuite) TestRemove(c *C) {
	m := NewHostRouter()
	rA := &ConstRouter{Location: &Loc{Name: "a"}}
	rB := &ConstRouter{Location: &Loc{Name: "b"}}
	m.SetRouter("google.com", rA)
	m.SetRouter("yahoo.com", rB)

	out, err := m.Route(request("google.com", "http://google.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, rA.Location)

	out, err = m.Route(request("yahoo.com", "http://yahoo.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, rB.Location)

	m.RemoveRouter("yahoo.com")

	out, err = m.Route(request("google.com", "http://google.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, rA.Location)

	out, err = m.Route(request("yahoo.com", "http://yahoo.com/"))
	c.Assert(err, IsNil)
	c.Assert(out, Equals, nil)
}

func request(hostname, url string) Request {
	u := MustParseUrl(url)
	hr := &http.Request{URL: u, Header: make(http.Header), Host: hostname}
	return &BaseRequest{
		HttpRequest: hr,
	}
}
