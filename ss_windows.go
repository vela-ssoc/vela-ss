package ss

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-process"
)

var notSupportNetlink uint32 = 0

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

	if !s.match(sock) {
		return
	}

	if s.ov.flag.Ref() {
		s.ref(sock)
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

func ByPID(pid int32, opt ...func(*option)) *ss {
	ov := newOption()
	ov.cnd = cond.New("pid = " + auxlib.ToString(pid))

	for _, fn := range opt {
		fn(ov)
	}

	return ssByOption(ov)
}

func ByProcess(p *process.Process, opt ...func(*option)) *ss {
	ov := newOption()
	ov.cnd = cond.New("pid = " + auxlib.ToString(p.Pid))
	for _, fn := range opt {
		fn(ov)
	}

	return ssByOption(ov)
}

func withKernel() {
}
