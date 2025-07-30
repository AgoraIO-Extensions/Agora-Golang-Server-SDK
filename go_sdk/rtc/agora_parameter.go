package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <stdlib.h>
// #include "agora_parameter.h"
import "C"
import "unsafe"

type AgoraParameter struct {
	cParameter unsafe.Pointer
}

func (p *AgoraParameter) SetParameters(params string) int {
	if p.cParameter == nil {
		return -1
	}
	cParameters := C.CString(params)
	defer C.free(unsafe.Pointer(cParameters))
	return int(C.agora_parameter_set_parameters(p.cParameter, cParameters))
}
