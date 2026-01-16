package agoraservice

// Addsei 创建Agora格式的SEI数据
// ref to sei agora standard: https://confluence.agoralab.co/pages/viewpage.action?pageId=1418953055
func Addsei(seiMsg []byte, isH264 bool) []byte {

	seiMsgLength := len(seiMsg)
	result := make([]byte, 1050)
	// 起始码: 0x00 0x00 0x00 0x01
	result[0] = 0x0
	result[1] = 0x0
	result[2] = 0x0
	result[3] = 0x1
	
	// NAL单元类型: SEI (0x06)
	if isH264 {
		result[4] = 0x6
	} else {
		result[4] = 0x27
	}
	
	// SEI载荷类型: Agora特定
	//result[5] = SEI_H264_AGORA
	result[5] = 101
	
	// SEI载荷大小（可变长度编码）
	lengtmp := seiMsgLength
	i := 6
	
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
	return result[:i]
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

	seiNal := Addsei(sei, is264)
	// try appedn
	result := append(encodedData, seiNal...)
	return result
}
