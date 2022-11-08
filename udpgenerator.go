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
	"time"
)

type Data struct {
	data          [][]byte
	chunksize     int
	typeoperation string
	filename      string
	rawlendata    int
	connect       string
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

func (d *Data) CreateConnect(CONNECT string) (*net.UDPConn, error) {

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	log.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	return c, nil

}

func (d *Data) SendInfiniteTraffic() error {

	c, err := d.CreateConnect(d.connect)
	if err != nil {
		return errors.New("Can't create connect")

	}
	defer c.Close()
	chunks, err := d.RandomData(d.chunksize)
	if err != nil {
		log.Println(err)
		os.Exit(0)

	}
	d.data = append(d.data, chunks)
	index := 0
	for {
		index++

		i, err := c.Write(d.data[0])
		time.Sleep(10 * time.Microsecond)

		if err != nil {
			log.Println(err)
			os.Exit(4)
		}
		log.Printf("Chunk #%d Sended %d bytes", index, i)
	}
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
	threads := flag.Int("t", 1, "Set threads")

	flag.Parse()

	log.Println("Currer size raw and threads", *sizeraw, threads)
	data := Data{}

	data.connect = fmt.Sprintf("%s:%d", *server, *port)
	data.chunksize = *chunksize
	if *threads == 1 {
		if *sizeraw > 0 {
			data.rawlendata = *sizeraw

		} else {
			data.filename = *filename
		}

		if *sizeraw != -1 {
			c, err := data.CreateConnect(data.connect)
			if err != nil {
				log.Println("Can't create connect")
				os.Exit(5)

			}
			defer c.Close()
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
		} else {
			err := data.SendInfiniteTraffic()
			if err != nil {
				log.Panicln(err)
				os.Exit(5)

			}

		}
	} else {

	}

}
