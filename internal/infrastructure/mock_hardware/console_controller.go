package mock_hardware

import "fmt"

type ConsoleLightController struct{}

func NewConsoleLightController() *ConsoleLightController {
	return &ConsoleLightController{}
}

func (c *ConsoleLightController) SetLightState(lightID string, isOn bool) error {
	status := "OFF"
	if isOn {
		status = "ON"
	}

	fmt.Printf("[HARDWARE MOCK] Light '%s' was turned %s.\n", lightID, status)
	return nil
}
