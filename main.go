package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

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

// Opencsv Open a CSV file as a slice of slices of strings
func Opencsv(file string) [][]string {
	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error, not a valid file", err)
	}
	r := csv.NewReader(csvfile)
	table := make([][]string, 0)
	i := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("fatal error", err)
		}
		//fmt.Println(record[0], " ", record[1])
		table = append(table, record)
		i++
	}
	//fmt.Print(table)
	return table
}

//
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

// DownloadFile gets a file and downloading using the file path and the url of the file
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	err := DownloadFile("iris.csv", "https://raw.githubusercontent.com/Fireinfern/Datasets/master/iris.csv")
	if err != nil {
		fmt.Println("Error")
	}
	table := Opencsv("iris.csv")
	X := make([][]float64, 0)
	Y := make([]string, 0)
	for i := 1; i < len(table); i++ {
		X = append(X, ConvertStringArrayToFloatArray(table[i][:len(table[i])-2]))
		Y = append(Y, table[i][len(table[i])-1])
	}
	catY := categorizeUniqueStrings(uniqueElementsStrings(Y), Y)
	// Prints de prueba para conversion y categorizacion
	// fmt.Println(catY)
	l_rate := 0.1
	n_epoch := 5
	binCat := binarizarCategorias(catY, len(uniqueElementsStrings(Y)))
	peso1 := make(chan []float64)
	peso2 := make(chan []float64)
	peso3 := make(chan []float64)
	go obtenerPesos(X, binCat[0], l_rate, n_epoch, peso1)
	go obtenerPesos(X, binCat[1], l_rate, n_epoch, peso2)
	go obtenerPesos(X, binCat[2], l_rate, n_epoch, peso3)
	pesos1, ok1 := <-peso1
	pesos2, ok2 := <-peso2
	pesos3, ok3 := <-peso3
	fmt.Println(pesos1, ok1)
	fmt.Println(pesos2, ok2)
	fmt.Println(pesos3, ok3)

}
