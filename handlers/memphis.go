package handlers

import (
	"fmt"
	"os"
	"sync"

	"github.com/memphisdev/memphis.go"
)

func Produce(wg *sync.WaitGroup, host, user, token, station string, count int64, msg []byte, sync bool) {
	conn, err := memphis.Connect(host, user, token)
	if err != nil {
		fmt.Println("Can not connect with memphis: " + err.Error())
		os.Exit(0)
	}
	defer conn.Close()
	p, err := conn.CreateProducer(station, "p", memphis.ProducerGenUniqueSuffix())
	if err != nil {
		fmt.Println("Can not create producer: " + err.Error())
		os.Exit(0)
	}

	if sync {
		for i := 0; i < int(count); i++ {
			p.Produce(msg)
		}
	} else {
		for i := 0; i < int(count); i++ {
			p.Produce(msg, memphis.AsyncProduce())
		}
	}
	wg.Done()
}
