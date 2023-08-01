package ss

import (
	"github.com/shirou/gopsutil/net"
)

func toSocket(item net.ConnectionStat) *Socket {
	sock := &Socket{
		LocalIP:    item.Laddr.IP,
		LocalPort:  int(item.Laddr.Port),
		RemoteIP:   item.Raddr.IP,
		RemotePort: int(item.Raddr.Port),
		State:      item.Status,
		Pid:        uint32(item.Pid),
		Inode:      item.Fd,
		Family:     item.Family,
		Protocol:   item.Type,
	}

	switch item.Status {

	case "ESTABLISHED":
		sock.State = "ESTAB"
	case "SYN_SENT":
		sock.State = "SYN-SENT"
	case "SYN_RECEIVED":
		sock.State = "SYN-RCVD"
	case "FIN_WAIT_1":
		sock.State = "FIN-WAIT-1"
	case "FIN_WAIT_2":
		sock.State = "FIN-WAIT-2"
	case "TIME_WAIT":
		sock.State = "TIME-WAIT"
	case "CLOSED":
		sock.State = "CLOSED"
	case "CLOSE_WAIT":
		sock.State = "CLOSE-WAIT"
	case "LAST_ACK":
		sock.State = "LAST-ACK"
	case "LISTEN":
		sock.State = "LISTEN"
	case "CLOSING":
		sock.State = "CLOSING"
	case "":
		sock.State = "unknown"
	default:
		sock.State = item.Status
	}

	return sock
}

func (s *ss) tcp4() {
	tbl, err := net.Connections("tcp4")
	if err != nil {
		s.Err = err
		return
	}

	for _, item := range tbl {
		s.handle(toSocket(item))
	}
}

func (s *ss) tcp6() {
	tbl, err := net.Connections("tcp6")
	if err != nil {
		s.Err = err
		return
	}

	for _, item := range tbl {
		s.handle(toSocket(item))
	}
}

func (s *ss) tcp() {
	s.tcp4()
	s.tcp6()
}
