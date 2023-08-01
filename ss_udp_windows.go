package ss

import (
	"github.com/shirou/gopsutil/net"
)

func (s *ss) udp4() {
	tbl, err := net.Connections("udp4")
	if err != nil {
		s.Err = err
		return
	}

	for _, item := range tbl {
		s.handle(toSocket(item))
	}
}

func (s *ss) udp6() {
	tbl, err := net.Connections("udp6")
	if err != nil {
		s.Err = err
		return
	}

	for _, item := range tbl {
		s.handle(toSocket(item))
	}
}

func (s *ss) udp() {
	s.udp4()
	s.udp6()
}
