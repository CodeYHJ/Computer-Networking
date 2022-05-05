package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ICMP struct {
	Type       uint8
	Code       uint8
	CheckSum   uint16
	Identifier uint16
	Seq        uint16
}

func main() {
	ipaddr, err := net.ResolveIPAddr("ip", "baidu.com")
	if err != nil {
		fmt.Printf("ResolveIPAddr报错：%v \n", err)
	}
	fmt.Print(ipaddr.String(), "\n")

	for i := 1; i < 6; i++ {
		icmp := getICMP(uint16(i))
		if err = sendIcmp(icmp, ipaddr); err != nil {
			fmt.Print("fail sentICMP \n")
		}
		time.Sleep(2 * time.Second)
	}

}

func sendIcmp(icmp ICMP, ipaddr *net.IPAddr) error {
	connect, err := net.DialIP("ip4:icmp", nil, ipaddr)
	if err != nil {
		fmt.Printf("fail connect remote host: %s \n", err)
		return err
	}
	defer connect.Close()
	// icmp 结构体转换为byte
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	// icmp数据流 写入连接
	_, err = connect.Write(buffer.Bytes())

	if err != nil {
		fmt.Printf("fail connect write icmp : %s \n", err)
		return err
	}
	timeStart := time.Now()
	timeDeadline := time.Now().Add(time.Second * 2)
	// 连接读取数据截至时间
	connect.SetReadDeadline(timeDeadline)

	// 接收数据
	recv := make([]byte, 1024)
	receive, err := connect.Read(recv)

	if err != nil {
		fmt.Printf("fail connect receive data : %s \n", err)
		return err
	}

	timeEnd := time.Now()

	duration := timeEnd.Sub(timeStart).Nanoseconds() / 1e6

	fmt.Printf("%d bytes from %s: seq=%d ttl=%dms \n", receive, ipaddr.String(), icmp.Seq, duration)

	return err

}
func getICMP(seq uint16) ICMP {
	// init ICMP数据结构
	icmp := ICMP{Type: 8, Code: 0, CheckSum: 0, Seq: seq}

	// struct 2 buffer
	buffer := bytes.Buffer{}

	// write icmp struct in buffter
	binary.Write(&buffer, binary.BigEndian, icmp)

	// caculate icmp checksum
	icmp.CheckSum = caculateCheckSum(buffer.Bytes())

	buffer.Reset()

	return icmp

}
func caculateCheckSum(icmpByte []byte) uint16 {
	var (
		checksum uint32 = 0
		index    int    = 0
		length   int    = len(icmpByte)
	)
	for length > 1 {

		sum := uint32(icmpByte[index])<<8 + uint32(icmpByte[index+1])
		checksum += sum
		length -= 2
		index += 2

	}
	// 长度为基数
	if length > 0 {
		checksum += uint32(icmpByte[index])
	}

	checksum += (checksum >> 16)

	return uint16(^checksum)
}
