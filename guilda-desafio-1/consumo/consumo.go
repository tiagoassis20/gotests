package consumo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Consumo struct {
	Kms    float64
	Litros float64
}
type ConsumoData []Consumo

func (c *ConsumoData) Load(filename string) error {

	if err := checkFile(filename); err != nil {
		return err
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", file)
	if file == nil {
		return nil
	}
	if err := json.Unmarshal(file, &c); err != nil {
		return err
	}

	return nil
}
func (c ConsumoData) Add(consumo Consumo) {
	c = append(c, consumo)
	fmt.Println(c)
}

func (c ConsumoData) Save(filename string) error {
	// Preparing the data to be marshalled and written.
	dataBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, dataBytes, 0644); err != nil {
		return err
	}
	return nil
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		f.Write([]byte("[]"))
	}
	return nil
}
