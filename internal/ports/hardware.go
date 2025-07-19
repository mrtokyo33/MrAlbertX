package ports

type LightControllerPort interface {
	SetLightState(lightID string, isOn bool) error
}