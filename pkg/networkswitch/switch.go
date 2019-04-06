package networkswitch

type Switch interface {
	GetEnvironmentDetails() (EnvironmentDetails, error)
}

type EnvironmentDetails interface {
	GetTemperatureCelicius() int16
}
