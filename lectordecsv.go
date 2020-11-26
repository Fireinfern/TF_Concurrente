package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
	fmt.Print(table)
}
