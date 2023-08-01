package ss

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/vela"
)

func define(r vela.Router) {
	r.GET("/api/v1/arr/agent/ss", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {

		cmd := "-a -p"
		if !ctx.QueryArgs().Has("process") {
			cmd = "-a"
		}

		f, err := NewOptionFlag(cmd)
		if err != nil {
			return err
		}

		v := By(Flag(f))
		chunk, err := json.Marshal(v)
		if err != nil {
			return err
		}
		ctx.Write(chunk)
		return nil
	}))

	r.GET("/api/v1/arr/agent/listen", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		f, _ := NewOptionFlag("-l -p")
		v := By(Flag(f))
		chunk, err := json.Marshal(v.Sockets)
		if err != nil {
			return err
		}
		ctx.Write(chunk)
		return nil
	}))

	r.GET("/api/v1/arr/agent/ss/pid", xEnv.Then(func(ctx *fasthttp.RequestCtx) error {
		v := ctx.QueryArgs().Peek("pid")
		if len(v) == 0 {
			return fmt.Errorf("got process pid empty")
		}

		pid, err := auxlib.ToInt32E(string(v))
		if err != nil {
			return err
		}

		flag, err := NewOptionFlag("-a -p")
		if err != nil {
			return err
		}

		s := ByPID(pid, Flag(flag), Pid([]int32{pid}))
		if s.Err != nil {
			return err
		}

		chunk, err := json.Marshal(s)
		if err != nil {
			return err
		}

		ctx.Write(chunk)
		return nil
	}))
}
