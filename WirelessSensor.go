package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tarm/goserial"
)

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}

func timer1() {
	timer1 := time.NewTicker(3 * time.Second)
	select {
	case <-timer1.C:
	}
}

func BytesCombine(pBytes ...[]byte) []byte { //字节拼接方法
	return bytes.Join(pBytes, []byte(""))
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

func readAll(b io.ReadWriteCloser, len int) ([]byte, error) { //读取内容

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
		// if one_len == 0 {

		// 	return []byte(""), fmt.Errorf("没读取到数据")
		// }
		readedLen = readedLen + one_len

		if readedLen == len {
			break
		}

	}
	return pack, nil
}

func readDOM(b io.ReadWriteCloser, s0 []byte, datalen int) ([]byte, error) {
	var data0 []byte
	var err error
	w, err := b.Write(s0)
	if err != nil && w <= 0 {
		fmt.Println("内部录入地址命令错误")
		return data0, err
	}

	buffer0 := make([]byte, 512)
	buffer0, err = readAll(b, datalen)

	if err != nil || len(buffer0) == 0 {
		return data0, err
	}
	data0 = buffer0[0:]
	return data0, err
}

func SerialPortSelect(comname string) (io.ReadWriteCloser, []byte, int) {

	var baudnum int
	baudnum = 9600
	var s100 []byte
	var address []byte
	a := &serial.Config{Name: comname, Baud: baudnum, ReadTimeout: 15 * time.Second}
	b, err := serial.OpenPort(a)
	if err != nil && b == nil {
		// fmt.Printf("未读取到该串口信息\n\n")
		// fmt.Printf("解决方案：1.")
		return b, address, -1
	}
	sendadderhead := []byte{0x20, 0x07, 0xff, 0xff, 0xff, 0xff}
	s100 = BytesCombine(sendadderhead, []byte{0x00})
	address, err = readDOM(b, s100, 6)

	if len(address) == 0 || err != nil {
		b.Close()
		//fmt.Printf("没有读取到数据，请检查电源线路\n")
		return b, address, -2
	}
	return b, address, 1

}

func sss(comname string) (io.ReadWriteCloser, []byte, int) {
	flag := 3
	s9 := true
	s := 1
	var b io.ReadWriteCloser
	var err int
	var address []byte
	for i := 0; i < flag; i++ {
		if i == 0 {
			fmt.Println("正在连接......")
		} else if i == 1 {
			fmt.Println("正在重新连接......")
		} else if i == 2 {
			fmt.Println("正在尝试最后一次重新链接......")
		}
		b, address, err = SerialPortSelect(comname)

		//fmt.Println(err,address)

		if err < 0 {
			time.Sleep(20 * time.Second)
			s9 = false
			continue
		} else {
			s9 = true
			break
		}
	}
	if !s9 {
		fmt.Printf("连接失败\n")
		fmt.Printf("请检查如下问题：1.端口编号选择是否正确   2.串口及电源线路是否链接正常   3.串口是否被占用   4.串口是否属于配置模式\n\n")
		s = -1

	}
	return b, address, s
}

func Configure1(s44 []byte, s55 []byte, b io.ReadWriteCloser) {
MEIMEI:
	var p int
	var s33 []byte
	fmt.Printf("1，永久修改  2，临时修改（选择编号回车键确定）")
	s33 = []byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff}
	fmt.Scanf("%d\n", &p)
	if p == 1 {
		var s2 []byte
		s2 = BytesCombine(s33, s44)
		w, err := b.Write(s2)
		if err != nil && w == 0 {
			fmt.Printf("数据修改失败")
		}
	} else if p == 2 {
		var s2 []byte
		s2 = BytesCombine(s33, s55)
		w, err := b.Write(s2)
		if err != nil && w == 0 {
			fmt.Printf("数据修改失败")
		}
	} else {
		fmt.Printf("输入错误\n")
		goto MEIMEI
	}

}

func Configure(comname string, b io.ReadWriteCloser) {
	fmt.Printf("1,请先长按config键进入配置模式（红绿灯同时闪烁-->进入成功)\n")
	fmt.Printf("注意:(进入配置模式，如果在一分钟之内没有任何配置操作，系统将会自动退出配置模式)\n\n")
LOOP:
	var k int
	fmt.Printf("请确定成功进入配置模式(红绿灯同时闪烁)，输入1，回车确定，请输入：")
	fmt.Scanf("%d\n", &k)
	if k == 1 {
		fmt.Printf("waiting……\n")
		var baudnum int
		baudnum = 38400
		a := &serial.Config{Name: comname, Baud: baudnum, ReadTimeout: 10 * time.Second}
		b, err := serial.OpenPort(a)
		if err != nil || b == nil {
			fmt.Printf("配置时，获取串口信息失败")
		}

		for i := 0; i < 3; {
			i++
			time.Sleep(5 * time.Second)
			s := []byte{0x20, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff}
			w, err := b.Write(s)
			if err != nil || w == 0 {
				fmt.Printf("配置时，数据写入错误")
			}
			fmt.Println(w)
		}
	HAHA:
		fmt.Printf("设置菜单：\n")
		fmt.Printf("1.节点地址  2.节点类型  3.网络类型  4.网络ID  5.无线频点  6.数据编码 \n ")
		fmt.Printf("7.发送模式  8.波特率  9.校验  10.数据源址  11.发射功率  12.主动上报\n\n")
		fmt.Printf("注意：不提供返回主菜单机制，因为修改完成后，仍然处于配置状态，所以需要断开电源重新进入\n\n")
		var s1 int
		//var s33 []byte
		fmt.Printf("请选择修改项的编号（回车键确定）： ")
		fmt.Scanf("%d\n", &s1)

		//修改节点地址

		if s1 == 1 {
			var s2 string
			fmt.Printf("节点地址（取值范围：00000001-FFFFFFFE）\n")
			fmt.Printf("请修改: ")
			fmt.Scanf("%s\n", &s2)
		HEIHEI:
			a, _ := strconv.ParseInt(s2, 16, 64) //16进制字符串转10进制数字
			b_buf := bytes.NewBuffer([]byte{})
			binary.Write(b_buf, binary.BigEndian, int32(a))
			row := b_buf.Bytes()
			var p int
			fmt.Printf("1，永久修改  2，临时修改（选择编号回车键确定）")
			fmt.Scanf("%d\n", &p)
			if p == 1 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0C, 0xff, 0xff, 0xff, 0xff, 0x01, 0x02}, row)
				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else if p == 2 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0C, 0xff, 0xff, 0xff, 0xff, 0x00, 0x02}, row)
				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else {
				fmt.Printf("输入错误！！！\n")
				goto HEIHEI
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA

		}

		//修改节点类型
		if s1 == 2 {
			fmt.Printf("节点类型：1.中心节点  2.中继路由  3.终端节点\n")
		LIULIU:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x07, 0x01}
				s55 := []byte{0x00, 0x07, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x07, 0x02}
				s55 := []byte{0x00, 0x07, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x07, 0x03}
				s55 := []byte{0x00, 0x07, 0x03}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto LIULIU
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}

		//修改网络类型
		if s1 == 3 {
			fmt.Printf("网络类型： 1.网状网  2.星型网  3.链型网  4.对等网")
		SENSEN:
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			var p1 int

			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x06, 0x01}
				s55 := []byte{0x00, 0x06, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x06, 0x02}
				s55 := []byte{0x00, 0x06, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x06, 0x03}
				s55 := []byte{0x00, 0x06, 0x03}
				Configure1(s44, s55, b)
			} else if p1 == 4 {
				s44 := []byte{0x01, 0x06, 0x04}
				s55 := []byte{0x00, 0x06, 0x04}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto SENSEN
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改网络ID
		if s1 == 4 {
			var s2 string
			fmt.Printf("网络ID： （修改范围：00-FF）\n")
			fmt.Printf("请修改：")
			fmt.Scanf("%s\n", &s2)
		LANLAN:
			a, _ := strconv.ParseInt(s2, 16, 64) //16进制字符串转10进制数字
			b_buf := bytes.NewBuffer([]byte{})
			binary.Write(b_buf, binary.BigEndian, int32(a))
			row := b_buf.Bytes()
			var rowData = map[string]int{
				"nouse": 2,
				"use":   2,
			}
			data := row[0:]
			nowlen := 0
			nouselen := rowData["nouse"]

			nowlen = nowlen + nouselen
			uselen := rowData["use"]
			usedata := data[nowlen:(uselen + nowlen)]
			fmt.Println(usedata)
			var p int
			fmt.Printf("1，永久修改  2，临时修改（选择编号回车键确定）")
			fmt.Scanf("%d\n", &p)
			if p == 1 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x01, 0x05}, usedata)

				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else if p == 2 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x00, 0x05}, usedata)
				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else {
				fmt.Printf("输入错误！！！\n")
				goto LANLAN
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改无线频点
		if s1 == 5 {
			fmt.Printf("无线频点   1.手动修改   2.自动\n")
			fmt.Printf("请选择修改编号（按回车键确定）：")
			var s22 int
			fmt.Scanf("%d\n", &s22)
			if s22 == 1 {

				fmt.Printf("请修改（修改范围0-F）：")
				var s2 string
				fmt.Scanf("%s\n", &s2)
			YANGYANG:
				a, _ := strconv.ParseInt(s2, 16, 64) //16进制字符串转10进制数字
				fmt.Println(a)

				b_buf := bytes.NewBuffer([]byte{})

				binary.Write(b_buf, binary.BigEndian, int32(a))
				row := b_buf.Bytes()
				var rowData = map[string]int{
					"nouse": 3,
					"use":   1,
				}
				data := row[0:]
				nowlen := 0
				nouselen := rowData["nouse"]

				nowlen = nowlen + nouselen
				uselen := rowData["use"]
				usedata := data[nowlen:(uselen + nowlen)]
				fmt.Println(usedata)

				var p int
				//var s33 []byte
				fmt.Printf("1，永久修改  2，临时修改（选择编号回车键确定）")
				//s33 = []byte{0x22,0x0C,0xff,0xff,0xff,0xff}
				fmt.Scanf("%d\n", &p)
				if p == 1 {
					var s2 []byte
					s2 = BytesCombine([]byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff, 0x01, 0x03}, usedata)
					fmt.Println(s2)
					fmt.Println(row)
					w, err := b.Write(s2)
					if err != nil && w == 0 {
						fmt.Printf("数据修改失败")
					}
				} else if p == 2 {
					var s2 []byte
					s2 = BytesCombine([]byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff, 0x00, 0x03}, usedata)
					w, err := b.Write(s2)
					if err != nil && w == 0 {
						fmt.Printf("数据修改失败")
					}
				} else {
					fmt.Printf("输入错误！！！\n")
					goto YANGYANG
				}
				fmt.Printf("\n")
				fmt.Printf("数据修改完成！！！\n")
				goto HAHA
			} else if s22 == 2 {
				fmt.Printf("协议没有说明无线频点“自动”命令！！！")
				goto HAHA

			}
		}
		//数据编码配置
		if s1 == 6 {
			fmt.Printf("数据编码： 1.ASCII  2.HEX\n")
		MANMAN:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x14, 0x01}
				s55 := []byte{0x00, 0x14, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x14, 0x02}
				s55 := []byte{0x00, 0x14, 0x02}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto MANMAN
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改发送模式
		if s1 == 7 {
			fmt.Printf("发送模式：1.广播  2.固定目标  3.SHUNCOM协议  4.MODBUS_RTU\n")
		XINXIN:
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			var p1 int
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x08, 0x01}
				s55 := []byte{0x00, 0x08, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				fmt.Printf("\n")
				fmt.Printf("输入目标地址：")

				s44 := []byte{0x01, 0x08, 0x02}
				s55 := []byte{0x00, 0x08, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x08, 0x03}
				s55 := []byte{0x00, 0x08, 0x03}
				Configure1(s44, s55, b)
			} else if p1 == 4 {
				s44 := []byte{0x01, 0x08, 0x04}
				s55 := []byte{0x00, 0x08, 0x04}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto XINXIN
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改波特率
		if s1 == 8 {
			fmt.Printf("波特率：1.1200  2.2400  3.4800  4.9600  5.19200  6.38400  7.57600  8.115200")
		GUGU:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x10, 0x01}
				s55 := []byte{0x00, 0x10, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x10, 0x02}
				s55 := []byte{0x00, 0x10, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x10, 0x03}
				s55 := []byte{0x00, 0x10, 0x03}
				Configure1(s44, s55, b)
			} else if p1 == 4 {
				s44 := []byte{0x01, 0x10, 0x04}
				s55 := []byte{0x00, 0x10, 0x04}
				Configure1(s44, s55, b)
			} else if p1 == 5 {
				s44 := []byte{0x01, 0x10, 0x05}
				s55 := []byte{0x00, 0x10, 0x05}
				Configure1(s44, s55, b)
			} else if p1 == 6 {
				s44 := []byte{0x01, 0x10, 0x06}
				s55 := []byte{0x00, 0x10, 0x06}
				Configure1(s44, s55, b)
			} else if p1 == 7 {
				s44 := []byte{0x01, 0x10, 0x07}
				s55 := []byte{0x00, 0x10, 0x07}
				Configure1(s44, s55, b)
			} else if p1 == 8 {
				s44 := []byte{0x01, 0x10, 0x08}
				s55 := []byte{0x00, 0x10, 0x08}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto GUGU
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改校验
		if s1 == 9 {
			fmt.Printf("校验： 1.None  2.Even  3.0dd  4.Mark  5.Space")
		ZHUZHU:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x11, 0x01}
				s55 := []byte{0x00, 0x11, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x11, 0x02}
				s55 := []byte{0x00, 0x11, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x11, 0x03}
				s55 := []byte{0x00, 0x11, 0x03}
				Configure1(s44, s55, b)
			} else if p1 == 4 {
				s44 := []byte{0x01, 0x11, 0x04}
				s55 := []byte{0x00, 0x11, 0x04}
				Configure1(s44, s55, b)
			} else if p1 == 5 {
				s44 := []byte{0x01, 0x11, 0x05}
				s55 := []byte{0x00, 0x11, 0x05}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto ZHUZHU
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！")
			goto HAHA
		}
		//修改数据源址
		if s1 == 10 {
			fmt.Printf("数据源址： 1.不输出  2.输出")
		NIUNIU:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x15, 0x01}
				s55 := []byte{0x00, 0x15, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x15, 0x02}
				s55 := []byte{0x00, 0x15, 0x02}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto NIUNIU
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//修改发射功率
		if s1 == 11 {
			fmt.Printf("发射功率： 1.最小  2.中等  3.最大")
		XINGXING:
			var p1 int
			fmt.Printf("请输入修改编号（回车键确定修改）： ")
			fmt.Scanf("%d\n", &p1)
			if p1 == 1 {
				s44 := []byte{0x01, 0x04, 0x01}
				s55 := []byte{0x00, 0x04, 0x01}
				Configure1(s44, s55, b)
			} else if p1 == 2 {
				s44 := []byte{0x01, 0x04, 0x02}
				s55 := []byte{0x00, 0x04, 0x02}
				Configure1(s44, s55, b)
			} else if p1 == 3 {
				s44 := []byte{0x01, 0x04, 0x03}
				s55 := []byte{0x00, 0x04, 0x03}
				Configure1(s44, s55, b)
			} else {
				fmt.Printf("输入错误！！！\n")
				goto XINGXING
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}
		//主动上报
		if s1 == 12 {
			var s2 string
			fmt.Printf("主动上报  （修改范围：0-65535）\n")
			fmt.Printf("请修改（修改范围0000-ffff）：")
			fmt.Scanf("%s\n", s2)
		BOBO:
			a, _ := strconv.ParseInt(s2, 16, 64) //16进制字符串转10进制数字
			fmt.Println(a)
			b_buf := bytes.NewBuffer([]byte{})
			binary.Write(b_buf, binary.BigEndian, int32(a))
			row := b_buf.Bytes()
			var p int
			//var s33 []byte
			fmt.Printf("1，永久修改  2，临时修改（选择编号回车键确定）")
			//s33 = []byte{0x22,0x0C,0xff,0xff,0xff,0xff}
			fmt.Scanf("%d\n", &p)
			if p == 1 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x01, 0x20}, row)
				fmt.Println(s2)
				fmt.Println(row)

				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else if p == 2 {
				var s2 []byte
				s2 = BytesCombine([]byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x00, 0x20}, row)
				w, err := b.Write(s2)
				if err != nil && w == 0 {
					fmt.Printf("数据修改失败")
				}
			} else {
				fmt.Printf("输入错误！！！\n")
				goto BOBO
			}
			fmt.Printf("\n")
			fmt.Printf("数据修改完成！！！\n")
			goto HAHA
		}

	} else {
		fmt.Printf("输入错误！！！\n")
		goto LOOP
	}
}

func Initi(comname string, b io.ReadWriteCloser, data1len int) {
	fmt.Printf("请长按config进入初始化模式（请在5秒内进入）\n\n")

	fmt.Printf("初始化将自动完成，请耐心等待\n")
	fmt.Printf("Loading……\n\n")
	var baudnum int
	baudnum = 38400
	a := &serial.Config{Name: comname, Baud: baudnum, ReadTimeout: 100 * time.Second}
	b, err := serial.OpenPort(a)
	if err != nil || b == nil {
		fmt.Printf("获取信息失败")
	}
	for i := 0; i < 3; {
		i++
		timer1()
		s := []byte{0x20, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff}
		w, err := b.Write(s)
		if err != nil && w == 0 {
			fmt.Printf("数据写入错误")
		}
	}

	// buffer1 := make([]byte, 256)
	// buffer1, err = readAll(b,data1len)
	// fmt.Println(data1len)
	// fmt.Println(buffer1)

	// if err != nil{
	// fmt.Printf("读取数据失败")
	// }
	// fmt.Printf("haha")
	// fmt.Println(b)
	timer1()

	s1 := []byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff, 0x01, 0x07, 0x03}
	w2, err := b.Write(s1)
	if err != nil && w2 == 0 {
		fmt.Printf("节点类型修改失败\n")
	}
	timer1()

	s2 := []byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff, 0x01, 0x08, 0x03}
	w3, err := b.Write(s2)
	if err != nil && w3 == 0 {
		fmt.Printf("发送模式修改失败\n")
	}
	timer1()
	s3 := []byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x01, 0x20, 0x00, 0x00}
	w4, err := b.Write(s3)
	if err != nil && w4 == 0 {
		fmt.Printf("主动上报修改失败\n")
	}

	timer1()
	s4 := []byte{0x22, 0x0c, 0xff, 0xff, 0xff, 0xff, 0x01, 0x02, 0x00, 0x00, 0x00, 0x01}
	w5, err := b.Write(s4)
	if err != nil && w5 == 0 {
		fmt.Printf("节点地址修改失败\n")
	}

	timer1()
	s5 := []byte{0x22, 0x09, 0xff, 0xff, 0xff, 0xff, 0x01, 0x14, 0x02}
	w6, err := b.Write(s5)
	if err != nil && w6 == 0 {
		fmt.Printf("数据编码修改失败\n")
	}

	timer1()
	s6 := []byte{0x22, 0x0a, 0xff, 0xff, 0xff, 0xff, 0x01, 0x05, 0x000, 0x000}
	w7, err := b.Write(s6)
	if err != nil && w7 == 0 {
		fmt.Printf("网络ID修改失败\n")
	}

	fmt.Println("初始化已完成、若要进入查询及配置模式，请断开电源重新进入\n\n")
}

func Exit() {

	os.Exit(1)

}

func GetSerialPortdata(b io.ReadWriteCloser, address []byte) {
	fmt.Println("正在努力加载 --> ┗|｀O′|┛ 嗷~~")
	fmt.Printf("\n")
	sendadderhead := []byte{0x20, 0x07, 0xff, 0xff, 0xff, 0xff}
	//节点地址
	s0 := BytesCombine(sendadderhead, []byte{0x00})
	var SensorData0 = map[string]int{
		"head_0":    1,
		"lenth_0":   1,
		"address_0": 4,
	}
	// data0, err := readDOM(b, s0, 6)

	data0 := address
	if len(data0) == 0 {
		return
	}
	nowlen0 := 0
	headlen0 := SensorData0["head_0"]
	nowlen0 = nowlen0 + headlen0
	lenthlen0 := SensorData0["lenth_0"]

	nowlen0 = nowlen0 + lenthlen0
	addresslen0 := SensorData0["address_0"]
	adressdata0 := data0[nowlen0:(addresslen0 + nowlen0)]
	fmt.Printf("节点地址    :    %02x  \n", adressdata0)

	// sendadderhead = BytesCombine(sendadderhead,adressdata0)
	// fmt.Printf("ceshi:")
	// fmt.Println(sendadderhead)

	sendadderhead = []byte{0x20, 0x07}
	sendadderhead = BytesCombine(sendadderhead, adressdata0)

	//节点类型
	s0 = BytesCombine(sendadderhead, []byte{0x07})
	var SensorData1 = map[string]int{
		"head_1":     1,
		"lenth_1":    1,
		"address_1":  4,
		"function_1": 1,
		"type_1":     1,
	}

	data1, err := readDOM(b, s0, 8)
	// fmt.Printf("b的值：")
	// fmt.Println(b)
	// fmt.Printf("***********************************")
	// fmt.Printf("data1的值")
	// fmt.Println(data1)

	if err != nil {
		return
	}
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
	// fmt.Println(typelen1)
	// fmt.Println(typedata1)
	number1 := dataConversion(typedata1)

	res1 := []string{"", "中心节点", "中继路由", "终端节点"}
	if number1 < 1 || number1 > 3 {
		fmt.Printf("节点类型    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("节点类型    ：   ")
		fmt.Println(res1[number1])
	}

	//网络类型					(注：漏掉了，再加上，所以编号为15)
	var SensorData15 = map[string]int{
		"head_15":     1,
		"lenth_15":    1,
		"address_15":  4,
		"function_15": 1,
		"type_15":     1,
	}
	s0 = BytesCombine(sendadderhead, []byte{0x06})

	data15, err := readDOM(b, s0, 8)

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
	res2 := []string{"", "网状网", "星型网", "链型网", "对等网"}
	if number15 < 1 || number15 > 4 {
		fmt.Printf("网络类型    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("网络类型    ：   ")
		fmt.Println(res2[number15])
	}

	//网络ID
	s0 = BytesCombine(sendadderhead, []byte{0x05})
	var SensorData2 = map[string]int{
		"head_2":     1,
		"lenth_2":    1,
		"address_2":  4,
		"function_2": 1,
		"NETID_2":    2,
	}

	data2, err := readDOM(b, s0, 9)

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

	fmt.Printf("  网络ID    :    ")
	for i := 0; i < 2; i++ {
		fmt.Printf("%02x ", NETIDdata2[i]) //
	}
	fmt.Println()

	//无线频点
	s0 = BytesCombine(sendadderhead, []byte{0x03})
	var SensorData3 = map[string]int{
		"head_3":           1,
		"lenth_3":          1,
		"address_3":        4,
		"function_3":       1,
		"frequencyPoint_3": 1,
	}

	data3, err := readDOM(b, s0, 8)

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
	fmt.Printf("无线频点    :    ")
	fmt.Printf("%02x ", FQdata3)
	fmt.Println()

	//数据编码
	s0 = BytesCombine(sendadderhead, []byte{0x14})
	var SensorData4 = map[string]int{
		"head_4":         1,
		"lenth_4":        1,
		"address_4":      4,
		"function_4":     1,
		"DataEncoding_4": 1,
	}
	data4, err := readDOM(b, s0, 8)

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
	res3 := []string{"", "ASCII", "HEX"}
	if number4 < 1 || number4 > 2 {
		fmt.Printf("数据编码    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("数据编码    ：   ")
		fmt.Println(res3[number4])
	}

	//发送模式
	s0 = BytesCombine(sendadderhead, []byte{0x08})
	var SensorData5 = map[string]int{
		"head_5":     1,
		"lenth_5":    1,
		"address_5":  4,
		"function_5": 1,
		"SendMod_5":  1,
	}

	data5, err := readDOM(b, s0, 8)

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
	res4 := []string{"", "广播", "固定协议", "SHUNCOM协议", "MODBUS_RTU"}
	if number5 < 1 || number5 > 5 {
		fmt.Printf("发送模式    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("发送模式    ：   ")
		fmt.Println(res4[number5])
	}

	//波特率
	s0 = BytesCombine(sendadderhead, []byte{0x10})
	var SensorData6 = map[string]int{
		"head_6":     1,
		"lenth_6":    1,
		"address_6":  4,
		"function_6": 1,
		"BaudRate_6": 1,
	}

	data6, err := readDOM(b, s0, 8)

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
	res5 := []string{"", "1200", "2400", "4800", "9600", "19200", "38400", "57600", "115200"}
	if number6 < 1 || number6 > 8 {
		fmt.Printf("  波特率    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("  波特率    ：   ")
		fmt.Println(res5[number6])
	}

	//校验
	s0 = BytesCombine(sendadderhead, []byte{0x11})
	var SensorData7 = map[string]int{
		"head_7":     1,
		"lenth_7":    1,
		"address_7":  4,
		"function_7": 1,
		"Check_7":    1,
	}

	data7, err := readDOM(b, s0, 8)

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
	res6 := map[int]string{
		1: "None",
		2: "Even",
		3: "0dd",
		4: "Mark",
		5: "Space",
	}
	if number7 < 1 || number7 > 5 {
		fmt.Printf("   校验    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {

		fmt.Printf("    校验    ：   ")
		fmt.Println(res6[number7])
	}

	//数据源址
	s0 = BytesCombine(sendadderhead, []byte{0x15})
	var SensorData9 = map[string]int{
		"head_9":        1,
		"lenth_9":       1,
		"address_9":     4,
		"function_9":    1,
		"DataAddress_9": 1,
	}

	data9, err := readDOM(b, s0, 8)

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
	res7 := []string{"", "不输出", "输出"}
	if number9 < 1 || number9 > 2 {
		fmt.Printf("数据源址    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("数据源址    ：   ")
		fmt.Println(res7[number9])
	}

	//发射功率
	s0 = BytesCombine(sendadderhead, []byte{0x04})
	var SensorData10 = map[string]int{
		"head_10":         1,
		"lenth_10":        1,
		"address_10":      4,
		"function_10":     1,
		"TansmitPower_10": 1,
	}

	data10, err := readDOM(b, s0, 8)

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
	res := []string{"", "最小", "中等", "最大"}
	if number10 < 1 || number10 > 3 {
		fmt.Printf("发射功率    ：   ")
		fmt.Println("出现溢出，接收到的数据与协议不相符")
	} else {
		fmt.Printf("发射功率    ：   ")
		fmt.Println(res[number10])
	}

	//主动上报
	s0 = BytesCombine(sendadderhead, []byte{0x20})
	var SensorData14 = map[string]int{
		"head_14":            1,
		"lenth_14":           1,
		"address_14":         4,
		"function_14":        1,
		"ActiveReporting_14": 2,
	}
	data14, err := readDOM(b, s0, 9)

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
	// b.Close()
	fmt.Printf("主动上报    ：   ")
	fmt.Println(number14)
	fmt.Printf("\n")
	fmt.Printf("加载完成！     ヾ(￣▽￣)Bye~Bye~\n")

}

func main() {
LOOP:
	var p1 int
	fmt.Printf("ComName:  1.O5  2.O1  3.O3  4.O2  5.COM4  6.COM5  7.COM6  8.COM7\n\n")
	fmt.Printf("请选择端口编号(回车键确定) :")
	fmt.Scanf("%d\n", &p1)
	res1 := []string{"", "/dev/ttyO5", "/dev/ttyO1", "/dev/ttyO3", "/dev/ttyO2", "COM4", "COM5", "COM6", "COM7"}
	comname := res1[p1]
	fmt.Printf("\n")

	var b io.ReadWriteCloser
	var address []byte
	var s int
	b, address, s = sss(comname)
	if s < 0 {
		goto LOOP
	}

	for {
		fmt.Printf("请选择进入方式：1，查询  2，配置  3，初始化  0，退出    【编号+Enter键】= 确认\n\n")
		fmt.Printf("请输入：")
		var types int
		fmt.Scanf("%d\n", &types)
	Looq:
		if types == 0 {
			Exit()
		} else if types == 1 {

			GetSerialPortdata(b, address)
			fmt.Println("0退出/1返回（回车键确定）")
			var p3 int
			fmt.Printf("请输入：")
			fmt.Scanf("%d\n", &p3)
			if p3 == 0 {
				os.Exit(1)
			} else if p3 == 1 {
				continue
			} else {
				goto Looq
			}
		} else if types == 2 {
			b.Close()

			Configure(comname, b)
		} else if types == 3 {
			b.Close()
			Initi(comname, b, 37)
			continue

		} else {
			continue
		}
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	fmt.Println(<-ch)
}
