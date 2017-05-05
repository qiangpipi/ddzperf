package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type SumBase struct {
	Res      int64
	Duration int64
}

type Fdata struct {
	TpsS int64
	TpsF int64
	Rq   int64
}

var ServerIp = "192.169.23.210"
var ServerPort = "80"
var BetUrl = "/mock/rankBet"
var Url = "http://192.169.23.210:7001/mock/rankBet"

func httpReq() []byte {
	c := &http.Client{}
	req, _ := http.NewRequest("GET", Url, strings.NewReader(""))
	res, err := c.Do(req)
	if err != nil {
		fmt.Println("HTTP Get Error:", err)
	} else {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
		} else {
			//			fmt.Println("RES Read end")
			return b
		}
		defer res.Body.Close()
	}
	return []byte{0}
}
func ReqWorker(tps chan<- SumBase) {
	var base SumBase
	//	c := &http.Client{}
	for {
		//		fmt.Println("HTTP Get Start.")
		start := time.Now().UnixNano()
		b := httpReq()
		end := time.Now().UnixNano()
		totalUs := (end - start) / 1000
		base.Duration = totalUs
		if strings.Contains(string(b), "\"000000\"") {
			base.Res = 0
		} else {
			base.Res = 1
		}
		tps <- base
		//		req, _ := http.NewRequest("GET", Url, strings.NewReader(""))
		//		start := time.Now().UnixNano()
		//		res, err := c.Do(req)
		//		if err != nil {
		//			fmt.Println("HTTP Get Error:", err)
		//		} else {
		//			fmt.Println("HTTP Get End.")
		//			end := time.Now().UnixNano()
		//			totalUs := (end - start) / 1000
		//			base.Duration = totalUs
		//			fmt.Println("RES Read start")
		//			b, err := ioutil.ReadAll(res.Body)
		//			if err != nil {
		//				fmt.Println("Error reading body:", err)
		//			} else {
		//				fmt.Println(strconv.Itoa(id), "RES Read end")
		//				if strings.Contains(string(b), "\"000000\"") {
		//					base.Res = 0
		//				} else {
		//					base.Res = 1
		//				}
		//			}
		//			tps <- base
		//		}
		//		defer res.Body.Close()
	}
}

func SumWorker(q <-chan SumBase, qf chan<- Fdata) {
	timer := time.After(5 * time.Second)
	countS := int64(0)
	countF := int64(0)
	duration := int64(0)
	for {
		d := <-q
		if d.Res == 0 {
			countS++
		} else if d.Res == 1 {
			countF++
		}
		duration += d.Duration
		//		fmt.Println("countS:", countS, "countF:", countF, "duration:", duration)
		if len(timer) == 1 {
			<-timer
			fmt.Println("Queue size:", len(q))
			data := Fdata{countS, countF, duration}
			qf <- data
			fmt.Println("countS:", countS, "countF:", countF, "duration:", duration, "AvgDura:", duration/((countS+countF)*1000))
			countS = 0
			countF = 0
			duration = 0
			timer = time.After(5 * time.Second)
		}
	}
}

func FileWorker(f *os.File, qf <-chan Fdata) {
	if _, err := f.WriteString("GoodT, BadT, Dura\n"); err != nil {
		fmt.Println("Error writing header:", err)
	} else {
		for {
			d := <-qf
			tmp := ""
			tmp += strconv.FormatInt(d.TpsS, 10)
			tmp += ", "
			tmp += strconv.FormatInt(d.TpsF, 10)
			tmp += ", "
			tmp += strconv.FormatInt(d.Rq, 10)
			tmp += "\n"
			if _, err := f.WriteString(tmp); err != nil {
				fmt.Println("Error writing file:", err)
			}
		}
	}
}
