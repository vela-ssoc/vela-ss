package ss

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/inode"
	"github.com/vela-ssoc/vela-process"
	"golang.org/x/sys/unix"
	"strconv"
)

var notNetlinkSupport = false

func (s *ss) ref(sock *Socket) {
	if s.process == nil {
		s.process = make(map[int32]*process.Process)
	}

	if len(sock.Process) > 0 {
		return
	}

	pid := int32(sock.Pid)

	if p, ok := s.process[pid]; ok {
		sock.Process = p.Executable
		sock.Username = p.Username
		return
	}

	pv, err := process.Fast(int32(sock.Pid))
	if err != nil {
		s.process[pid] = &process.Process{
			Executable: "",
			Username:   "",
		}
		sock.Process = ""
		sock.Username = ""
		return
	}

	sock.Process = pv.Executable
	sock.Username = pv.Username

	s.process[pid] = pv
}

func (s *ss) handle(sock *Socket) {
	//限速
	s.wait()

	//
	sock.Private()

	if s.ov.inode != nil {
		sock.Pid = uint32(s.ov.inode.FindPid(sock.Inode))
	}

	if s.ov.flag.Ref() {
		s.ref(sock)
	}

	if !s.match(sock) {
		return
	}

	if s.ov.hook != nil {
		if s.ov.hook(sock) {
			s.over = true
		}
	}

	if s.ov.vsh != nil {
		s.ov.vsh.Do(sock)
	}

	s.append(sock)
}

func (s *ss) do(family, protocol uint8) {
	if supportNetLink {
		err := s.netlink(family, protocol)
		if err != nil {
			xEnv.Infof("netlink fail %v", err)
		}
		return
	}

	err := s.proc(family, protocol)
	if err != nil {
		s.Err = err
		return
	}
	s.Err = nil
}

func (s *ss) tcp4() {
	s.do(unix.AF_INET, unix.IPPROTO_TCP)
}

func (s *ss) tcp6() {
	s.do(unix.AF_INET6, unix.IPPROTO_TCP)
}

func (s *ss) udp4() {
	s.do(unix.AF_INET, unix.IPPROTO_UDP)
}

func (s *ss) udp6() {
	s.do(unix.AF_INET6, unix.IPPROTO_UDP)
}

func ByPID(pid int32, opt ...func(*option)) *ss {
	ov := newOption()
	ov.cnd = cond.New("pid = " + auxlib.ToString(pid))

	for _, fn := range opt {
		fn(ov)
	}

	/*
		ov.hook = func(s *Socket) (stop bool) {
			s.Pid = uint32(pid)
			return
		}
	*/

	return ssByOption(ov)
}

func ByProcess(p *process.Process, opt ...func(*option)) *ss {
	ov := newOption()
	for _, fn := range opt {
		fn(ov)
	}

	ov.inode = inode.New([]int32{p.Pid})

	ov.hook = func(s *Socket) (stop bool) {
		s.Pid = uint32(p.Pid)
		s.Process = p.Executable
		return
	}

	return ssByOption(ov)

}

func withKernel() {
	info := xEnv.Kernel()
	if info == "" {
		supportNetLink = false
		return
	}

	n, _ := strconv.Atoi(string(info[0]))
	if n >= 3 {
		xEnv.Errorf("%s socket support netlink", info)
		supportNetLink = true
		return
	}

	supportNetLink = false
	xEnv.Errorf("%s socket not support netlink", info)
}
