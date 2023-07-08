package overwatch

import "log"

type LogAlerter struct{}

func NewLogAlerter() *LogAlerter {
	return &LogAlerter{}
}

func (a *LogAlerter) ExcessPowerAlert() error {
	log.Print("Excess power is being generated. Don't let it get away!")
	return nil
}
