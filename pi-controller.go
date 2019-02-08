package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func signalToFire() {
	// Connect to Raspberry PI gpio (i.e. pins on the board)
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	// Setup pins connected to the relays responsible for spinning the launch motor and the firing motor
	spinPin := rpio.Pin(21)
	firePin := rpio.Pin(20)
	spinPin.Output()
	firePin.Output()

	// Spin up the firing motor so it is ready to fire
	spinPin.High()
	fmt.Println("gpio SPIN")
	time.Sleep(time.Second)

	// Fire the ball by sending the singal to start the fire relay
	firePin.High()
	fmt.Println("gpio FIRE")
	time.Sleep(time.Second)

	// Turn off the motors
	firePin.Low()
	spinPin.Low()

	fmt.Println("gpio OFF")
}
