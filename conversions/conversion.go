package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

func convertStringToFloat(number string) float64 {
	s, err := strconv.ParseFloat(number, 64)
	if err != nil {
		fmt.Println("fatal error ", err)
	}
	return s
}

// SetConvertionInArray Convert a value in an string array to float64
func SetConvertionInArray(org []string, pos int) float64 {
	return convertStringToFloat(org[pos])
}

// ConvertStringArrayToFloatArray converts an array of strings to an array of float64
func ConvertStringArrayToFloatArray(stringArray []string) []float64 {
	var newArray []float64
	for i := 0; i < len(stringArray); i++ {
		newArray = append(newArray, SetConvertionInArray(stringArray, i))
	}
	return newArray
}

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

func Recive(irisChanel chan []Iris) {
	ln, _ := net.Listen("tcp", "localhost:8001") // <- crea conexion local
	defer ln.Close()
	con, _ := ln.Accept() // <- crea servidor
	defer con.Close()
	var iris []Iris
	dec := json.NewDecoder(con)
	dec.Decode(&iris)
	irisChanel <- iris
	close(irisChanel)
}

func SendConvertion(x []X) {
	conn, _ := net.Dial("tcp", "localhost:8000")
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(x)
}

func main() {
	irisChanel := make(chan []Iris)
	go Recive(irisChanel)
	iris, _ := <-irisChanel
	//fmt.Print(iris)
	x := make([][]float64, 0)
	for _, val := range iris {
		s := make([]string, 0)
		s = append(s, val.SepalLength, val.SepalWidth, val.PetalLength, val.PetalWidth)
		x = append(x, ConvertStringArrayToFloatArray(s))
	}
	//fmt.Print(x)
	var trainX []X
	var aux X
	for _, val := range x {
		aux.SepalLength = val[0]
		aux.SepalWidth = val[1]
		aux.PetalLength = val[2]
		aux.PetalWidth = val[3]
		trainX = append(trainX, aux)
	}
	SendConvertion(trainX)
	fmt.Print("Completado")
	// enc := json.NewEncoder(conn)
	// enc.Encode(trainX)
}
