package message

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xuyp1991/libp2p_exercise/p2pnet"
	"log"
	"time"
)

var (
	msgStore map[string]MessageInfo
	Owner    string
)

// //构造一个结构体,存储消息,校验消息有没有接收过
type MessageInfo struct {
	Timestamp string
	Hash      string
	From      string
	To        string
	Data      []byte
}

func calculateHash(msg MessageInfo) string {
	record := msg.Timestamp + msg.From + msg.To + string(msg.Data)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func NewMsgInfo(to string, data []byte) MessageInfo {
	var result MessageInfo
	t := time.Now()

	result.From = Owner
	result.To = to
	result.Data = data
	result.Timestamp = t.String()
	result.Hash = calculateHash(result)
	return result
}

func InitMsgStore(owner string, outChan <-chan string) {
	msgStore = make(map[string]MessageInfo)
	Owner = owner
	go ListenData(outChan)
}

func (m MessageInfo) SendData() {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	p2pnet.SendData(string(bytes))
}

func UnmarshalMsgInfo(data string) (MessageInfo, bool) {
	var msgInfo MessageInfo
	if err := json.Unmarshal([]byte(data), &msgInfo); err != nil {
		log.Fatal(err)
	}

	if _, ok := msgStore[msgInfo.Hash]; ok {
		return msgInfo, false
	}

	msgStore[msgInfo.Hash] = msgInfo
	fmt.Printf("\x1b[32m%s\x1b[0m> ", data)
	return msgInfo, true
}

func ListenData(outChan <-chan string) {
	for {
		select {
		case outdata := <-outChan:
			UnmarshalMsgInfo(outdata)
		}
	}
}
