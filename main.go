package main

import (
	"bufio"
	// "context"
	// "crypto/rand"
	// "crypto/sha256"
	// "encoding/hex"
	// "encoding/json"
	"flag"
	"fmt"
	//   "io"
	"log"
	//    mrand "math/rand"
	"os"
	// "strconv"
	golog "github.com/ipfs/go-log"
	"strings"

	gologv2 "github.com/ipfs/go-log/v2"
	"github.com/xuyp1991/libp2p_exercise/message"
	"github.com/xuyp1991/libp2p_exercise/p2pnet"
)

func getStdInData() {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		if err != nil {
			log.Fatal(err)
		}
		newMsg := message.NewMsgInfo("all", []byte(sendData))
		newMsg.SendData()
	}
}

func main() {

	golog.SetAllLoggers(gologv2.LevelInfo)

	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	owner := flag.String("o", "", "owner of the node")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	outchan := p2pnet.StartNet(*listenF, *secio, *seed, *target)
	message.InitMsgStore(*owner, outchan)

	go getStdInData()

	select {} // hang forever
}
