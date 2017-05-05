package main

import (
	"fmt"
)

// RTP configuration
//Pattern 7-3
var Prizes = map[float64]float64{
	8:   0.02,
	3.2: 0.05,
	1.6: 0.1,
	0.8: 0.2,
	0.4: 0.4,
	0:   0.23,
}

//Pattern 7-1
//var Prizes = map[float64]float64{
//	10:  0.016,
//	5:   0.032,
//	2:   0.08,
//	1:   0.16,
//	0.5: 0.32,
//	0:   0.392,
//}

//Pattern 7-2
//var Prizes = map[float64]float64{
//	100: 0.0016,
//	50:  0.0032,
//	20:  0.008,
//	10:  0.016,
//	5:   0.032,
//	0:   0.9392,
//}

//Pattern 7-3
//var Prizes = map[float64]float64{
//	1.2: 0.133333,
//	1.1: 0.1455,
//	1:   0.16,
//	0.9: 0.1778,
//	0.8: 0.2,
//	0:   0.1834,
//}

var Num = 30

func calculateB(p map[float64]float64, init float64, n int) (pn map[float64]float64, res float64) {
	res = init
	if res == 1 {
		for k, v := range p {
			if k < 1 {
				res = res + (1-k)*v
			}
		}
		fmt.Println(res)
	}
	pn = make(map[float64]float64)
	if n-2 == 0 {
		return nil, res
	} else {
		kt := 0.0
		for kx, cx := range p {
			for ky, cy := range Prizes {
				if kx >= 1 {
					kt = kx - 1 + ky
				} else if kx < 1 {
					kt = ky
				}
				//				if kt == 1 {
				//					fmt.Println("kx:", kx, "cx:", cx, "ky:", ky, "cy:", cy)
				//				}
				if _, ok := pn[kt]; ok {
					pn[kt] = pn[kt] + cx*cy
				} else {
					pn[kt] = cx * cy
				}
			}
		}
		//		fmt.Println(pn)
		n = n - 1
		for k, v := range pn {
			if k < 1 {
				res = res + (1-k)*v
			}
		}
		//		fmt.Println(pn)
		fmt.Println(res)
		return calculateB(pn, res, n)
	}
}

func calculateA(p map[float64]float64, b float64, n int) (a float64) {
	rtp := 0.0
	for k, v := range p {
		rtp = rtp + k*v
	}
	baseDiff := 1 - rtp
	diff := float64(n) * baseDiff
	return b - diff
}

func main() {
	_, b := calculateB(Prizes, 1.0, Num)
	fmt.Println("b:", b)
	a := calculateA(Prizes, b, Num)
	fmt.Println("a:", a)
	fmt.Println("RPT:", a/b)
}
