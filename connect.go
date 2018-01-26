package socks5

import (
	"log"
	"net"
)

// Connect remote conn which u want to connect with your dialer
// Error or OK both replied.
func DialWithSource(network, address string, local_address *net.TCPAddr) (net.Conn, error) {
	var d net.Dialer

	d.LocalAddr = local_address

	return d.Dial(network, address)
}

func (r *Request) Connect(c *net.TCPConn) (*net.TCPConn, error) {
	if Debug {
		log.Println("Call X:", r.Address())
		log.Println("Local X:", c.LocalAddr())
	}

	host, _, err := net.SplitHostPort(c.LocalAddr().String())

	log.Println("Local X IP:", host)

	ip := net.ParseIP(host)

	source_addr := &net.TCPAddr{IP: ip}

	tmp, err := DialWithSource("tcp", r.Address(), source_addr)
	if err != nil {
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepHostUnreachable, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepHostUnreachable, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}
	rc := tmp.(*net.TCPConn)

	if Debug {
		log.Println("Dial Remote X:", tmp.RemoteAddr())
		log.Println("Dial Local X:", tmp.LocalAddr())
	}

	a, addr, port, err := ParseAddress(rc.LocalAddr().String())
	if err != nil {
		var p *Reply
		if r.Atyp == ATYPIPv4 || r.Atyp == ATYPDomain {
			p = NewReply(RepHostUnreachable, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = NewReply(RepHostUnreachable, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(c); err != nil {
			return nil, err
		}
		return nil, err
	}
	p := NewReply(RepSuccess, a, addr, port)
	if err := p.WriteTo(c); err != nil {
		return nil, err
	}

	return rc, nil
}
