package cmd

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"

	lapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	vapi "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/sirupsen/logrus"
)

func newHandler(ctx context.Context, vAPI vapi.FullNode, lAPI lapi.FullNode) *handler {
	h := &handler{
		ctx: ctx,

		vAPI: apiInfo{
			rv: reflect.ValueOf(vAPI),
			rt: reflect.TypeOf(vAPI),
		},
		lAPI: apiInfo{
			rv: reflect.ValueOf(lAPI),
			rt: reflect.TypeOf(lAPI),
		},

		receiver: make(chan *req, 20),
	}

	go h.start()

	return h
}

type handler struct {
	ctx context.Context

	vAPI apiInfo
	lAPI apiInfo

	receiver chan *req
}

type apiInfo struct {
	rv reflect.Value
	rt reflect.Type
}

func (h *handler) start() {
	for {
		select {
		case <-h.ctx.Done():
			logrus.Warn("context done, stop handler req")
			return
		case r := <-h.receiver:
			go func() {
				r.err <- h.call(r)
				close(r.err)
			}()
		}
	}
}

func (h *handler) call(r *req) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Fatalf("call %s panic: %v\n %s", r.methodName, err, string(debug.Stack()))
		}
	}()
	logrus.Debugf("start handler compare %v", r.methodName)
	defer func() {
		logrus.Debugf("end handler compare %v", r.methodName)
	}()
	vm, ok := h.vAPI.rv.Type().MethodByName(r.methodName)
	if !ok {
		return fmt.Errorf("not found method %s", r.methodName)
	}
	lm, ok := h.lAPI.rv.Type().MethodByName(r.methodName)
	if !ok {
		return fmt.Errorf("not found method %s", r.methodName)
	}

	inParams := make([]reflect.Value, 0, len(r.in))
	inParams2 := make([]reflect.Value, 0, len(r.in))
	for i, param := range r.in {
		v := reflect.ValueOf(param)
		inParams = append(inParams, v)
		// The first parameter is usually context.Context
		if i == 0 {
			inParams2 = append(inParams2, v)
			continue
		}
		inParams2 = append(inParams2, reflect.ValueOf(tryConvertParam(param)))
	}

	var vRes, lRes []reflect.Value
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Fatalf("call %s panic: %v\n %s", r.methodName, err, string(debug.Stack()))
			}
		}()
		defer wg.Done()
		vRes = vm.Func.Call(append([]reflect.Value{h.vAPI.rv}, inParams...))
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Fatalf("call %s panic: %v\n %s", r.methodName, err, string(debug.Stack()))
			}
		}()
		defer wg.Done()
		lRes = lm.Func.Call(append([]reflect.Value{h.lAPI.rv}, inParams2...))
	}()
	wg.Wait()

	if len(vRes) == 0 {
		return tryAppendParamsAndError(h.handleError(vRes[0], lRes[0]), inParams)
	}

	if len(vRes) == 2 {
		if err := h.handleError(vRes[1], lRes[1]); err != nil {
			if r.expectCallAPIError {
				logrus.Infof("call %s : %v", r.methodName, err)
				return nil
			}
			return tryAppendParamsAndError(err, inParams)
		}
	}
	logrus.Tracef("call %s result: \n%+v\n%+v", r.methodName, vRes[0].Interface(), lRes[0].Interface())

	if r.resultChecker != nil {
		return tryAppendParamsAndError(r.resultChecker(lRes[0].Interface(), vRes[0].Interface()), inParams)
	}

	return tryAppendParamsAndError(checkByJSON(vRes[0].Interface(), lRes[0].Interface()), inParams)
}

// todo: not check each param
func tryConvertParam(param interface{}) interface{} {
	key, ok := param.(types.TipSetKey)
	if ok {
		return toLoutsTipsetKey(types.TipSetKey(key))
	}
	msg, ok := param.(*types.Message)
	if ok {
		return toLotusMsg(msg)
	}
	num, ok := param.(types.EthUint64)
	if ok {
		return ethtypes.EthUint64(num)
	}
	hash, ok := param.(types.EthHash)
	if ok {
		return ethtypes.EthHash(hash)
	}
	ptrHash, ok := param.(*types.EthHash)
	if ok {
		return (*ethtypes.EthHash)(ptrHash)
	}
	addr, ok := param.(types.EthAddress)
	if ok {
		return ethtypes.EthAddress(addr)
	}
	call, ok := param.(types.EthCall)
	if ok {
		return toLotusEthCall(call)
	}
	bytes, ok := param.(types.EthBytes)
	if ok {
		return ethtypes.EthBytes(bytes)
	}
	msgMatch, ok := param.(*types.MessageMatch)
	if ok {
		return toLotusEthMessageMatch(msgMatch)
	}
	bh, ok := param.(types.EthBlockNumberOrHash)
	if ok {
		return toLotusEthBlockNumberOrHash(bh)
	}

	return param
}

func (h *handler) handleError(vErr, lErr reflect.Value) error {
	v := vErr.Interface()
	l := lErr.Interface()

	if v != nil || l != nil {
		if v != nil && l != nil {
			return fmt.Errorf("venus and lotus all return error: \n%v\n%v", v, l)
		}
		return fmt.Errorf("venus error: %v, lotus error: %v", v, l)
	}

	return nil
}

func tryAppendParamsAndError(err error, params []reflect.Value) error {
	if err == nil {
		return nil
	}

	str := "params:"
	// skip context
	for i := 1; i < len(params); i++ {
		str += fmt.Sprintf("%d:%v", i, params[i])
		if i < len(params)-1 {
			str += ", "
		}
	}

	return fmt.Errorf("%v\n %s", err, str)
}

func (h *handler) send(r *req) {
	select {
	case <-h.ctx.Done():
		r.err <- h.ctx.Err()
		return
	default:
	}

	h.receiver <- r
}
