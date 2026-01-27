package agoraservice

import (
	"bytes"
	"fmt"
)

// Addsei 创建Agora格式的SEI数据
// ref to sei agora standard: https://confluence.agoralab.co/pages/viewpage.action?pageId=1418953055
/*
For H.265, the NAL Unit Header is 2 bytes long:
NAL Unit Header:
. Forbidden bit (1 bit): Usually 0, indicating the data contains no errors.
. NAL Unit Type (6 bits): Identifies the type of the NAL unit, such as VCL (Video Coding Layer) or non-VCL data.
. Layer ID (6 bits): Used to identify the layer to which the NAL unit belongs.
. Temporal ID (3 bits): Used to identify the temporal hierarchy of the NAL unit.
For H.264, the NAL Unit Header is 1 byte, typically 0x06.
0 Forbidden bit (F) (forbidden_bit)​
	1比特禁止位。通常为0，表示NAL单元无错误。当网络传输中发现错误时，可将其置1，以便接收方丢弃该单元 。


1-2 NRI (nal_ref_idc)​
2比特重要性指示。取值00到11，指示当前NALU的重要性。值越大，重要性越高。例如，参考帧的片、SPS、PPS等关键单元必须大于0；非参考帧或SEI等信息可设为00，解码器在负荷过高时可选择丢弃 。
	
3-7 Type (nal_unit_type)​
5比特 NALU类型。这是最关键的部分，定义了该NALU承载的数据内容。具体类型见下文详解, shoud be 0x06 for 264,
*/
func AddSEI(seiMsg []byte, isH264 bool) []byte {

	seiMsgLength := len(seiMsg)
	result := make([]byte, 1050)
	i := 0
	// 起始码: 0x00 0x00 0x00 0x01
	result[i] = 0x0
	i++
	result[i] = 0x0
	i++
	result[i] = 0x0
	i++
	result[i] = 0x1
	i++
	
	// NAL单元类型: for 264 1 byte; for 265 2 bytes;
	if isH264 {
		result[i] = 0x6
		i++
	} else {
		result[i] = 0x4e
		i++
		result[i] = 0x01
		i++
	}
	
	// SEI载荷类型: Agora特定
	//result[5] = SEI_H264_AGORA
	result[i] = 0x65
	i++
	
	// SEI载荷大小（可变长度编码）
	lengtmp := seiMsgLength
	
	// 对于大小 >= 255，使用多个0xFF字节
	for lengtmp >= 255 {
		result[i] = 0xff // 255
		i++
		lengtmp -= 255
	}
	
	// 写入剩余大小
	result[i] = byte(lengtmp)
	i++
	
	// 复制SEI消息数据
	copy(result[i:], seiMsg[:seiMsgLength])
	i += seiMsgLength
	
	// RBSP尾部比特
	result[i] = 0x80
	
	// 设置总SEI长度
	return result[:i+1]
}
// for 264, the nalu header is 1 byte; for 265, the nalu header is 2 bytes;
func is264NaluHeader(data []byte, index int) bool {
	header := data[index]
	header = header & 0x1f  // 3~7 bits, total 5 bits
	return header == 0x06 && data[index+1] == 0x65

}

// 2 types of nalu header for 265:
func is265NaluHeader(data []byte, index int) bool {
	header := data[index]
	header = header & 0x7E  // 1~7 bits, total 7 bits
	header = header >> 1  // 2~7 bits, total 6 bits
	return (header == 0x27 || header == 0x28) && data[index+2] == 0x65
}
/*
find sei data in the encoded data
bit info: sps pps sei 
for 264: sps pps sei
for 265: sps pps vps sei

sum(sps+pps+vps) < 20+50+20 < 200bytes

*/
/*
Parameters:
- encoded: the encoded data
- codecType: the codec type

Returns:
- seiData: the sei data
- seiLength: the length of the sei data
*/
func FindSEI(encoded []byte, codecType VideoCodecType) ([]byte, int) {
	if encoded == nil || len(encoded) == 0 {
		return nil, -1000
	}
	len := len(encoded)
	if len < 6 {
		return nil, -1001
	}

	startcode4 := []byte{0x00, 0x00, 0x00, 0x01}

	i := 0
	startIndex := 0
	endIndex := len - 6
	reverseSearch := false
	// end floag is 0x80
	if encoded[len-1] != 0x80 {
		// from begin to end
		startIndex = 0
		endIndex = len - 6 //max value
		endIndex = 200*4 //actual max value
		if endIndex > len - 6 {
			endIndex = len - 6
		}

	} else {
		// from end to begin: reverse search
		startIndex = len - 6
		endIndex = 6  //max value
		reverseSearch = true
	}
	// only support 264 and 265 now
	if codecType != VideoCodecTypeH264 && codecType != VideoCodecTypeH265 {
		return nil, -1003
	}
	fmt.Printf("startIndex: %d, endIndex: %d, reverseSearch: %t\n", startIndex, endIndex, reverseSearch)

	// find start flag
	data := encoded
	
	offset := 0
	find := false
	
	if reverseSearch == false {
		find = false
		for i = startIndex; i < endIndex; i++ {
			if bytes.Equal(data[i:i+4], startcode4) {
				if codecType == VideoCodecTypeH264 {
					if is264NaluHeader(data, i+4) {
						find = true
						offset = 6
						break
					}
				} else if codecType == VideoCodecTypeH265 {
					if is265NaluHeader(data, i+4) {
						find = true
						offset = 7
						break
					}
				}
			}
		}
	}

	// froce reverse search
	if reverseSearch == true {
		startIndex = len-7
		endIndex = 6
		find = false
		for i = startIndex; i >= endIndex; i-- {
			if bytes.Equal(data[i:i+4], startcode4) {
				if codecType == VideoCodecTypeH264 {
					if is264NaluHeader(data, i+4) {
						find = true
						offset = 6
						break
					}
				} else if codecType == VideoCodecTypeH265 {
					if is265NaluHeader(data, i+4) {
						find = true
						offset = 7
						break
					}
				}
			}
		}
	}
	
	if !find {
		return nil, -1004
	}
	// then parse the sei data
	index := i+offset
	seiLength := 0
	for true {
		if data[index] == 0xff {
			seiLength += 255
			index++
		} else {
			seiLength += int(data[index])
			break
		}
	}

	fmt.Printf("sei length: %d\n", seiLength)
	fmt.Printf("data[index+seiLength]: %d, index+seiLength+1: %d\n", data[index+seiLength], index+seiLength+1)
	
	// then check end-flag
	index += 1
	if data[index+seiLength] != 0x80 {
		return nil, -1003
	}

	// copy and return the sei data
	seiData := make([]byte, seiLength)
	copy(seiData, data[index:index+seiLength])
	return seiData,  seiLength
}


// InsertSEITo264 inserts SEI into H.264 encoded data
//
// Parameters:
//   - h264EncodedData: Input H.264 encoded data
//   - sei: SEI string to insert
//
// Returns:
//   - Output data with SEI inserted, or nil on error
func InsertSEIToEncodedData(encodedData []byte, sei []byte, codecType VideoCodecType) []byte {

	if codecType != VideoCodecTypeH264 && codecType != VideoCodecTypeH265 {
		return nil
	}

	if sei == nil || len(sei) == 0 {
		return encodedData
	}

	is264 := codecType == VideoCodecTypeH264


	if codecType == VideoCodecTypeH265 {
		is264 = false
	}

	seiNal := AddSEI(sei, is264)
	// try appedn
	result := append(encodedData, seiNal...)
	return result
}
