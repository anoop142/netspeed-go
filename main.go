/*
Anoop S
Simple tool to print network speed written in go
Made for personal dwmblocks script
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	sysPath = "/sys/class/net/"
	rxFile  = "statistics/rx_bytes"
	txFile  = "statistics/tx_bytes"
)

func getIntFromFile(filePath string, intRead *uint64) {
	File, err := os.Open(filePath)

	if err != nil {
		log.Fatalln(err)
	}

	defer File.Close()

	fmt.Fscanf(File, "%d", intRead)

}

func printStats(device string, up, down bool) {
	var txBytesOld, txBytesNew, rxBytesOld, rxBytesNew uint64

	// check if interface exist
	if _, err := os.Stat(sysPath + device); os.IsNotExist(err) {
		log.Fatalf("interface %v not found\n", device)
	}

	rxFilePath := sysPath + device + "/" + rxFile
	txFilePath := sysPath + device + "/" + txFile

	getIntFromFile(rxFilePath, &rxBytesOld)
	getIntFromFile(txFilePath, &txBytesOld)

	// sleep for 1 second
	time.Sleep(1 * time.Second)

	getIntFromFile(rxFilePath, &rxBytesNew)
	getIntFromFile(txFilePath, &txBytesNew)

	// calculate speed
	downSpeed := rxBytesNew - rxBytesOld
	upSpeed := txBytesNew - txBytesOld

	// skip negligible changes
	if downSpeed < 1024 {
		downSpeed = 0
	}

	if upSpeed < 1024 {
		upSpeed = 0
	}
	// convert to human friendly format
	downSpeedString := humanize.Bytes(downSpeed)
	upSpeedString := humanize.Bytes(upSpeed)

	// print both if none specified
	if (up && down) || (!up && !down) {
		fmt.Printf("%v:%v", upSpeedString, downSpeedString)
	} else if up {
		fmt.Printf("%v", upSpeedString)
	} else {
		fmt.Printf("%v", downSpeedString)
	}

}

func main() {
	device := flag.String("i", "", "interface")
	up := flag.Bool("up", false, "up stat")
	down := flag.Bool("down", false, "down stat")

	flag.Parse()

	if *device != "" {
		printStats(*device, *up, *down)
	} else {
		fmt.Printf("usage: netspeed-go -i interface [-up] [-down]\n")
		fmt.Printf("\nexample: netspeed-go -i wlan0 -down\n\n")
		flag.PrintDefaults()
	}

}
