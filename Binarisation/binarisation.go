package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Y struct {
	Setosa     float64 `json:"Setosa"`
	Versicolor float64 `json:"Versicolor"`
	Virginica  float64 `json:"Virginica"`
}

type Iris struct {
	SepalLength    string `json:"sepal_lenght"`
	SepalWidth     string `json:"sepal_width"`
	PetalLength    string `json:"petal_length"`
	PetalWidth     string `json:"petal_width"`
	Classification string `json:"classification"`
}

func Reciver(iris chan []Iris) {
	ln, _ := net.Listen("tcp", "localhost:8002") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	irisAux := make([]Iris, 0)
	dec := json.NewDecoder(con)
	dec.Decode(&irisAux)
	//fmt.Print(irisAux)
	iris <- irisAux
	close(iris)
	return
}

func uniqueElementsStrings(y []string) []string {
	valoresUnicos := make([]string, 0)
	valoresUnicos = append(valoresUnicos, y[0])
	for i := 1; i < len(y); i++ {
		existe := false
		for j := 0; j < len(valoresUnicos); j++ {
			if y[i] == valoresUnicos[j] {
				existe = true
			}
		}
		if !existe {
			valoresUnicos = append(valoresUnicos, y[i])
		}
	}
	return valoresUnicos
}

func categorizeUniqueStrings(unique, inputArray []string) []float64 {
	categorize := make([]float64, 0)
	for i := 0; i < len(inputArray); i++ {
		for j := 0; j < len(unique); j++ {
			if inputArray[i] == unique[j] {
				categorize = append(categorize, float64(j))
				break
			}
		}
	}
	return categorize
}

func binarizarCategorias(y []float64, cant int) [][]float64 {
	catBin := make([][]float64, 0)
	for i := 0; i < cant; i++ {
		aux := make([]float64, 0)
		for j := 0; j < len(y); j++ {
			if i == int(y[j]) {
				aux = append(aux, 1.0)
			} else {
				aux = append(aux, 0.0)
			}
		}
		catBin = append(catBin, aux)
	}
	return catBin
}

func Send(y []Y) {
	conn, _ := net.Dial("tcp", "localhost:8000")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(y)
}

func main() {
	irisChanel := make(chan []Iris)
	go Reciver(irisChanel)
	iris, _ := <-irisChanel
	// fmt.Println(iris)
	y := make([]string, 0)
	for _, val := range iris {
		y = append(y, val.Classification)
	}
	catY := categorizeUniqueStrings(uniqueElementsStrings(y), y)
	binCat := binarizarCategorias(catY, len(uniqueElementsStrings(y)))
	// fmt.Print(binCat)
	yf := make([]Y, 0)
	var aux Y
	// for _, val := range binCat {
	// 	aux.Setosa = val[0]
	// 	aux.Versicolor = val[1]
	// 	aux.Virginica = val[2]
	// 	yf = append(yf, aux)
	// }
	for i := 0; i < len(binCat[0]); i++ {
		aux.Setosa = binCat[0][i]
		aux.Versicolor = binCat[1][i]
		aux.Virginica = binCat[2][i]
		yf = append(yf, aux)
	}
	Send(yf)
	fmt.Println("Completado")
}
