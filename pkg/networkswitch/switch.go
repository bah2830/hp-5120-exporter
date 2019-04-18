package networkswitch

import "fmt"

type Switch interface {
	GetEnvironmentDetails() (*EnvironmentDetails, error)
	Disconnect()
}

type EnvironmentDetails struct {
	Sensors []Sensor
}

type Sensor struct {
	Name        string
	TempCelsius int
	Limits      Limits
}

type Limits struct {
	Lower    int
	Warning  int
	Alarm    int
	Critical int
}

func (s Sensor) String() string {
	return fmt.Sprintf("%s: %d", s.Name, s.TempCelsius)
}
