package ss

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"net"
)

type Socket struct {
	Pid           uint32 `json:"pid"`
	Family        uint32 `json:"family"`
	Protocol      uint32 `json:"protocol"`
	LocalIP       string `json:"local_ip"`
	LocalPort     int    `json:"local_port"`
	RemoteIP      string `json:"remote_ip"`
	RemotePort    int    `json:"remote_port"`
	Path          string `json:"path"`
	State         string `json:"state"`
	Process       string `json:"process"`
	UID           uint32 `json:"uid"`
	IFace         uint32 `json:"iface"`
	Inode         uint32 `json:"inode"`
	LocalPrivate  bool   `json:"local_private"`
	RemotePrivate bool   `json:"remote_private"`
	Username      string `json:"username"`
}

func (sock *Socket) Private() {
	is := func(v string) bool {
		ip := net.ParseIP(v)
		if ip.IsPrivate() || ip.IsLoopback() || v == "0.0.0.0" || v == "::" || len(v) == 0 {
			return true
		}
		return false
	}

	if sock.LocalIP != "" && is(sock.LocalIP) {
		sock.LocalPrivate = true
	}

	if sock.RemoteIP != "" && is(sock.RemoteIP) {
		sock.RemotePrivate = true
	}
}

func (sock *Socket) Marshal(enc *kind.JsonEncoder) {
	enc.Tab("")

	enc.KV("state", sock.State)
	enc.KV("family", sock.Family)
	enc.KV("protocol", sock.Protocol)
	enc.KV("local_addr", sock.LocalIP)
	enc.KV("local_port", sock.LocalPort)
	enc.KV("local_private", sock.LocalPrivate)
	enc.KV("remote_addr", sock.RemoteIP)
	enc.KV("remote_port", sock.RemotePort)
	enc.KV("remote_private", sock.RemotePrivate)
	enc.KV("pid", sock.Pid)
	enc.KV("inode", sock.Inode)
	enc.KV("process", sock.Process)
	enc.KV("username", sock.Username)
	enc.End("},")
}

func (sock *Socket) Byte() []byte {
	buf := kind.NewJsonEncoder()
	sock.Marshal(buf)
	buf.End("")
	return buf.Bytes()
}

func (sock *Socket) String() string {
	return lua.B2S(sock.Byte())
}
