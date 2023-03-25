package adb

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// 用于执行 ADB 命令，例如:
// run("shell ls")
type ADBRunner = func(args string) ([]byte, error)

type Server struct {
	Cmd ADBRunner
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
			msg := err.Error()
			if errStr := stdErr.String(); errStr != "" {
				msg += ": " + errStr[:len(errStr)-1]
			}
			return out, fmt.Errorf("exec error [%s]: %s", c.String(), msg)
		}
		if IsErrorOutput(out) {
			msg := string(out[:len(out)-1])
			return out, fmt.Errorf("exec error [%s]: %s", c.String(), msg)
		}
		return out, err
	}
}

func (s *Server) KillServer() error {
	_, err := s.Cmd("kill-server")
	return err
}

// 获取所有已经连接的设备
func (s *Server) Devices() ([]Device, error) {
	out, err := s.Cmd("devices -l")
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
		// cmd 增加设备前缀
		d.Cmd = func(_args string) ([]byte, error) {
			_args = "-s " + args[0] + " " + _args
			return s.Cmd(_args)
		}
		d.Input = &input{cmd: d.Cmd}
		result = append(result, d)
	}
	return result, nil
}

// 连接一个网络设备，如果 err==nil 那么连接成功
func (s *Server) Connect(host string) error {
	_, err := s.Cmd("connect " + host)
	return err
}

// 断开网络设备的连接，如果 host=="" 则会断开所有网络设备连接
func (s *Server) Disconnect(host string) error {
	_, err := s.Cmd("disconnect " + host)
	return err
}

// 使用机器自带的 adb 命令，需要安装 adb
func DefaultServer() Server {
	return NewServer(NewADBRunner("adb"))
}

// cmd 见 NewADBRunner()，你也可以自己实现 ADBRunner
func NewServer(cmd ADBRunner) Server {
	return Server{
		Cmd: cmd,
	}
}

// 命令运行时可能会出现异常并在控制台输出，此函数就是为了识别这些可能的异常
// 另一方面，有些时候命令执行是成功的，但是执行的结果我将其视为 "失败"
// 例如执行 adb connect "" 时 ，输出是 "empty address..."，我认为将其视为 error 是合适的
// 此函数并未涵盖所有的可能的异常输出，我只是添加了我遇到的
func IsErrorOutput(output []byte) bool {
	if len(output) < 2 {
		return false
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
			return true
		}
	}
	return false
}
