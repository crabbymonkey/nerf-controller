package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func signalToFire() {
	// Connect to Raspberry PI gpio (i.e. pins on the board)
	logMessage("opening gpio")
	err := rpio.Open()
	if err != nil {
		fmt.Println("GPIO is not available")
		return
	}

	defer rpio.Close()

	// Setup pins connected to the relays responsible for spinning the launch motor and the firing motor
	spinPin := rpio.Pin(21)
	firePin := rpio.Pin(20)
	spinPin.Output()
	firePin.Output()

	// Spin up the firing motor so it is ready to fire
	spinPin.High()
	logMessage("gpio SPIN")
	time.Sleep(500 * time.Millisecond)

	// Fire the ball by sending the singal to start the fire relay
	firePin.High()
	logMessage("gpio FIRE")
	time.Sleep(500 * time.Millisecond)

	// Turn off the motors
	firePin.Low()
	spinPin.Low()

	logMessage("gpio OFF")
}
