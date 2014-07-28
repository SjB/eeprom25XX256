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

func main() {
	flag.Parse()
	in := os.Stdin

	buf := make([]byte, 32768)
	in.Read(buf)

	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	bus := embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)
	defer bus.Close()

	eeprom := eeprom25XX256.New(bus)

	if _, err := eeprom.Write(buf); err != nil {
		panic(err)
	}
}
