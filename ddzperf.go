package main

import (
	"ddzperf/worker"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	//	"time"
)

var para = 200

//var timer = time.Duration(10)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cs := make(chan os.Signal, 1)
	signal.Notify(cs, os.Interrupt, os.Kill)
	dataQ := make(chan worker.SumBase, 5000)
	fileQ := make(chan worker.Fdata, 100)
	fn := "p" + strconv.Itoa(para) + ".csv"
	csv, err := os.OpenFile(fn, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err == nil {
		go worker.FileWorker(csv, fileQ)
	} else {
		fmt.Println("File open error: ", err)
	}
	for i := 0; i < para; i++ {
		go worker.ReqWorker(dataQ)
	}
	go worker.SumWorker(dataQ, fileQ)
	//	t := time.After(timer * time.Minute)
	//	for {
	//		if len(t) == 1 {
	//			<-t
	//			break
	//		}
	//	}
	s := <-cs
	fmt.Println("\nSignal received:", s)
	defer csv.Close()
}
