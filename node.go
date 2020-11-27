package main

import (
	"encoding/json"
	"fmt"
	"net"
)

var (
	local string
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

type CompleteWeights struct {
	Setosa     []float64 `json:"Setosa"`
	Versicolor []float64 `json:"Versicolor"`
	Virginica  []float64 `json:"Virginica"`
}

func GetBasicIris(iris chan []Iris) {
	ln, _ := net.Listen("tcp", "localhost:8000") // <- crea conexion local
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

func SendIrisToConvert(irisData []Iris) {
	conn, _ := net.Dial("tcp", "localhost:8001")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(irisData)
}

func SendIrisToBin(irisData []Iris) {
	conn, _ := net.Dial("tcp", "localhost:8002")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(irisData)
}

func SendX(x []X) {
	conn, _ := net.Dial("tcp", "localhost:8004")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(x)
}

func SendY(y []Y) {
	conn, _ := net.Dial("tcp", "localhost:8005")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(y)
}

func GetTrainingData(xChanel chan []X) {
	ln, _ := net.Listen("tcp", "localhost:8000") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	x := make([]X, 0)
	dec := json.NewDecoder(con)
	dec.Decode(&x)
	xChanel <- x
	close(xChanel)
	return
}

func GetBinCat(yChanel chan []Y) {
	ln, _ := net.Listen("tcp", "localhost:8000") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	y := make([]Y, 0)
	dec := json.NewDecoder(con)
	dec.Decode(&y)
	yChanel <- y
	close(yChanel)
	return
}

func GetWeights(weightsChanel chan CompleteWeights) {
	ln, _ := net.Listen("tcp", "localhost:8000") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	var auxWeight CompleteWeights
	dec := json.NewDecoder(con)
	dec.Decode(&auxWeight)
	weightsChanel <- auxWeight
	close(weightsChanel)
	return
}

func main() {
	iris := make(chan []Iris)
	go GetBasicIris(iris)
	irisData, ok := <-iris
	if ok {
		// Esto es porque renzo ravelli no dejaba de decirme que lo escriba bien
		fmt.Println("Iris R3c1vid0")
	}
	SendIrisToConvert(irisData)
	xChanel := make(chan []X)
	go GetTrainingData(xChanel)
	x, _ := <-xChanel
	fmt.Println(x)
	SendIrisToBin(irisData)
	yChanel := make(chan []Y)
	go GetBinCat(yChanel)
	y, _ := <-yChanel
	fmt.Println(y)
	SendX(x)
	SendY(y)
	weightsChanel := make(chan CompleteWeights)
	go GetWeights(weightsChanel)
	weights := <-weightsChanel
	fmt.Println(weights)
}
