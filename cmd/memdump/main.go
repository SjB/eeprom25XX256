package main

import (
	"flag"
	"os"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"

	"github.com/SjB/eeprom25XX256"
)

const (
	channel = 0
	speed   = 500000
	bpw     = 8
	delay   = 0
)

func detectHost() {
	board := os.Getenv("EMBD_HOST")
	if board == "RPi" || board == "RPi2" {
		embd.SetHost(embd.HostRPi, 2)
	} else if board == "BBB" {
		embd.SetHost(embd.HostBBB, 0)
	}
}

func main() {
	flag.Parse()
	out := os.Stdout

	detectHost()

	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	bus := embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)
	defer bus.Close()

	eeprom := eeprom25XX256.New(bus)

	buf := make([]byte, 32768)
	if _, err := eeprom.Read(buf); err != nil {
		panic(err)
	}
	out.Write(buf)
}
