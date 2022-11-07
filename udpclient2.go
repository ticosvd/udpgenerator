package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Data struct {
	data      [][]byte
	chunksize int
}

//func (m *Data) ReadFile()

func ReadFile(f *string) ([][]byte, error) {

	var data [][]byte

	chunksize := 100

	alldata, err := os.ReadFile(*f)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for i := 0; i < len(alldata); i += chunksize {

		end := i + chunksize
		if end > len(alldata) {
			end = len(alldata) - 1
		}

		data = append(data, alldata[i:end])
	}

	return data, nil

}

func main() {
	var port = flag.Int("p", 5555, "Set port name")
	var chunksize = flag.Int("c", 5555, "Set port name")
	var server = flag.String("s", "10.199.100.100", "Set server ")
	var file = flag.String("f", "1.txt", "Set text file ")

	flag.Parse()

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
	data.data, err = ReadFile(file)

	if err != nil {
		log.Println(err)
	}

	for index := 0; index < len(data.data); index++ {
		i, err := c.Write(data.data[index])

		if err != nil {
			log.Println(err)
			os.Exit(4)
		}

		log.Printf("Sended %d bytes", i)
	}

}
