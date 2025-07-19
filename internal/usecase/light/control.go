package light

import (
	"MrAlbertX/server/internal/ports"
	"fmt"
)

type ControlUseCase struct {
	repo       ports.LightRepositoryPort
	controller ports.LightControllerPort
}

func NewControlUseCase(repo ports.LightRepositoryPort, controller ports.LightControllerPort) *ControlUseCase {
	return &ControlUseCase{
		repo:       repo,
		controller: controller,
	}
}

func (uc *ControlUseCase) Execute(lightID string, isOn bool) error {
	fmt.Printf("[BUSINESS LOGIC] Executing command for light '%s'.\n", lightID)
	light, err := uc.repo.FindByID(lightID)
	if err != nil {
		return err
	}
	if light == nil {
		return fmt.Errorf("light with id '%s' not found", lightID)
	}
	light.IsOn = isOn
	if err := uc.controller.SetLightState(light.ID, light.IsOn); err != nil {
		return err
	}
	return uc.repo.Update(*light)
}