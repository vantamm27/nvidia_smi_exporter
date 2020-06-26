package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {

	metricList := []string{
		"temperature.gpu", "utilization.gpu",
		"utilization.memory", "memory.total",
		"memory.free", "memory.used",
		"power.draw", "fan.speed"}
	strQuery := "name,index,"
	for i := 0; i < len(metricList); i++ {
		strQuery += metricList[i] + ","
	}
	strQuery = strings.Trim(strQuery, ",")
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu="+strQuery,
		"--format=csv,noheader,nounits").Output()

	// out, err := exec.Command(
	// 	"nvidia-smi",
	// 	"--query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,power.draw,fan.speed",
	// 	"--format=csv,noheader,nounits").Output()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	result := ""
	for _, row := range records {
		name := fmt.Sprintf("%s[%s]", row[0], row[1])
		for idx, value := range row[2:] {
			result = fmt.Sprintf(
				"%s%s{gpu=\"%s\"} %s\n", result,
				strings.Replace(metricList[idx], ".", "_", -1), name, value)
		}
	}

	fmt.Fprintf(response, result)
}

func main() {
	addr := ":9101"
	if len(os.Args) > 1 {
		addr = ":" + os.Args[1]
	}

	http.HandleFunc("/metrics", metrics)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
