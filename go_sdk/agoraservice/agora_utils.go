package agoraservice

// #cgo CFLAGS: -I${SRCDIR}/../../agora_sdk/include/c/api2 -I${SRCDIR}/../../agora_sdk/include/c/base
// #include <stdlib.h>
// #include "agora_parameter.h"
import (
    "sync"
    "time"
    "bytes"
)

// AudioConsumer provides utility functions for the Agora SDK.


// AudioConsumer handles PCM data consumption and sending
type AudioConsumer struct {
    mu              sync.Mutex
    startTime       int64  // in ms
    buffer          *bytes.Buffer
    consumedPackets int
    pcmSender       *AudioPcmDataSender
    frame           *AudioFrame
    
    // Audio parameters
    bytesPerFrame    int
    samplesPerChannel int
    
    // State
    isInitialized bool
}

// NewAudioConsumer creates a new AudioConsumer instance
func NewAudioConsumer(pcmSender *AudioPcmDataSender, sampleRate int, channels int) *AudioConsumer {
    if pcmSender == nil {
        return nil
    }
    
    bytesPerFrame := (sampleRate / 100) * channels * 2 // 2 bytes per sample
    
    consumer := &AudioConsumer{
        buffer:           bytes.NewBuffer(make([]byte, 0, bytesPerFrame*20)), // Pre-allocate buffer
        pcmSender:       pcmSender,
        frame: &AudioFrame{
            SamplesPerSec:        sampleRate,
            Channels:         channels,
            BytesPerSample:   2,
            Buffer:            make([]byte, 0, bytesPerFrame), // Pre-allocate frame buffer
        },
        bytesPerFrame:    bytesPerFrame,
        samplesPerChannel: sampleRate / 100 * channels,
        isInitialized:    true,
    }

    return consumer
}

// PushPCMData adds PCM data to the buffer
func (ac *AudioConsumer) PushPCMData(data []byte) {
    if !ac.isInitialized || data == nil || len(data) == 0 {
        return
    }
    
    ac.mu.Lock()
	defer ac.mu.Unlock()
    ac.buffer.Write(data)
}

// reset resets the consumer's timing state
func (ac *AudioConsumer) reset() {
    if !ac.isInitialized {
        return
    }
    
    
    ac.startTime = (time.Now().UnixMilli())
    ac.consumedPackets = 0
}

// Consume processes and sends audio data
func (ac *AudioConsumer) Consume() int {
    if !ac.isInitialized {
        return -1
    }
    
    now := time.Now().UnixMilli()
    elapsedTime := now - ac.startTime
    expectedTotalPackets := int(elapsedTime / 10) //change type from 
    toBeSentPackets := expectedTotalPackets - ac.consumedPackets

    dataLen := ac.buffer.Len()
    
    // Handle underflow
    if toBeSentPackets > 18 && dataLen/ac.bytesPerFrame < 18 {
        return -2 // should wait for more data
    }
    
    // Reset state if necessary
    if toBeSentPackets > 18 {
        ac.reset()
        toBeSentPackets = min(18, dataLen/ac.bytesPerFrame)
        ac.consumedPackets = (-toBeSentPackets)
    }
    
    // Calculate actual packets to send
    actualPackets := min(toBeSentPackets, dataLen/ac.bytesPerFrame)
    if actualPackets < 1 {
        return -3
    }
    
    // Prepare and send frame
    bytesToSend := ac.bytesPerFrame * actualPackets
    
    ac.mu.Lock()
    frameData := make([]byte, bytesToSend)
    n, _ := ac.buffer.Read(frameData)
    ac.mu.Unlock()
    
    if n > 0 {
        ac.frame.Buffer = frameData[:n]
        
        ac.frame.SamplesPerChannel = ac.samplesPerChannel *actualPackets
        
        ac.consumedPackets += actualPackets
        
        ret := ac.pcmSender.SendAudioPcmData(ac.frame)
        return ret
    }
    
    return -5
}

// Len returns the current buffer length
func (ac *AudioConsumer) Len() int {
    return ac.buffer.Len()
}

// Clear empties the buffer
func (ac *AudioConsumer) Clear() {
    ac.mu.Lock()
    defer ac.mu.Unlock()
    ac.buffer.Reset()
}

// Release frees resources
func (ac *AudioConsumer) Release() {
    if !ac.isInitialized {
        return
    }
    
    ac.isInitialized = false
    
    ac.mu.Lock()
    defer ac.mu.Unlock()
    
    // Clear references to allow GC
    ac.buffer = nil
    ac.frame = nil
    ac.pcmSender = nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}