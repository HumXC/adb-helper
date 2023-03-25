package adb

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// 用于执行 ADB 命令，例如:
// run("shell ls")
type ADBRunner = func(args string) ([]byte, error)

type ADBHelp struct {
	ADBRunner ADBRunner
}

// adb 是运行 adb 时使用的命令，可以使用指定 adb 的路径，例如 "/usr/bin/adb"
func NewADBRunner(adb string) ADBRunner {
	return func(args string) ([]byte, error) {
		cmd := strings.Split(args, " ")
		c := exec.Command(adb, cmd...)
		var stdErr bytes.Buffer
		c.Stderr = &stdErr
		out, err := c.Output()
		if err != nil {
			errStr := stdErr.String()
			err = errors.New(errStr[:len(errStr)-1])
			return out, err
		}
		err = CheckError(out)
		return out, err
	}
}

func (a *ADBHelp) KillServer() error {
	_, err := a.ADBRunner("kill-server")
	return err
}

// 获取所有已经连接的设备
func (a *ADBHelp) Devices() ([]Device, error) {
	out, err := a.ADBRunner("devices -l")
	if err != nil {
		return nil, err
	}
	result := make([]Device, 0)
	buf := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buf)
	r := regexp.MustCompile(`\s+`)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" || strings.HasPrefix(text, "*") || strings.HasPrefix(text, "List of devices attached") {
			continue
		}
		text = r.ReplaceAllString(text, " ")
		args := strings.Split(text, " ")
		d := Device{
			ID: args[0],
		}
		for i := 1; i < len(args); i++ {
			arg := strings.Split(args[i], ":")
			switch arg[0] {
			case "usb":
				d.USB = arg[1]
			case "product":
				d.Product = arg[1]
			case "model":
				d.Model = arg[1]
			case "device":
				if i == 1 {
					d.IsOnline = true
					continue
				}
				d.Device = arg[1]
			case "transport_id":
				d.TransportID, _ = strconv.Atoi(arg[1])
			}
		}
		d.runner = a.ADBRunner
		d.preArg = "-s " + args[0] + " "
		result = append(result, d)
	}
	return result, nil
}

// 连接一个网络设备，如果 err==nil 那么连接成功
func (a *ADBHelp) Connect(host string) error {
	_, err := a.ADBRunner("connect " + host)
	return err
}

// 断开网络设备的连接，如果 host=="" 则会断开所有网络设备连接
func (a *ADBHelp) Disconnect(host string) error {
	_, err := a.ADBRunner("disconnect " + host)
	return err
}

// 使用机器自带的 adb 命令，需要安装 adb
func Default() ADBHelp {
	return New(NewADBRunner("adb"))
}

// / ADBRunner 见 NewADBRunner()，你也可以自己实现 ADBRunner
func New(adbRunner ADBRunner) ADBHelp {
	return ADBHelp{
		ADBRunner: adbRunner,
	}
}

// 命令运行时可能会出现异常并在控制台输出，此函数就是为了识别这些可能的异常
// 另一方面，有些时候命令执行是成功的，但是执行的结果我将其视为 "失败"
// 例如执行 adb connect "" 时 ，输出是 "empty address..."，我认为将其视为 error 是合适的
// 此函数并未涵盖所有的可能的异常输出，我只是添加了我遇到的
func CheckError(output []byte) (err error) {
	if len(output) < 2 {
		return
	}
	outStr := string(output)
	errFlag := []string{
		"Exception",
		"error:",
		"empty address",
		"failed",
	}
	for _, flag := range errFlag {
		if strings.HasPrefix(outStr, flag) {
			return errors.New(outStr[:len(outStr)-1])
		}
	}
	return nil
}
