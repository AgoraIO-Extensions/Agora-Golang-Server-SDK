package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <stdlib.h>
// #include "agora_parameter.h"
import (
	"sync"
)
/*
* AudioVadManager
* RegisterAudioFrameObserver::RegisterAudioFrameObserver(observer, configure)
然后保持configure 在Connection中，并对connection中做VadManager的初始化工作
在call back的时候，从Connection中获取vadmanager，调用Process
修改call back回调的参数，增加2个
* 
*/

type AudioVadManager struct {
	vadInstance sync.Map  // only access inside
	isInitialized bool  // only access inside
	vadConfigure *AudioVadConfigV2
}
func NewAudioVadManager(config *AudioVadConfigV2) *AudioVadManager {
	return &AudioVadManager{
		isInitialized: true,
		vadConfigure: config,
		vadInstance: sync.Map{},
	}
}
func (m *AudioVadManager) makeKey(channel string, uid string) string {
	return channel + uid
}
func (m *AudioVadManager) Process(channel string, uid string, frame *AudioFrame) (*AudioFrame, VadState) {

	// 1. make key
	key := m.makeKey(channel, uid)
	// 2. get instance from map
	vad, ok := m.vadInstance.Load(key)
	if !ok {
		// 2.1. create instance
		vad = NewAudioVadV2(m.vadConfigure)
		
		// and add to map
		m.vadInstance.Store(key, vad)
	}
	
	// 3. do process
	return vad.(*AudioVadV2).Process(frame)
}
func (m *AudioVadManager) Release() {
	
	if !m.isInitialized {
		return
	}
	m.isInitialized = false
	m.vadConfigure = nil
	// 释放所有的 vad 实例
	m.vadInstance.Range(func(key, value interface{}) bool {
		if vad, ok := value.(*AudioVadV2); ok {
			vad.Release() // 释放资源
			return true
		}
		return false
	})
	
	m.vadInstance = sync.Map{}

}
