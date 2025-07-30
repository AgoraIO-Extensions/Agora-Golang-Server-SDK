package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include "agora_service.h"
import "C"
import "unsafe"

type MediaNodeFactory struct {
	cFactory unsafe.Pointer
}

func newMediaNodeFactory() *MediaNodeFactory {
	factory := C.agora_service_create_media_node_factory(agoraService.service)
	if factory == nil {
		return nil
	}
	return &MediaNodeFactory{
		cFactory: factory,
	}
}

func (factory *MediaNodeFactory) release() {
	if factory.cFactory == nil {
		return
	}
	C.agora_media_node_factory_destroy(factory.cFactory)
	factory.cFactory = nil
}
