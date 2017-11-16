package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	//"byte"

	"github.com/tarm/goserial"
)

func readAll(b io.ReadWriteCloser, len int) ([]byte, error) {
	pack := make([]byte, len)
	readedLen := 0
	try := 0
	for {
		one_len, err := b.Read(pack[readedLen:])
		if err != nil {
			e, ok := err.(net.Error)
			if !ok || !e.Temporary() || try >= 3 {
				return []byte(""), err
			}
			try++
		}

		readedLen = readedLen + one_len

		if readedLen == len {
			break
		}
	}
	return pack, nil
}

func dataConversion(a []byte) int { //定义一个函数，单字节十六进制转十进制型
	var tmout uint8
	tmoutBuff := bytes.NewBuffer(a)
	binary.Read(tmoutBuff, binary.BigEndian, &tmout)
	return int(tmout)
}

func dataConversion1(a []byte) int { //定义一个函数，处理接收的双字节数据，转为十进制
	var tmout1 uint16
	tmoutBuff1 := bytes.NewBuffer(a)
	binary.Read(tmoutBuff1, binary.BigEndian, &tmout1)
	return int(tmout1)
}

func test1() {

	a := &serial.Config{Name: "COM4", Baud: 9600}
	b, err := serial.OpenPort(a)
	if err != nil {
		fmt.Println("端口打开失败", err)
		return
	}

	//节点地址
	var s0 []byte
	s0 = append(s0, 0x20)
	s0 = append(s0, 0x07)
	s0 = append(s0, 0xff)
	s0 = append(s0, 0xff)
	s0 = append(s0, 0xff)
	s0 = append(s0, 0xff)
	s0 = append(s0, 0x00)

	var SensorData0 = map[string]int{
		"head_0":    1,
		"lenth_0":   1,
		"address_0": 4,
	}

	w, err := b.Write(s0)
	if err != nil && w <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer0 := make([]byte, 64)
	buffer0, err = readAll(b, 6)
	if err != nil {
		fmt.Println(err)
		return
	}

	data0 := buffer0[0:]
	nowlen0 := 0
	headlen0 := SensorData0["head_0"]
	// headdata := data[nowlen:headlen]
	nowlen0 = nowlen0 + headlen0
	lenthlen0 := SensorData0["lenth_0"]

	nowlen0 = nowlen0 + lenthlen0
	addresslen0 := SensorData0["address_0"]
	adressdata0 := data0[nowlen0:(addresslen0 + nowlen0)]
	//number0 := dataConversions(adressdata)
	fmt.Printf("节点地址: ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%02x ", adressdata0[i]) //节点地址
	}

	fmt.Println()

	//节点类型
	var s1 []byte
	s1 = append(s1, 0x20)
	s1 = append(s1, 0x07)
	s1 = append(s1, 0xff)
	s1 = append(s1, 0xff)
	s1 = append(s1, 0xff)
	s1 = append(s1, 0xff)
	s1 = append(s1, 0x07)

	var SensorData1 = map[string]int{
		"head_1":     1,
		"lenth_1":    1,
		"address_1":  4,
		"function_1": 1,
		"type_1":     1,
	}

	w1, err := b.Write(s1)
	if err != nil && w1 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer1 := make([]byte, 64)
	buffer1, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}
	data1 := buffer1[0:]
	nowlen1 := 0
	headlen1 := SensorData1["head_1"]

	nowlen1 = nowlen1 + headlen1
	lenthlen1 := SensorData1["lenth_1"]

	nowlen1 = nowlen1 + lenthlen1
	addresslen1 := SensorData1["address_1"]

	nowlen1 = nowlen1 + addresslen1
	functioncodelen1 := SensorData1["function_1"]

	nowlen1 = nowlen1 + functioncodelen1
	typelen1 := SensorData1["type_1"]
	typedata1 := data1[nowlen1:(nowlen1 + typelen1)]
	number1 := dataConversion(typedata1)
	fmt.Printf("节点类型：")
	if number1 == 1 {
		fmt.Printf("中心节点")
	}
	if number1 == 2 {
		fmt.Printf("中继路由")

	}
	if number1 == 3 {
		fmt.Printf("终端节点")

	}
	// for i := 0; i < 8; i++ {
	// 	fmt.Printf("%02x ", buffer1[i]) //节点地址
	// }
	fmt.Println()

	//网络类型(注：漏掉了，再加上所以编排上，是按照最后一个编排的)
	var s15 []byte
	s15 = append(s15, 0x20)
	s15 = append(s15, 0x07)
	s15 = append(s15, 0xff)
	s15 = append(s15, 0xff)
	s15 = append(s15, 0xff)
	s15 = append(s15, 0xff)
	s15 = append(s15, 0x06)

	var SensorData15 = map[string]int{
		"head_15":     1,
		"lenth_15":    1,
		"address_15":  4,
		"function_15": 1,
		"type_15":     1,
	}

	w15, err := b.Write(s15)
	if err != nil && w15 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer15 := make([]byte, 64)
	buffer15, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}
	data15 := buffer15[0:]
	nowlen15 := 0
	headlen15 := SensorData15["head_15"]

	nowlen15 = nowlen15 + headlen15
	lenthlen15 := SensorData15["lenth_15"]

	nowlen15 = nowlen15 + lenthlen15
	addresslen15 := SensorData15["address_15"]

	nowlen15 = nowlen15 + addresslen15
	functioncodelen15 := SensorData15["function_15"]

	nowlen15 = nowlen15 + functioncodelen15
	typelen15 := SensorData15["type_15"]
	typedata15 := data15[nowlen15:(nowlen15 + typelen15)]
	number15 := dataConversion(typedata15)
	fmt.Printf("网络类型：")
	if number15 == 1 {
		fmt.Printf("网状网")
	}
	if number15 == 2 {
		fmt.Printf("星型网")

	}
	if number15 == 3 {
		fmt.Printf("链型网")

	}
	if number15 == 4 {
		fmt.Printf("对等网")

	}

	fmt.Println()

	//网络ID
	var s2 []byte
	s2 = append(s2, 0x20)
	s2 = append(s2, 0x07)
	s2 = append(s2, 0xff)
	s2 = append(s2, 0xff)
	s2 = append(s2, 0xff)
	s2 = append(s2, 0xff)
	s2 = append(s2, 0x05)

	var SensorData2 = map[string]int{
		"head_2":     1,
		"lenth_2":    1,
		"address_2":  4,
		"function_2": 1,
		"NETID_2":    2,
	}

	w2, err := b.Write(s2)
	if err != nil && w2 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer2 := make([]byte, 64)
	buffer2, err = readAll(b, 9)
	if err != nil {
		fmt.Println(err)
		return
	}
	data2 := buffer2[0:]
	nowlen2 := 0
	headlen2 := SensorData2["head_2"]

	nowlen2 = nowlen2 + headlen2
	lenthlen2 := SensorData2["lenth_2"]

	nowlen2 = nowlen2 + lenthlen2
	addresslen2 := SensorData2["address_2"]

	nowlen2 = nowlen2 + addresslen2
	functioncodelen2 := SensorData2["function_2"]

	nowlen2 = nowlen2 + functioncodelen2
	NETIDlen2 := SensorData2["NETID_2"]
	NETIDdata2 := data2[nowlen2:(nowlen2 + NETIDlen2)]

	fmt.Printf("网络ID: ")
	for i := 0; i < 2; i++ {
		fmt.Printf("%02x ", NETIDdata2[i]) //
	}
	fmt.Println()

	//无线频点
	var s3 []byte
	s3 = append(s3, 0x20)
	s3 = append(s3, 0x07)
	s3 = append(s3, 0xff)
	s3 = append(s3, 0xff)
	s3 = append(s3, 0xff)
	s3 = append(s3, 0xff)
	s3 = append(s3, 0x03)

	var SensorData3 = map[string]int{
		"head_3":           1,
		"lenth_3":          1,
		"address_3":        4,
		"function_3":       1,
		"frequencyPoint_3": 1,
	}

	w3, err := b.Write(s3)
	if err != nil && w3 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer3 := make([]byte, 64)
	buffer3, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data3 := buffer3[0:]
	nowlen3 := 0
	headlen3 := SensorData3["head_3"]

	nowlen3 = nowlen3 + headlen3
	lenthlen3 := SensorData3["lenth_3"]

	nowlen3 = nowlen3 + lenthlen3
	addresslen3 := SensorData3["address_3"]

	nowlen3 = nowlen3 + addresslen3
	functioncodelen3 := SensorData3["function_3"]

	nowlen3 = nowlen3 + functioncodelen3
	FQlen3 := SensorData3["frequencyPoint_3"]
	FQdata3 := data3[nowlen3:(nowlen3 + FQlen3)]
	fmt.Printf("无线频点: ")
	fmt.Printf("%02x ", FQdata3)
	fmt.Println()

	//数据编码
	var s4 []byte
	s4 = append(s4, 0x20)
	s4 = append(s4, 0x07)
	s4 = append(s4, 0xff)
	s4 = append(s4, 0xff)
	s4 = append(s4, 0xff)
	s4 = append(s4, 0xff)
	s4 = append(s4, 0x14)

	var SensorData4 = map[string]int{
		"head_4":         1,
		"lenth_4":        1,
		"address_4":      4,
		"function_4":     1,
		"DataEncoding_4": 1,
	}

	w4, err := b.Write(s4)
	if err != nil && w4 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer4 := make([]byte, 64)
	buffer4, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data4 := buffer4[0:]
	nowlen4 := 0
	headlen4 := SensorData4["head_4"]

	nowlen4 = nowlen4 + headlen4
	lenthlen4 := SensorData4["lenth_4"]

	nowlen4 = nowlen4 + lenthlen4
	addresslen4 := SensorData4["address_4"]

	nowlen4 = nowlen4 + addresslen4
	functioncodelen4 := SensorData4["function_4"]

	nowlen4 = nowlen4 + functioncodelen4
	DElen4 := SensorData4["DataEncoding_4"]
	DEdata4 := data4[nowlen4:(nowlen4 + DElen4)]
	number4 := dataConversion(DEdata4)
	fmt.Printf("数据编码：")
	if number4 == 1 {
		fmt.Printf("ASCII")
	}
	if number4 == 2 {
		fmt.Printf("HEX")
	}
	fmt.Println()
	//发送模式
	var s5 []byte
	s5 = append(s5, 0x20)
	s5 = append(s5, 0x07)
	s5 = append(s5, 0xff)
	s5 = append(s5, 0xff)
	s5 = append(s5, 0xff)
	s5 = append(s5, 0xff)
	s5 = append(s5, 0x08)

	var SensorData5 = map[string]int{
		"head_5":     1,
		"lenth_5":    1,
		"address_5":  4,
		"function_5": 1,
		"SendMod_5":  1,
	}

	w5, err := b.Write(s5)
	if err != nil && w5 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer5 := make([]byte, 64)
	buffer5, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data5 := buffer5[0:]
	nowlen5 := 0
	headlen5 := SensorData5["head_5"]

	nowlen5 = nowlen5 + headlen5
	lenthlen5 := SensorData5["lenth_5"]

	nowlen5 = nowlen5 + lenthlen5
	addresslen5 := SensorData5["address_5"]

	nowlen5 = nowlen5 + addresslen5
	functioncodelen5 := SensorData5["function_5"]

	nowlen5 = nowlen5 + functioncodelen5
	SMlen5 := SensorData5["SendMod_5"]
	SMdata5 := data5[nowlen5:(nowlen5 + SMlen5)]
	number5 := dataConversion(SMdata5)
	fmt.Printf("发送模式：")
	if number5 == 1 {
		fmt.Printf("广播")
	}
	if number5 == 2 {
		fmt.Printf("固定协议")
	}
	if number5 == 3 {
		fmt.Printf("SHUNCOM协议")
	}
	if number5 == 4 {
		fmt.Printf("MODBUS_RTU")
	}
	fmt.Println()

	//波特率
	var s6 []byte
	s6 = append(s6, 0x20)
	s6 = append(s6, 0x07)
	s6 = append(s6, 0xff)
	s6 = append(s6, 0xff)
	s6 = append(s6, 0xff)
	s6 = append(s6, 0xff)
	s6 = append(s6, 0x10)

	var SensorData6 = map[string]int{
		"head_6":     1,
		"lenth_6":    1,
		"address_6":  4,
		"function_6": 1,
		"BaudRate_6": 1,
	}

	w6, err := b.Write(s6)
	if err != nil && w6 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer6 := make([]byte, 64)
	buffer6, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data6 := buffer6[0:]
	nowlen6 := 0
	headlen6 := SensorData6["head_6"]

	nowlen6 = nowlen6 + headlen6
	lenthlen6 := SensorData6["lenth_6"]

	nowlen6 = nowlen6 + lenthlen6
	addresslen6 := SensorData6["address_6"]

	nowlen6 = nowlen6 + addresslen6
	functioncodelen6 := SensorData6["function_6"]

	nowlen6 = nowlen6 + functioncodelen6
	BRlen6 := SensorData6["BaudRate_6"]
	BRdata6 := data6[nowlen6:(nowlen6 + BRlen6)]
	number6 := dataConversion(BRdata6)
	fmt.Printf("波特率：")
	if number6 == 1 {
		fmt.Printf("1200")
	}
	if number6 == 2 {
		fmt.Printf("2400")
	}
	if number6 == 3 {
		fmt.Printf("4800")
	}
	if number6 == 4 {
		fmt.Printf("9600")
	}
	if number6 == 5 {
		fmt.Printf("19200")
	}
	if number6 == 6 {
		fmt.Printf("38400")
	}
	if number6 == 7 {
		fmt.Printf("57600")
	}
	if number6 == 8 {
		fmt.Printf("115200")
	}
	fmt.Println()

	//校验
	var s7 []byte
	s7 = append(s7, 0x20)
	s7 = append(s7, 0x07)
	s7 = append(s7, 0xff)
	s7 = append(s7, 0xff)
	s7 = append(s7, 0xff)
	s7 = append(s7, 0xff)
	s7 = append(s7, 0x11)

	var SensorData7 = map[string]int{
		"head_7":     1,
		"lenth_7":    1,
		"address_7":  4,
		"function_7": 1,
		"Check_7":    1,
	}
	w7, err := b.Write(s7)
	if err != nil && w7 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer7 := make([]byte, 64)
	buffer7, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data7 := buffer7[0:]
	nowlen7 := 0
	headlen7 := SensorData7["head_7"]

	nowlen7 = nowlen7 + headlen7
	lenthlen7 := SensorData7["lenth_7"]

	nowlen7 = nowlen7 + lenthlen7
	addresslen7 := SensorData7["address_7"]

	nowlen7 = nowlen7 + addresslen7
	functioncodelen7 := SensorData7["function_7"]

	nowlen7 = nowlen7 + functioncodelen7
	Checklen7 := SensorData7["Check_7"]
	Checkdata7 := data7[nowlen7:(nowlen7 + Checklen7)]
	number7 := dataConversion(Checkdata7)
	fmt.Printf("校验：")
	if number7 == 1 {
		fmt.Printf("None")
	}
	if number7 == 2 {
		fmt.Printf("Even")
	}
	if number7 == 3 {
		fmt.Printf("0dd")
	}
	if number7 == 4 {
		fmt.Printf("Mark")
	}
	if number7 == 5 {
		fmt.Printf("Space")
	}

	fmt.Println()

	// //数据位
	// var s8 []byte
	// s8 = append(s8, 0x20)
	// s8 = append(s8, 0x07)
	// s8 = append(s8, 0xff)
	// s8 = append(s8, 0xff)
	// s8 = append(s8, 0xff)
	// s8 = append(s8, 0xff)
	// s8 = append(s8, 0x12)

	// var SensorData8 = map[string]int{
	// 	"head_8":     1,
	// 	"lenth_8":    1,
	// 	"address_8":  4,
	// 	"function_8": 1,
	// 	"Data_8":     1,
	// }

	// w8, err := b.Write(s8)
	// if err != nil && w8 <= 0 {
	// 	fmt.Println("内部录入地址命令错误")
	// }

	// buffer8 := make([]byte, 64)
	// buffer8, err = readAll(b, 8)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// data8 := buffer8[0:]
	// nowlen8 := 0
	// headlen8 := SensorData8["head_8"]

	// nowlen8 = nowlen8 + headlen8
	// lenthlen8 := SensorData8["lenth_8"]

	// nowlen8 = nowlen8 + lenthlen8
	// addresslen8 := SensorData8["address_8"]

	// nowlen8 = nowlen8 + addresslen8
	// functioncodelen8 := SensorData8["function_8"]

	// nowlen8 = nowlen8 + functioncodelen8
	// Datalen8 := SensorData8["Data_8"]
	// Datadata8 := data8[nowlen8:(nowlen8 + Datalen8)]
	// number8 := dataConversion(Datadata8)
	// fmt.Println(number8)
	// fmt.Printf("数据位：")
	// if number8 == 8 {
	// 	fmt.Printf("8+0+1")
	// }
	// if number8 == 2 {
	// 	fmt.Printf("8+0+1")
	// }
	// if number8 == 3 {
	// 	fmt.Printf("8+1+1")
	// }
	// if number8 == 4 {
	// 	fmt.Printf("8+0+2")
	// }
	// fmt.Println()

	//数据源址
	var s9 []byte
	s9 = append(s9, 0x20)
	s9 = append(s9, 0x07)
	s9 = append(s9, 0xff)
	s9 = append(s9, 0xff)
	s9 = append(s9, 0xff)
	s9 = append(s9, 0xff)
	s9 = append(s9, 0x15)

	var SensorData9 = map[string]int{
		"head_9":        1,
		"lenth_9":       1,
		"address_9":     4,
		"function_9":    1,
		"DataAddress_9": 1,
	}

	w9, err := b.Write(s9)
	if err != nil && w9 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer9 := make([]byte, 64)
	buffer9, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data9 := buffer9[0:]
	nowlen9 := 0
	headlen9 := SensorData9["head_9"]

	nowlen9 = nowlen9 + headlen9
	lenthlen9 := SensorData9["lenth_9"]

	nowlen9 = nowlen9 + lenthlen9
	addresslen9 := SensorData9["address_9"]

	nowlen9 = nowlen9 + addresslen9
	functioncodelen9 := SensorData9["function_9"]

	nowlen9 = nowlen9 + functioncodelen9
	DAlen9 := SensorData9["DataAddress_9"]
	DAdata9 := data9[nowlen9:(nowlen9 + DAlen9)]
	number9 := dataConversion(DAdata9)
	fmt.Printf("数据源址：")
	if number9 == 1 {
		fmt.Printf("不输出")
	}
	if number9 == 2 {
		fmt.Printf("输出")
	}

	fmt.Println()

	//发射功率
	var s10 []byte
	s10 = append(s10, 0x20)
	s10 = append(s10, 0x07)
	s10 = append(s10, 0xff)
	s10 = append(s10, 0xff)
	s10 = append(s10, 0xff)
	s10 = append(s10, 0xff)
	s10 = append(s10, 0x04)

	var SensorData10 = map[string]int{
		"head_10":         1,
		"lenth_10":        1,
		"address_10":      4,
		"function_10":     1,
		"TansmitPower_10": 1,
	}

	w10, err := b.Write(s10)
	if err != nil && w10 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer10 := make([]byte, 64)
	buffer10, err = readAll(b, 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	data10 := buffer10[0:]
	nowlen10 := 0
	headlen10 := SensorData10["head_10"]

	nowlen10 = nowlen10 + headlen10
	lenthlen10 := SensorData10["lenth_10"]

	nowlen10 = nowlen10 + lenthlen10
	addresslen10 := SensorData10["address_10"]

	nowlen10 = nowlen10 + addresslen10
	functioncodelen10 := SensorData10["function_10"]

	nowlen10 = nowlen10 + functioncodelen10
	TSlen10 := SensorData10["TansmitPower_10"]
	TSdata10 := data10[nowlen10:(nowlen10 + TSlen10)]
	number10 := dataConversion(TSdata10)
	fmt.Printf("发射功率：")
	if number10 == 1 {
		fmt.Printf("最小")
	}
	if number10 == 2 {
		fmt.Printf("中等")
	}
	if number10 == 3 {
		fmt.Printf("最大")
	}

	fmt.Println()

	// //休眠控制
	// var s11 []byte
	// s11 = append(s11, 0x20)
	// s11 = append(s11, 0x07)
	// s11 = append(s11, 0xff)
	// s11 = append(s11, 0xff)
	// s11 = append(s11, 0xff)
	// s11 = append(s11, 0xff)
	// s11 = append(s11, 0x16)

	// w11, err := b.Write(s11)
	// if err != nil && w11 <= 0 {
	// 	fmt.Println("内部录入地址命令错误")
	// }

	// buffer11 := make([]byte, 64)
	// buffer11, err = readAll(b, 8)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for i := 0; i < 8; i++ {
	// 	fmt.Printf("%02x ", buffer11[i])
	// }
	// fmt.Println()

	// //休眠时间
	// var s12 []byte
	// s12 = append(s12, 0x20)
	// s12 = append(s12, 0x07)
	// s12 = append(s12, 0xff)
	// s12 = append(s12, 0xff)
	// s12 = append(s12, 0xff)
	// s12 = append(s12, 0xff)
	// s12 = append(s12, 0x19)

	// w12, err := b.Write(s12)
	// if err != nil && w12 <= 0 {
	// 	fmt.Println("内部录入地址命令错误")
	// }

	// buffer12 := make([]byte, 64)
	// buffer12, err = readAll(b, 9)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for i := 0; i < 9; i++ {
	// 	fmt.Printf("%02x ", buffer12[i])
	// }
	// fmt.Println()

	// //工作时间
	// var s13 []byte
	// s13 = append(s13, 0x20)
	// s13 = append(s13, 0x07)
	// s13 = append(s13, 0xff)
	// s13 = append(s13, 0xff)
	// s13 = append(s13, 0xff)
	// s13 = append(s13, 0xff)
	// s13 = append(s13, 0x18)

	// w13, err := b.Write(s13)
	// if err != nil && w13 <= 0 {
	// 	fmt.Println("内部录入地址命令错误")
	// }

	// buffer13 := make([]byte, 64)
	// buffer13, err = readAll(b, 8)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for i := 0; i < 8; i++ {
	// 	fmt.Printf("%02x ", buffer13[i])
	// }
	// fmt.Println()

	//主动上报
	var s14 []byte
	s14 = append(s14, 0x20)
	s14 = append(s14, 0x07)
	s14 = append(s14, 0xff)
	s14 = append(s14, 0xff)
	s14 = append(s14, 0xff)
	s14 = append(s14, 0xff)
	s14 = append(s14, 0x20)

	var SensorData14 = map[string]int{
		"head_14":            1,
		"lenth_14":           1,
		"address_14":         4,
		"function_14":        1,
		"ActiveReporting_14": 2,
	}

	w14, err := b.Write(s14)
	if err != nil && w14 <= 0 {
		fmt.Println("内部录入地址命令错误")
	}

	buffer14 := make([]byte, 64)
	buffer14, err = readAll(b, 9)
	if err != nil {
		fmt.Println(err)
		return
	}
	data14 := buffer14[0:]
	nowlen14 := 0
	headlen14 := SensorData14["head_14"]

	nowlen14 = nowlen14 + headlen14
	lenthlen14 := SensorData14["lenth_14"]

	nowlen14 = nowlen14 + lenthlen14
	addresslen14 := SensorData14["address_14"]

	nowlen14 = nowlen14 + addresslen14
	functioncodelen14 := SensorData14["function_14"]

	nowlen14 = nowlen14 + functioncodelen14
	ARlen14 := SensorData14["ActiveReporting_14"]
	ARdata14 := data14[nowlen14:(nowlen14 + ARlen14)]
	number14 := dataConversion1(ARdata14)
	fmt.Printf("主动上报：")
	fmt.Println(number14)

	//IO功能

	//量程模式
}

func main() {

	test1()
}
