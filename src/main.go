package main

import (
	Versioner "github.com/FraktalDeFiDAO/Versioner/Versioner"
)

func main() {
	versioner := *Versioner.NewVersioner()

	versioner.Run()
}
