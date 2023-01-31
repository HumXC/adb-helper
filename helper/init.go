package helper

import (
	"errors"
	"strings"
)

// 异常退出，运行命令时被人为打断可能会引起此错误
// 例如在运行 Input.Press(100,100,3000) 时，会长按屏幕 3 秒，如果途中用手指触摸屏幕，则操作被打断
var ErrExit = errors.New("exit")

// 用于执行 ADB 命令，例如:
// run("shell ls")
type ADBRunner = func(args string) ([]byte, error)
type ADBHelp struct {
	adbRunner ADBRunner
	Input     InputMan
}

// 直接截图传输图片，截图过程中如果触碰屏幕，可能会导致失败
func (a *ADBHelp) Screencap() ([]byte, error) {
	return a.adbRunner("shell screencap -p")
}

func New(adbRunner ADBRunner) ADBHelp {
	return ADBHelp{
		adbRunner: adbRunner,
		Input:     InputMan{adbRunner: adbRunner},
	}
}

// 命令运行时可能会出现异常并在控制台输出，此函数就是为了识别这些可能的异常
func ChexkError(output []byte) (err error) {
	if len(output) < 4 && len(output) > 128 {
		return
	}
	outStr := string(output)
	if strings.HasPrefix(outStr, "Exception") {
		return ErrExit
	}
	return nil
}
