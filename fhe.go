package main

/*
#cgo LDFLAGS: -L. -lfhe_bridge -L/usr/local/lib -lseal-4.3 -lstdc++ -lm
#cgo CFLAGS: -I.
#include "fhe_c.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

type FHEContext struct {
	ptr unsafe.Pointer
	mu  sync.Mutex
}

var (
	globalCtx *FHEContext
	once      sync.Once
)

func InitFHE() error {
	var initErr error
	once.Do(func() {
		ptr := C.fhe_init()
		if ptr == nil {
			initErr = errors.New("failed to initialize FHE context")
			return
		}
		globalCtx = &FHEContext{ptr: ptr}
	})
	return initErr
}

func GetFHE() *FHEContext {
	return globalCtx
}

func (ctx *FHEContext) Encrypt(value int64) (string, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cResult := C.fhe_encrypt(ctx.ptr, C.long(value))
	if cResult == nil {
		return "", errors.New("encryption failed")
	}
	defer C.fhe_free_string(cResult)

	return C.GoString(cResult), nil
}

func (ctx *FHEContext) Decrypt(encrypted string) (int64, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cStr := C.CString(encrypted)
	defer C.free(unsafe.Pointer(cStr))

	result := C.fhe_decrypt(ctx.ptr, cStr)
	return int64(result), nil
}

func (ctx *FHEContext) Add(enc1, enc2 string) (string, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cEnc1 := C.CString(enc1)
	defer C.free(unsafe.Pointer(cEnc1))
	cEnc2 := C.CString(enc2)
	defer C.free(unsafe.Pointer(cEnc2))

	cResult := C.fhe_add(ctx.ptr, cEnc1, cEnc2)
	if cResult == nil {
		return "", errors.New("addition failed")
	}
	defer C.fhe_free_string(cResult)

	return C.GoString(cResult), nil
}

func (ctx *FHEContext) Multiply(enc1, enc2 string) (string, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cEnc1 := C.CString(enc1)
	defer C.free(unsafe.Pointer(cEnc1))
	cEnc2 := C.CString(enc2)
	defer C.free(unsafe.Pointer(cEnc2))

	cResult := C.fhe_multiply(ctx.ptr, cEnc1, cEnc2)
	if cResult == nil {
		return "", errors.New("multiplication failed")
	}
	defer C.fhe_free_string(cResult)

	return C.GoString(cResult), nil
}

func (ctx *FHEContext) NoiseBudget(encrypted string) (int, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cStr := C.CString(encrypted)
	defer C.free(unsafe.Pointer(cStr))

	budget := C.fhe_noise_budget(ctx.ptr, cStr)
	if budget < 0 {
		return 0, errors.New("failed to get noise budget")
	}
	return int(budget), nil
}

func (ctx *FHEContext) ZANSAdd(encrypted string, iterations int) (string, error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	cStr := C.CString(encrypted)
	defer C.free(unsafe.Pointer(cStr))

	cResult := C.fhe_zans_add(ctx.ptr, cStr, C.int(iterations))
	if cResult == nil {
		return "", errors.New("ZANS addition failed")
	}
	defer C.fhe_free_string(cResult)

	return C.GoString(cResult), nil
}

func (ctx *FHEContext) Close() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.ptr != nil {
		C.fhe_free(ctx.ptr)
		ctx.ptr = nil
	}
}
