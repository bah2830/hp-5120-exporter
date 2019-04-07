package networkswitch

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
