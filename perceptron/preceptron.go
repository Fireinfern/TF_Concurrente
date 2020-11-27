package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
)

type Iris struct {
	SepalLength    string `json:"sepal_lenght"`
	SepalWidth     string `json:"sepal_width"`
	PetalLength    string `json:"petal_length"`
	PetalWidth     string `json:"petal_width"`
	Classification string `json:"classification"`
}

type X struct {
	SepalLength float64 `json:"sepal_length"`
	SepalWidth  float64 `json:"sepal_width"`
	PetalLength float64 `json:"petal_length"`
	PetalWidth  float64 `json:"petal_width"`
}

type Y struct {
	Setosa     float64 `json:"Setosa"`
	Versicolor float64 `json:"Versicolor"`
	Virginica  float64 `json:"Virginica"`
}

func predecir(fila, pesos []float64) float64 {
	activador := pesos[0]
	for i := 0; i < len(fila); i++ {
		activador = activador + (pesos[i+1] * fila[i])
	}
	if activador >= 0.0 {
		return 1.0
	} else {
		return 0.0
	}
}

func entrenarPesos(entrenamiento_X [][]float64, entrenamiento_Y []float64, l_rate float64, n_epoch int) []float64 {
	pesos := make([]float64, 0)
	for i := 0; i < len(entrenamiento_X[0])+1; i++ {
		pesos = append(pesos, 0.0)
	}
	for epoch := 0; epoch < (n_epoch); epoch++ {
		sum_error := 0.0
		for i, row := range entrenamiento_X {
			prediction := predecir(row, pesos)
			error := entrenamiento_Y[i] - prediction
			sum_error = sum_error + math.Pow(error, 2.0)
			pesos[0] = pesos[0] + l_rate*error
			for j := 0; j < len(row); j++ {
				pesos[j+1] = pesos[j+1] + l_rate*error*row[j]
			}
		}
	}
	return pesos
}

func obtenerPesos(X [][]float64, y []float64, l_rate float64, n_epochs int, c chan []float64) {
	c <- entrenarPesos(X, y, l_rate, n_epochs)
	close(c)
}

func GetX(xChanel chan []X) {
	ln, _ := net.Listen("tcp", "localhost:8004") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	xAux := make([]X, 0)
	dec := json.NewDecoder(con)
	dec.Decode(&xAux)
	//fmt.Print(irisAux)
	xChanel <- xAux
	close(xChanel)
	return
}

func GetY(yChanel chan []Y) {
	ln, _ := net.Listen("tcp", "localhost:8005") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	yAux := make([]Y, 0)
	dec := json.NewDecoder(con)
	dec.Decode(&yAux)
	//fmt.Print(irisAux)
	yChanel <- yAux
	close(yChanel)
	return
}

func main() {
	xChanel := make(chan []X)
	yChanel := make(chan []Y)
	go GetX(xChanel)
	go GetY(yChanel)
	x, _ := <-xChanel
	y, _ := <-yChanel

	// Tratar X para procesamiento
	xTrain := make([][]float64, 0)
	for _, val := range x {
		aux := make([]float64, 0)
		aux = append(aux, val.SepalLength, val.SepalWidth, val.PetalLength, val.PetalWidth)
		xTrain = append(xTrain, aux)
	}
	var y1, y2, y3 []float64
	for _, val := range y {
		y1 = append(y1, val.Setosa)
		y2 = append(y2, val.Versicolor)
		y3 = append(y3, val.Virginica)
	}
	fmt.Println(y1)
	fmt.Println(y2)
	fmt.Println(y3)
	l_rate := 0.1
	n_epoch := 5
	peso1 := make(chan []float64)
	peso2 := make(chan []float64)
	peso3 := make(chan []float64)
	go obtenerPesos(xTrain, y1, l_rate, n_epoch, peso1)
	go obtenerPesos(xTrain, y2, l_rate, n_epoch, peso2)
	go obtenerPesos(xTrain, y3, l_rate, n_epoch, peso3)
	// pesos1, _ := <-peso1
	// pesos2, _ := <-peso2
	// pesos3, _ := <-peso3
	println(<-peso1)
	println(<-peso2)
	println(<-peso3)
}
