package main

import (
	"fmt"
	"time"

	"github.com/HumXC/adb-helper"
)

func main() {
	// 运行之前确保 adb 已经连接到了设备
	server := adb.DefaultServer()
	devices, err := server.Devices()
	if err != nil {
		fmt.Println(err)
	}
	if len(devices) == 0 {
		return
	}
	err = devices[0].Input.Power()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second)
	devices[0].ScreencapTo("./sc.png")
}
