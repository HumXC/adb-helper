package adb

import "fmt"

func (d *Device) Press(x, y, hold int) (err error) {
	return d.Swipe(x, y, x, y, hold)
}

func (d *Device) Click(x, y int) (err error) {
	return d.Press(x, y, 50)
}

func (d *Device) Swipe(startX, startY, endX, endY, hold int) error {
	_, err := d.runner(fmt.Sprintf(d.preArg+"shell input touchscreen swipe %d %d %d %d %d", startX, startY, endX, endY, hold))
	return err
}

func (d *Device) Text(str string) (err error) {
	_, err = d.runner(d.preArg + "shell input text " + str)
	return
}

func (d *Device) Power() error {
	return d.Keyevent("POWER")
}

func (d *Device) VolumeUp() error {
	return d.Keyevent("VOLUME_UP")
}

func (d *Device) VolumeDown() error {
	return d.Keyevent("VOLUME_DOWN")
}

func (d *Device) Del(count int) error {
	var err error
	for j := 0; j < count; j++ {
		err = d.Keyevent("DEL")
		if err != nil {
			return err
		}
	}
	return err
}

func (d *Device) Home() error {
	return d.Keyevent("HOME")
}

func (d *Device) Back() error {
	return d.Keyevent("BACK")
}

func (d *Device) Menu() error {
	return d.Keyevent("MENU")
}

func (d *Device) Keyevent(keycode string) (err error) {
	_, err = d.runner(d.preArg + "shell input keyevent KEYCODE_" + keycode)
	return
}
