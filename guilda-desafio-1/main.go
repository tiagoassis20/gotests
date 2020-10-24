package main

import (
	"flag"
	"fmt"

	"example.com/guilda/golang/desafio1/consumo"
)

func main() {
	kms := flag.Float64("km", 0, "kilometragem atual do carro")
	litros := flag.Float64("litros", 0, "quantidade de litros na bomba")
	flag.Parse()
	filename := "consumo.json"
	novoConsumo := consumo.Consumo{Kms: *kms, Litros: *litros}
	consumoData := make(consumo.ConsumoData, 1)
	if err := consumoData.Load(filename); err != nil {
		panic(err)
	}
	consumoData = append(consumoData, novoConsumo)
	fmt.Println((consumoData))
	if err := consumoData.Save(filename); err != nil {
		panic(err)
	}
}
