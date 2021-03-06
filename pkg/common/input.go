package common

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func gamepadUp() bool {
	return ebiten.GamepadAxis(0, 1) < -0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton11)
}

func gamepadDown() bool {
	return ebiten.GamepadAxis(0, 1) > 0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton13)
}

func gamepadLeft() bool {
	return ebiten.GamepadAxis(0, 0) < -0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton14)
}

func gamepadRight() bool {
	return ebiten.GamepadAxis(0, 0) > 0.1 || ebiten.IsGamepadButtonPressed(0, ebiten.GamepadButton12)
}

func movingUp() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyW) || gamepadUp()
}

func movingDown() bool {
	return ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyS) || gamepadDown()
}

func movingLeft() bool {
	return ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyA) || gamepadLeft()
}

func movingRight() bool {
	return ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyD) || gamepadRight()
}

func accept() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton1)
}

//TODO: Don't forget the gamepad
func back() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyShift) || inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyTab)
}

func start() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}
