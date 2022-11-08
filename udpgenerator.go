package main

/*
UDP generator
Version 0.0.1
Date : 07-11-2022

*/

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Data struct {
	data          [][]byte
	chunksize     int
	typeoperation string
	filename      string
	rawlendata    int
}

func (d *Data) SplitData(adata *[]byte) error {
	alldata := *adata

	for i := 0; i < len(alldata); i += d.chunksize {

		end := i + d.chunksize

		if end > len(alldata) {
			end = len(alldata) - 1
		}

		d.data = append(d.data, alldata[i:end])
	}

	if len(d.data) == 0 {
		return errors.New("Can't split alldata")
	} else {
		return nil
	}

}

func (d *Data) RandomData(lendata int) ([]byte, error) {
	alldata := make([]byte, lendata)
	_, err := rand.Read(alldata)
	if err != nil {
		return nil, errors.New("Error in creating alldata`")
	}

	return alldata, nil

}

func (d *Data) ReadRawFile() ([][]byte, error) {
	var alldata []byte
	var err error
	if d.rawlendata > 0 {
		alldata, err = d.RandomData(d.rawlendata)
	} else {

		alldata, err = os.ReadFile(d.filename)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = d.SplitData(&alldata)

	if err != nil {
		log.Println(err)
		os.Exit(5)
	}

	return d.data, nil

}

func main() {
	var port = flag.Int("p", 5555, "Set port name")
	var chunksize = flag.Int("c", 100, "Set chunksize name")
	var server = flag.String("s", "10.199.100.100", "Set server ")
	//file := flag.NewFlagSet("file", flag.ExitOnError)
	//filename := file.String("n", "", "Set text file ")
	filename := flag.String("f", "", "Set text file ")

	//raw := flag.NewFlagSet("raw", flag.ExitOnError)
	//sizeraw := raw.Int("l", 0, "Set size length")
	sizeraw := flag.Int("l", 0, "Set size length")

	flag.Parse()

	log.Println("Currer size raw", *sizeraw)
	data := Data{}

	data.chunksize = *chunksize

	CONNECT := fmt.Sprintf("%s:%d", *server, *port)
	log.Println(CONNECT)

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()
	if *sizeraw > 0 {
		data.rawlendata = *sizeraw

	} else {
		data.filename = *filename
	}
	data.data, err = data.ReadRawFile()
	if err != nil {
		log.Println(err)
	}

	for index := 0; index < len(data.data); index++ {
		i, err := c.Write(data.data[index])

		if err != nil {
			log.Println(err)
			os.Exit(4)
		}

		log.Printf("Chunk #%d Sended %d bytes", index, i)
	}

}
