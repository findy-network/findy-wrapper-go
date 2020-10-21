package c2go

//#cgo LDFLAGS: -lindy
//#include <stdio.h>
//#include <stdlib.h>
//#include "findy_glue.h"
import "C"
import (
	"unsafe"

	"github.com/findy-network/findy-wrapper-go/internal/ctx"
	"github.com/golang/glog"

	"github.com/findy-network/findy-wrapper-go/dto"
)

func buildResult(cmdHandle uint32, err int32, setter func(r *dto.Result)) {
	r := dto.Result{}
	if err == C.Success {
		setter(&r)
	} else {
		r.SetErrCode(int(err))
		errorMsg := C.GoString(C.findy_get_current_error())
		if glog.V(1) {
			glog.Info(errorMsg)
		}
		r.SetErrorJSON(errorMsg)
	}
	ch := ctx.CmdContext.Pop(cmdHandle, r)
	ch <- r
}

// MARK: C-callbacks, these are called from findy_glue.c

//export strStrUllHandler
func strStrUllHandler(cmdHandle uint32, err int32, cstr1, cstr2 *C.char, ull C.ulonglong) {
	str1 := C.GoString(cstr1)
	str2 := C.GoString(cstr2)
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetStr1(str1)
		r.SetStr2(str2)
		r.SetU64(uint64(ull))
	})
}

//export strStrStrHandler
func strStrStrHandler(cmdHandle uint32, err int32, cstr1, cstr2, cstr3 *C.char) {
	str1 := C.GoString(cstr1)
	str2 := C.GoString(cstr2)
	str3 := C.GoString(cstr3)
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetStr1(str1)
		r.SetStr2(str2)
		r.SetStr3(str3)
	})
}

//export strStrHandler
func strStrHandler(cmdHandle uint32, err int32, cstr1, cstr2 *C.char) {
	str1 := C.GoString(cstr1)
	str2 := C.GoString(cstr2)
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetStr1(str1)
		r.SetStr2(str2)
	})
}

//export strHandler
func strHandler(cmdHandle uint32, err int32, cstr *C.char) {
	str := C.GoString(cstr)
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetStr1(str)
	})
}

//export handleHandler
func handleHandler(cmdHandle uint32, err, handle int32) {
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetHandle(int(handle))
	})
}

//export boolHandler
func boolHandler(cmdHandle uint32, err int32, value C.bool) {
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetYes(bool(value))
	})
}

//export handleU32Handler
func handleU32Handler(cmdHandle uint32, err, handle int32, l uint32) {
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetHandle(int(handle))
		r.SetU64(uint64(l))
	})
}

//export handler
func handler(cmdHandle uint32, err int32) {
	buildResult(cmdHandle, err, func(r *dto.Result) {
		// nothing to set in this version
	})
}

//export u8ptrU32Handler
func u8ptrU32Handler(cmdHandle uint32, err int32, data *C.uchar, dataLen uint32) {
	bytes := C.GoBytes(unsafe.Pointer(data), C.int(dataLen))
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetBytes(bytes)
	})
}

//export strU8ptrU32Handler
func strU8ptrU32Handler(cmdHandle uint32, err int32, cstr *C.char, data *C.uchar, dataLen uint32) {
	str := C.GoString(cstr)
	bytes := C.GoBytes(unsafe.Pointer(data), C.int(dataLen))
	buildResult(cmdHandle, err, func(r *dto.Result) {
		r.SetStr1(str)
		r.SetBytes(bytes)
	})
}

// MARK: indy system SDK level helper functions

func FindySetRuntimeConfig(config string) int {
	configInC := C.CString(config)
	defer C.free(unsafe.Pointer(configInC))
	return int(C.indy_set_runtime_config(configInC))
}

// todo: consider moving these helper functions to internal package there are no need that they are visible
//		at least this moment. when we have rest api things might change.
// MARK: pool helper functions
func PoolCreateConfig(name, configJSON string) ctx.Channel {
	nameInC := C.CString(name)
	configJSONInC := C.CString(configJSON)
	defer C.free(unsafe.Pointer(nameInC))
	defer C.free(unsafe.Pointer(configJSONInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_create_pool_ledger_config(C.int(cmdHandle), nameInC, configJSONInC)
	return ch
}

func PoolList() ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_list_pools(C.int(cmdHandle))
	return ch
}

func PoolOpenLedger(name string) ctx.Channel {
	nameInC := C.CString(name)
	configInC := C.CString("{}")
	defer C.free(unsafe.Pointer(nameInC))
	defer C.free(unsafe.Pointer(configInC))
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_open_pool_ledger(C.int(cmdHandle), nameInC, configInC)
	return ch
}

func PoolCloseLedger(handle int) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_close_pool_ledger(C.int(cmdHandle), C.int(handle))
	return ch
}

func PoolSetProtocolVersion(version uint64) ctx.Channel {
	cmdHandle, ch := ctx.CmdContext.Push()
	C.findy_set_protocol_version(C.int(cmdHandle), C.ulonglong(version))
	return ch
}
