/**
 * @author uuxia
 * @date 15:58 2023/3/10
 * @description 通过拦截，处理数据，任务向下传递
 **/

package zcode

import (
	"github.com/aceld/zinx/ziface"
)

type LengthFieldFrameInterceptor struct {
	decoder ziface.ILengthField
}

func NewLengthFieldFrameInterceptor(maxFrameLength uint64, lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip int) *LengthFieldFrameInterceptor {
	return &LengthFieldFrameInterceptor{
		decoder: NewLengthFieldFrameDecoder(maxFrameLength, lengthFieldOffset, lengthFieldLength, lengthAdjustment, initialBytesToStrip),
	}
}

func (l *LengthFieldFrameInterceptor) Intercept(chain ziface.Chain) ziface.IcResp {
	req := chain.Request()

	if req == nil || l.decoder == nil {
		goto END
	}

	switch req.(type) {
	case ziface.IRequest:
		iRequest := req.(ziface.IRequest)
		iMessage := iRequest.GetMessage()

		if iMessage == nil {
			break
		}

		data := iMessage.GetData()

		bytebuffers := l.decoder.Decode(data)
		size := len(bytebuffers)
		if size == 0 { //半包，或者其他情况，任务就不要往下再传递了
			return nil
		}

		for i := 0; i < size; i++ {
			buffer := bytebuffers[i]
			if buffer == nil {
				continue
			}
			bufferSize := len(buffer)
			iMessage.SetData(buffer)
			iMessage.SetDataLen(uint32(bufferSize))

			if i < size-1 {
				chain.Proceed(chain.Request())
			}
		}
	}

END:
	return chain.Proceed(chain.Request())
}
