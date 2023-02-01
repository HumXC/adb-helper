package main

import (
	"fmt"
	"time"

	"github.com/HumXC/adb-helper/helper"
)

func main() {
	// 运行之前确保 adb 已经连接到了设备
	helper := helper.Default()
	devices, err := helper.Devices()
	if err != nil {
		fmt.Println(err)
	}
	if len(devices) == 0 {
		return
	}
	devices[0].Power()
	time.Sleep(time.Second)
	devices[0].ScreencapTo("./sc.png")
}
