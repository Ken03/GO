package main

import (
	//hi
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"net"
	"os"
)

func main() {

	//建立Scoket监听端口
	netListen, err := net.Listen("tcp", ":6002")
	CheckError(err)
	defer netListen.Close()

	Log("等待客户端连接")
	for {
		connect, err := netListen.Accept()
		if err != nil {
			continue
		}

		Log(connect.RemoteAddr().String(), "tcp连接成功")
		handleConnection(connect)

	}
}

func dataConversion(a []byte) int { //定义一个函数，单字节转为int型
	var tmout uint8
	tmoutBuff := bytes.NewBuffer(a)
	binary.Read(tmoutBuff, binary.BigEndian, &tmout)
	return int(tmout)
}

func dataConversions(a []byte) []int { //定义一个函数，数据可循环，且提取
	var res []int
	for i := 0; i < len(a); i++ {
		var tmout uint8
		var b []byte
		b = append(b, a[i])
		tmoutBuff := bytes.NewBuffer(b)
		binary.Read(tmoutBuff, binary.BigEndian, &tmout)
		res = append(res, int(tmout))
	}
	return res
	//number5 = res[0]*256 + res[1]

}

func dataConversion1(a []byte) int { //定义一个函数，处理双字节数据,转为Int
	var tmout1 uint16
	tmoutBuff1 := bytes.NewBuffer(a)
	binary.Read(tmoutBuff1, binary.BigEndian, &tmout1)
	return int(tmout1)
}

var BatteryData = map[string]int{
	"Address":                   6,  //mac地址
	"GroupNumber":               1,  //组序号
	"GroupVoltage":              2,  //组电压值
	"GroupCurrent":              2,  //组电流值
	"GroupState":                2,  //组状态
	"ReservedBytes":             6,  //保留字节
	"BatteryNumber":             1,  //组电池数量
	"BinaryContent":             22, //单节电池数据长度
	"MonitoringSerialNumber":    1,  //监测模块序号
	"Model":                     1,  //监测模块型号
	"BatteryStatus":             2,  //电池状态
	"BatteryInternalResistance": 2,  //电池内阻
	"BatteryVoltage":            2,  //电池电压
	"BatteryTemperature1":       2,  //电池1温度
	"BatteryTemperature2":       2,  //电池2温度
	"BatteryTemperature3":       2,  //电池3温度
	"BatteryTemperature4":       2,  //电池4温度
	"Ripple":                    2,  //纹波
	"ReservedBytes2":            4,  //保留字节

}

//处理链接

func handleConnection(connect net.Conn) {

	buffer := make([]byte, 4096)

	for {
		src, err := connect.Read(buffer)

		if err != nil {
			Log(connect.RemoteAddr().String(), "连接错误：", err)
			return
		}

		fmt.Println("心跳包", src)

		if src <= 7 {
			continue
		}

		data := buffer[:src] //处理MAC地址
		nowlen := 0
		addrlen := BatteryData["Address"]
		addrdata := data[nowlen:addrlen]
		fmt.Printf("MAC地址：% 2x\n", addrdata)

		nowlen = nowlen + addrlen           //处理组序号数据
		GNlen := BatteryData["GroupNumber"] //GN组序号  ,  GroupNumber组电流
		GNdata := data[nowlen:(GNlen + nowlen)]
		GNnumber := dataConversion(GNdata) //数据转换
		fmt.Println("组序号：", GNnumber)
		//fmt.Println(GNdata)

		nowlen = nowlen + GNlen              //处理组电压数据
		GVlen := BatteryData["GroupVoltage"] //GV组电压 , GroupVoltage组电压
		GVdata := data[nowlen:(GVlen + nowlen)]
		GVnumber1 := dataConversion1(GVdata)
		GVnumber2 := GVnumber1 / 10 //协议操作
		fmt.Println("组电压：", GVnumber2, "V")
		//fmt.Println(GVdata)

		nowlen = nowlen + GVlen              //处理组电流数据
		GClen := BatteryData["GroupCurrent"] //GC组电流，GroupCurrent组电流
		GCdata := data[nowlen:(GClen + nowlen)]
		GCnumber1 := dataConversion1(GCdata)
		GCnumber2 := GCnumber1 / 10
		fmt.Println("组电流: ", GCnumber2, "A")
		//fmt.Println(GCdata)

		nowlen = nowlen + GClen            //处理组状态数据
		GSlen := BatteryData["GroupState"] //GS组状态  ，GroupState组状态
		GSdata := data[nowlen:(GSlen + nowlen)]
		GSnumber := dataConversion1(GSdata)
		fmt.Printf("组状态： ")
		Binary1 := (GSnumber >> 0) & 1
		Binary2 := (GSnumber >> 1) & 1
		Binary3 := (GSnumber >> 2) & 1
		Binary4 := (GSnumber >> 3) & 1
		Binary5 := (GSnumber >> 4) & 1
		Binary6 := (GSnumber >> 5) & 1
		// Binary7 := (GSnumber >> 6) & 1
		// Binary8 := (GSnumber >> 7) & 1
		// Binary9 := (GSnumber >> 8) & 1
		// Binary10 := (GSnumber >> 9) & 1
		// Binary11 := (GSnumber >> 10) & 1
		// Binary12 := (GSnumber >> 11) & 1
		// Binary13 := (GSnumber >> 12) & 1
		// Binary14 := (GSnumber >> 13) & 1
		// Binary15 := (GSnumber >> 14) & 1
		// Binary16 := (GSnumber >> 15) & 1
		if Binary1 == 0 && Binary2 == 0 && Binary3 == 0 {
			fmt.Printf("电池组放电  ")
		}
		if Binary1 == 1 && Binary2 == 0 && Binary3 == 0 {
			fmt.Printf("电池组充电  ")
		}
		if Binary1 == 0 && Binary2 == 1 && Binary3 == 0 {
			fmt.Printf("电池组涓流充电  ")
		}
		if Binary3 == 0 && Binary4 == 0 {
			fmt.Printf("组电压正常  ")
		}
		if Binary3 == 1 && Binary4 == 0 {
			fmt.Printf("组电压过压  ")
		}
		if Binary3 == 0 && Binary4 == 1 {
			fmt.Printf("组电压欠压  ")
		}
		if Binary5 == 0 && Binary6 == 0 {
			fmt.Printf("组电流正常  ")
		}
		if Binary5 == 1 && Binary6 == 0 {
			fmt.Printf("组充电电流过大  ")
		}
		if Binary5 == 0 && Binary6 == 1 {
			fmt.Printf("组放电电流过大  ")
		}
		fmt.Printf("\n")
		//fmt.Println("组状态:", GSnumber)

		nowlen = nowlen + GSlen               //处理保留字节数据
		RBlen := BatteryData["ReservedBytes"] //RB保留字节 ，ReservedBytes保留字节
		RBdata := data[nowlen:(RBlen + nowlen)]
		fmt.Println("保留字节：", RBdata)

		nowlen = nowlen + RBlen               //处理电池组数量数据
		BNlen := BatteryData["BatteryNumber"] //BN 电池组数量 ， BatteryNumber电池组数量
		BNdata := data[nowlen:(BNlen + nowlen)]
		//var BNdata1 int
		BNdatanumber := dataConversion(BNdata)
		fmt.Println("组电池数量：", BNdatanumber)
		fmt.Printf("\n")
		fmt.Println("电池数据：\n")

		// var BinaryContent []byte
		nowlen = nowlen + BNlen               //循环处理电池数据                     // 20
		BCdata := data[nowlen:src]            //  BCdata=电池数据内容
		BClen := BatteryData["BinaryContent"] //每节电池长度
		nowlen = 0
		for i := 1; i <= BNdatanumber; i++ { //电池组数量  2
			BCdataOne := BCdata[nowlen:BClen]
			nowlen = BClen + nowlen
			BClen = BClen * (i + 1)
			fmt.Printf("第%d电池数据：", i)
			fmt.Printf("\n")
			//fmt.Println(BCdataOne)
			//fmt.Printf("\n")
			analysis(BCdataOne)
		}

	}
}

func analysis(b []byte) { //Analysis解析
	nowlen := 0
	MSNlen := BatteryData["MonitoringSerialNumber"] //MSN 检测模块序号MonitoringSerialNumber
	MSNdata := b[nowlen:MSNlen]
	nowlen = nowlen + MSNlen
	MSNnumber := dataConversion(MSNdata)
	fmt.Println("监测模块序号：", MSNnumber)
	//var BNdata1 int
	Modellen := BatteryData["Model"]
	Modeldata := b[nowlen:(nowlen + Modellen)]
	nowlen = nowlen + Modellen
	Modelnumber := dataConversion(Modeldata)
	fmt.Println("监测模块型号：", Modelnumber)

	BSlen := BatteryData["BatteryStatus"]
	BSdata := b[nowlen:(nowlen + BSlen)]
	nowlen = nowlen + BSlen
	BSnumber := dataConversion1(BSdata)
	fmt.Printf("电池状态：")
	BinaryT1 := (BSnumber >> 0) & 1
	BinaryT2 := (BSnumber >> 1) & 1
	BinaryT3 := (BSnumber >> 2) & 1
	BinaryT4 := (BSnumber >> 3) & 1
	BinaryT5 := (BSnumber >> 4) & 1
	BinaryT6 := (BSnumber >> 5) & 1
	BinaryT7 := (BSnumber >> 6) & 1
	// BinaryT8 := (BSnumber >> 7) & 1
	// BinaryT9 := (BSnumber >> 8) & 1
	// BinaryT10 := (BSnumber >> 9) & 1
	// BinaryT11 := (BSnumber >> 10) & 1
	// BinaryT12 := (BSnumber >> 11) & 1
	// BinaryT13 := (BSnumber >> 12) & 1
	// BinaryT14 := (BSnumber >> 13) & 1
	BinaryT15 := (BSnumber >> 14) & 1
	BinaryT16 := (BSnumber >> 15) & 1

	if BinaryT1 == 1 {
		fmt.Printf("内阻高  ")
	}
	if BinaryT2 == 1 {
		fmt.Printf("温度高  ")
	}
	if BinaryT3 == 1 {
		fmt.Printf("温度低  ")
	}
	if BinaryT4 == 1 {
		fmt.Printf("电压低  ")
	}
	if BinaryT5 == 1 {
		fmt.Printf("电压高  ")
	}
	if BinaryT6 == 1 {
		fmt.Printf("电池离线  ")
	}
	if BinaryT7 == 1 {
		fmt.Printf("温传坏  ")
	}
	if BinaryT15 == 1 {
		fmt.Printf("温度2是负温度  ")
	}
	if BinaryT16 == 1 {
		fmt.Printf("温度1是负温度  ")
	}
	fmt.Printf("\n")
	// fmt.Println("电池状态:", BSdata)

	BIRlen := BatteryData["BatteryInternalResistance"]
	BIRdata := b[nowlen:(nowlen + BIRlen)]
	nowlen = nowlen + BIRlen
	BIRnumber := dataConversions(BIRdata)
	ConversionValue := (float64(BIRnumber[0])*256 + float64(BIRnumber[1])) / 100
	fmt.Println("电池内阻：", ConversionValue, "Ω")

	BVlen := BatteryData["BatteryVoltage"]
	BVdata := b[nowlen:(nowlen + BVlen)]
	nowlen = nowlen + BVlen
	BVnumber := dataConversions(BVdata)
	ConversionValue0 := (float64(BVnumber[0])*256 + float64(BVnumber[1])) / 100
	fmt.Println("电池电压：", ConversionValue0, "V")

	BT1len := BatteryData["BatteryTemperature1"]
	BT1data := b[nowlen:(nowlen + BT1len)]
	nowlen = nowlen + BT1len
	BT1number := dataConversions(BT1data)
	ConversionValue1 := (float64(BT1number[0])*256 + float64(BT1number[1])) / 100
	fmt.Println("电池1温度：", ConversionValue1, "℃")

	BT2len := BatteryData["BatteryTemperature2"]
	BT2data := b[nowlen:(nowlen + BT2len)]
	nowlen = nowlen + BT2len
	BT2number := dataConversions(BT2data)
	ConversionValue2 := (float64(BT2number[0])*256 + float64(BT2number[1])) / 100
	fmt.Println("电池2温度：", ConversionValue2, "℃")

	BT3len := BatteryData["BatteryTemperature3"]
	BT3data := b[nowlen:(nowlen + BT3len)]
	nowlen = nowlen + BT3len
	BT3number := dataConversions(BT3data)
	ConversionValue3 := (float64(BT3number[0])*256 + float64(BT3number[1])) / 100
	fmt.Println("电池3温度：", ConversionValue3, "℃")

	BT4len := BatteryData["BatteryTemperature4"]
	BT4data := b[nowlen:(nowlen + BT4len)]
	nowlen = nowlen + BT4len
	BT4number := dataConversions(BT4data)
	ConversionValue4 := (float64(BT4number[0])*256 + float64(BT4number[1])) / 100
	fmt.Println("电池4温度：", ConversionValue4, "℃")

	Ripplelen := BatteryData["Ripple"]
	Rippledata := b[nowlen:(nowlen + Ripplelen)]
	nowlen = nowlen + Ripplelen
	Ripplenumber := dataConversions(Rippledata)
	ConversionValue5 := (float64(Ripplenumber[0])*256 + float64(Ripplenumber[1])) / 100
	fmt.Println("纹波：", ConversionValue5)

	RB2len := BatteryData["ReservedBytes2"]
	RB2data := b[nowlen:(nowlen + RB2len)]
	nowlen = nowlen + RB2len
	fmt.Println("保留字节：", RB2data)
	fmt.Printf("\n")

}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "严重错误: %s", err.Error())
		os.Exit(1)
	}
}
