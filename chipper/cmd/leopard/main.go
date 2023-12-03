package main

import (
	"github.com/BrutalHex/wire-pod/chipper/pkg/initwirepod"
	stt "github.com/BrutalHex/wire-pod/chipper/pkg/wirepod/stt/leopard"
)

func main() {
	initwirepod.StartFromProgramInit(stt.Init, stt.STT, stt.Name)
}
