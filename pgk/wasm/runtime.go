package wasm

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/Stolkerve/kappa/pgk/types"
	"github.com/fxamacker/cbor/v2"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func RunFunction(functionWasm []byte, ctx context.Context, req *types.RequestWrapper) (*types.RuntimeResult, error) {
	var resWrapper types.ResponseWrapper
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	var t1 = time.Now()

	runtimeConfig := wazero.NewRuntimeConfig().
		WithMemoryLimitPages(50)
	runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	reqBuff, err := cbor.Marshal(*req)
	if err != nil {
		return nil, err
	}

	_, err = runtime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(
		func(_ context.Context, m api.Module) uint32 {
			return uint32(len(reqBuff))
		}).Export("getEncodedRequestSize").
		NewFunctionBuilder().WithFunc(
		func(_ context.Context, m api.Module, pointer uint32) {
			m.Memory().Write(pointer, reqBuff)
		}).Export("readEncodedRequestToPointer").
		NewFunctionBuilder().WithFunc(
		func(_ context.Context, m api.Module, pointer uint32, size uint32) {
			resEncodedBuf, _ := m.Memory().Read(pointer, size)
			cbor.Unmarshal(resEncodedBuf, &resWrapper)
		}).Export("writeResponseFromPointer").
		Instantiate(ctx)

	if err != nil {
		return nil, err
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	modConfig := wazero.NewModuleConfig().
		WithStdout(stdout).
		WithStderr(stderr)

	endExecutionChannel := make(chan struct {
		res types.RuntimeResult
		err error
	}, 1)

	go func(endExecutionChannel chan struct {
		res types.RuntimeResult
		err error
	}) {
		mod, err := runtime.InstantiateWithConfig(ctx, functionWasm, modConfig)
		defer mod.Close(ctx)
		endExecutionChannel <- struct {
			res types.RuntimeResult
			err error
		}{
			res: types.RuntimeResult{
				ResponseWrapper: resWrapper,
				Duration:        time.Duration(time.Now().Sub(t1).Milliseconds()),
				Stdout:          stdout.String(),
				Memory:          mod.Memory().Size(),
				Stderr:          stderr.String(),
			},
			err: err,
		}
	}(endExecutionChannel)

	select {
	case r := <-endExecutionChannel:
		{
			runtime.CloseWithExitCode(ctx, 0)
			if r.err != nil {
				return nil, err
			}
			return &r.res, nil
		}
	case <-time.After(10 * time.Second):
		{
			runtime.CloseWithExitCode(ctx, 2)
			return nil, errors.New("Timeout")
		}
	}
}
