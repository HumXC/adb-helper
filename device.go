package adb

import (
	"os"
)

type Device struct {
	IsOnline    bool
	ID          string
	USB         string
	Product     string
	Model       string
	Device      string
	TransportID int
	// 此处的 Cmd 已经带有指定设备的参数，例如 "-s 192.168.1.3"
	// 当调用 Cmd("ls") 时，实际运行的命令是 "adb -s 192.168.1.3 shell ls"
	Cmd ADBRunner
	// adb 执行文件的路径
	ADBPath string
	Input   Input
}

// 直接截图传输图片，截图过程中如果触碰屏幕，可能会导致失败
func (d *Device) Screencap() ([]byte, error) {
	return d.Cmd("shell screencap -p")
}

// 截图并保存文件到 fileName
func (d *Device) ScreencapTo(fileName string) (*os.File, error) {
	b, err := d.Screencap()
	if err != nil {
		return nil, err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	_, err = f.Write(b)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	return f, nil
}
