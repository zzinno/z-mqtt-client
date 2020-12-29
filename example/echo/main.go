package main

import (
	"fmt"
	"github.com/zzinno/z-mqtt-client/client"
)

func main() {
	c := client.ZClient{}
	err := c.New(client.Config{
		Port:      1883,
		Broker:    "47.100.242.110",
		ClientID:  "testClient",
		Username:  "dev",
		Password:  "123123",
		SelfTopic: "testCli",
		CallBack: func(msg *client.RequestMsg) (*[]byte, error) {
			fmt.Println(msg.MsgId)
			a := []byte("3333")
			return &a, nil
		},
		Logger: nil,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = c.SendMsgAndWaitReply(client.Data{
		Key:     []byte("222"),
		Columns: "222",
	}, "testCli", 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Loop()
}
