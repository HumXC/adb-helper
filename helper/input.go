package helper

import "fmt"

type InputMan struct {
	adbRunner ADBRunner
}

func (i *InputMan) Press(x, y, hold int) (err error) {
	return i.Swipe(x, y, x, y, hold)
}

func (i *InputMan) Click(x, y int) (err error) {
	return i.Press(x, y, 50)
}

func (i *InputMan) Swipe(startX, startY, endX, endY, hold int) error {
	out, err := i.adbRunner(fmt.Sprintf("shell input touchscreen swipe %d %d %d %d %d", startX, startY, endX, endY, hold))
	if err == nil {
		err = ChexkError(out)
	}
	return err
}

func (i *InputMan) Text(str string) (err error) {
	_, err = i.adbRunner("shell input text " + str)
	return
}

func (i *InputMan) Power() error {
	return i.Keyevent("POWER")
}

func (i *InputMan) VolumeUp() error {
	return i.Keyevent("VOLUME_UP")
}

func (i *InputMan) VolumeDown() error {
	return i.Keyevent("VOLUME_DOWN")
}

func (i *InputMan) Del(count int) error {
	var err error
	for j := 0; j < count; j++ {
		err = i.Keyevent("DEL")
		if err != nil {
			return err
		}
	}
	return err
}

func (i *InputMan) Home() error {
	return i.Keyevent("HOME")
}

func (i *InputMan) Back() error {
	return i.Keyevent("BACK")
}

func (i *InputMan) Menu() error {
	return i.Keyevent("MENU")
}

func (i *InputMan) Keyevent(keycode string) (err error) {
	_, err = i.adbRunner("shell input keyevent KEYCODE_" + keycode)
	return
}
