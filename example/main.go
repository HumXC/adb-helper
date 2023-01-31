package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/HumXC/adb-helper/helper"
)

func main() {
	// 运行之前确保 adb 已经连接到了设备
	helper := helper.New(RunADB)
	err := helper.Input.Power()
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1000 * time.Millisecond)
	b, err := helper.Screencap()
	if err != nil {
		fmt.Println(err)
	}
	f, _ := os.Create("./img.png")
	defer f.Close()
	f.Write(b)
}

func runCMD(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
func RunADB(args string) ([]byte, error) {
	cmd := strings.Split(args, " ")
	return runCMD("adb", cmd...)
}
