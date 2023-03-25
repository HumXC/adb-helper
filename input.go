package adb

import (
	"fmt"
)

type Input interface {
	Press(x, y, duration int) error
	Power() error
	Click(x, y int) error
	Swipe(x, y, x2, y2, duration int) error
	Text(text string) error
	VolumeDown() error
	VolumeUp() error
	Home() error
	Back() error
	Menu() error
	Del(count int) error
	Keyevent(event string) error
}

type input struct {
	cmd ADBRunner
}

func (i *input) Press(x, y, duration int) (err error) {
	return i.Swipe(x, y, x, y, duration)
}

func (i *input) Click(x, y int) (err error) {
	return i.Press(x, y, 50)
}

func (i *input) Swipe(startX, startY, endX, endY, duration int) error {
	_, err := i.cmd(fmt.Sprintf("shell input touchscreen swipe %d %d %d %d %d", startX, startY, endX, endY, duration))
	return err
}

func (i *input) Text(str string) (err error) {
	_, err = i.cmd("shell input text " + str)
	return
}

func (i *input) Power() error {
	return i.Keyevent("POWER")
}

func (i *input) VolumeUp() error {
	return i.Keyevent("VOLUME_UP")
}

func (i *input) VolumeDown() error {
	return i.Keyevent("VOLUME_DOWN")
}

func (i *input) Del(count int) error {
	var err error
	for j := 0; j < count; j++ {
		err = i.Keyevent("DEL")
		if err != nil {
			return err
		}
	}
	return err
}

func (i *input) Home() error {
	return i.Keyevent("HOME")
}

func (i *input) Back() error {
	return i.Keyevent("BACK")
}

func (i *input) Menu() error {
	return i.Keyevent("MENU")
}

func (i *input) Keyevent(keycode string) (err error) {
	_, err = i.cmd("shell input keyevent KEYCODE_" + keycode)
	return
}
