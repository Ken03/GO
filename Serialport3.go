package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/tarm/goserial"
	//"strings"
)

func readAll(b io.ReadWriteCloser, len int) ([]byte, error) {
	//如果解析出来的len长度明显不符合要求，则返回
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

var SensorData = map[string]int{
	"Address":       1,
	"gongnengma":    1,
	"shujushuliang": 1,
	"wendu":         2,
	"shidu":         2,
	"jiaoyanma":     2,
}

func test() {
	a := &serial.Config{Name: "COM4", Baud: 9600} //设置端口，波特率
	b, err := serial.OpenPort(a)
	if err != nil {
		fmt.Println("端口打开失败", err)
		return
	}
	defer b.Close()
	var s []byte //向传感器写入正确命令格式
	s = append(s, 0X01)
	s = append(s, 0X04)
	s = append(s, 0x00)
	s = append(s, 0x00)
	s = append(s, 0x00)
	s = append(s, 0x02)
	s = append(s, 0x71)
	s = append(s, 0xCB)

	w, err := b.Write(s)
	if err != nil && w <= 0 {
		fmt.Println("写入数据错误")
	}

	buffer := make([]byte, 256) //读取传感器返回数据
	buffer, err = readAll(b, 9)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := buffer[0:]               //接收的所有数据赋给data
	nowlen := 0                      //当前长度为0
	addrlen := SensorData["Address"] // key = "Address" , value = "1",截取Address的长度
	addrdata := data[nowlen:addrlen] //截取当前长度  到  addrlen  之间的数据
	addnumber := dataConversion(addrdata)
	fmt.Println("地址：", addnumber)

	nowlen = nowlen + addrlen
	gnmlen := SensorData["gongnengma"]
	// gnmdata := data[nowlen:(gnmlen + nowlen)]
	// fmt.Println("功能码：", gnmdata)

	nowlen = nowlen + gnmlen
	sjsllen := SensorData["shujushuliang"]
	// sjsldata := data[nowlen:(sjsllen + nowlen)]
	// fmt.Println("数据字节：", sjsldata)

	nowlen = nowlen + sjsllen
	wblen := SensorData["wendu"]
	wbdata := data[nowlen:(wblen + nowlen)]
	wbNumber := dataConversion1(wbdata)
	wbNumber2 := float64(wbNumber) / 10
	fmt.Println("温度：", wbNumber2, "℃")

	nowlen = nowlen + wblen
	sdlen := SensorData["shidu"]
	sddata := data[nowlen:(sdlen + nowlen)]
	sdNumber := dataConversion1(sddata)
	sdNumber2 := float64(sdNumber) / 10
	fmt.Println("湿度：", sdNumber2, "%rh")

	// nowlen = nowlen + sdlen
	// jymlen := SensorData["jiaoyanma"]
	// jymdata := data[nowlen:(jymlen + nowlen)]
	// fmt.Println("校验码：", jymdata)
	poststrdata := "dev_test,tag_name=TH_01_tem value=" + strconv.FormatFloat(wbNumber2, 'E', -1, 64)
	httpPostForm(poststrdata)
}

func timer1() {
	ticker := time.NewTicker(1 * time.Second)
	timer1 := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		fmt.Println(time.Now())

		select {
		case <-timer1.C:
			test()
			fmt.Printf("\n")
		}
	}

}

func httpPostForm(postdata string) {
	url := "http://192.168.1.13:8186/write"

	// A := bytes.NewReader([]byte("dev_test,tag_name=TH_01_tem value=27"))
	A := bytes.NewReader([]byte(postdata))
	request, err := http.NewRequest("POST", url, A)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	request.Header.Set("Content-Type", "text/plain;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(respBytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
func main() {

	timer1()
}
