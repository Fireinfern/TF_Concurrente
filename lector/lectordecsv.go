package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type Iris struct {
	SepalLength    string `json:"sepal_lenght"`
	SepalWidth     string `json:"sepal_width"`
	PetalLength    string `json:"petal_length"`
	PetalWidth     string `json:"petal_width"`
	Classification string `json:"classification"`
}

func send(remote string, dataset []Iris) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(dataset)
}

func Opencsv(file string) [][]string {
	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error, not a valid file", err)
	}
	defer csvfile.Close()
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
	dataset := []Iris{}
	var iris Iris
	for _, rec := range table[1:] {
		iris.SepalLength = rec[0]
		iris.SepalWidth = rec[1]
		iris.PetalLength = rec[2]
		iris.PetalWidth = rec[3]
		iris.Classification = rec[4]
		dataset = append(dataset, iris)
		// fmt.Println(iris)
	}
	// fmt.Println(dataset)
	// jsonData, err := json.MarshalIndent(dataset, "", "\t")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	con, _ := net.Dial("tcp", "localhost:8000")
	defer con.Close()
	enc := json.NewEncoder(con)
	enc.Encode(dataset)
	//enc := json.NewEncoder(os.Stdout)
	//enc.Encode(dataset)
}
