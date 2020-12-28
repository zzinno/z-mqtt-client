// @Title  msg.go
// @Description  消息结构体
// @Author   loveward  2020/12/28 12:35
package client

type RequestMsg struct {
	MsgId     string
	Data      Data
	FromTopic string
}

type Data struct {
	Key     []byte
	Columns string
}

type RespondMsg struct {
	RequestMsg     RequestMsg
	Code           float64
	Data           []byte
	ReplyFromTopic string
	Err            error
}

/* @Description
   此处收到的msg
   如果type为0，就是请求
   如果type为1，就是应答
*/
type ZMsg struct {
	MsgType    int
	MsgContent []byte
}
