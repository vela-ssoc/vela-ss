package ss

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-process"
	vswitch "github.com/vela-ssoc/vela-switch"
)

var xEnv vela.Environment

func newListenSnapshotL(L *lua.LState) int {
	snap := newListenSnapshot(L)
	proc := L.NewVelaData(snap.Name(), lnTypeof)
	proc.Set(snap)
	L.Push(proc)
	return 1
}

/*
	vela.ss("-p -t" , "addr = 127.0.0.1")
*/

// vela.ss.pid("-t -a -l" , 10 , cnd)
func pidL(L *lua.LState) int {
	of := NewOptionFlagL(L, 1)
	of.ref = true

	pid := L.IsInt(2)
	cnd := cond.CheckMany(L, cond.Seek(3))
	L.Push(ByPID(pid, Flag(of), Pid([]int{pid}), Cnd(cnd)))
	return 1
}

func processL(L *lua.LState) int {
	of := NewOptionFlagL(L, 1)
	of.ref = true //关联进程
	pro := process.CheckById(L, 2)
	cnd := cond.CheckMany(L, cond.Seek(3))
	L.Push(ByProcess(pro, Flag(of), Cnd(cnd), Pid([]int{pro.Pid})))
	return 1
}

// vela.ss.switch('-p -s -v" , {})

func switchL(L *lua.LState) int {
	of := NewOptionFlagL(L, 1)
	vsh := vswitch.NewL(L)
	vsh.NewSwitchByLTab(L.CheckTable(2))

	sum := By(Flag(of), Switch(vsh))
	if sum == nil {
		L.Push(&ss{Err: fmt.Errorf("invalid options")})
	} else {
		L.Push(sum)
	}
	return 1
}

func indexL(L *lua.LState) int {
	v := L.Get(1).String()
	of, _ := NewOptionFlag("-a")
	var ret *Socket

	By(Flag(of), Hook(func(sock *Socket) (stop bool) {
		if sock == nil {
			return
		}

		if auxlib.ToString(sock.Inode) == v {
			ret = sock
			stop = true
		}
		return
	}))

	if ret == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(ret)
	}

	return 1
}

func call(L *lua.LState) int {
	of := NewOptionFlagL(L, 1)
	cnd := cond.CheckMany(L, cond.Seek(1))
	sum := By(Flag(of), Cnd(cnd))
	if sum == nil {
		L.Push(&ss{Err: fmt.Errorf("invalid options")})
	} else {
		L.Push(sum)
	}
	return 1
}

// ss.inode(123123)
func WithEnv(env vela.Environment) {
	xEnv = env
	withKernel()
	xEnv.Mime(&listen{}, encode, decode)
	kv := lua.NewUserKV()
	kv.Set("pid", lua.NewFunction(pidL))
	kv.Set("process", lua.NewFunction(processL))
	kv.Set("switch", lua.NewFunction(switchL))
	kv.Set("listen_snapshot", lua.NewFunction(newListenSnapshotL))
	kv.Set("inode", lua.NewFunction(indexL))

	xEnv.Set("ss",
		lua.NewExport("vela.ss.export",
			lua.WithTable(kv),
			lua.WithFunc(call)))
}
