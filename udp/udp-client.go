package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

type Payload struct {
	SerialID uint32     `json:"serialid"`
	Channels [8]Channel `json:"channels"`
}

type Channel struct {
	ID    uint16 `json:"id"`
	Value int16  `json:"value"`
}

var (
	host = flag.String("host", "127.0.0.1", "host")
	port = flag.String("port", "9527", "port")

	data = [8]Channel{
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 4},
		{5, 5},
		{6, 6},
		{7, 7},
		{8, 8},
	}

	count = 0
)

func main() {
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		return
	}
	defer conn.Close()

	ticker1 := time.NewTicker(time.Millisecond * 50)

	ticker2 := time.NewTimer(time.Second * 60)

	for {
		select {
		case <-ticker1.C:

			payload := Payload{
				SerialID: uint32(count),
				Channels: data,
			}

			data, _ := payload.Marshal()

			buffs := new(bytes.Buffer)
			binary.Write(buffs, binary.BigEndian, data)

			_, err = conn.Write(buffs.Bytes())
			if err != nil {
				fmt.Println("failed:", err)
				return
			}

			count++
		case <-ticker2.C:
			return
		}
	}
}

func (p *Payload) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, p.SerialID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Channels)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
