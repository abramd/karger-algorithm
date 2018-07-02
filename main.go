package main

import (
	"math/rand"
	"github.com/abramd/log"
	"sync"
	"fmt"
	"flag"
	"strings"
	"strconv"
)

type matrix [][]uint

//	m := matrix{{1, 0, 0, 0, 1, 0, 0, 1}, {1, 1, 0, 0, 0, 1, 0, 0}, {0, 1, 1, 0, 0, 0, 0, 0}, {0, 0, 1, 1, 0, 0, 1, 1}, {0, 0, 0, 1, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 0, 1, 0},}
// 1
//
//	m = matrix{{1, 1, 0}, {1, 0, 1}, {0, 1, 1}}
// 2

var argv struct {
	Matrix matrix
	Accuracy int
}

var mu sync.Mutex

func init() {
	tmp := flag.String("m", "", "matrix of graph. e.g. 1 0 1,0 0 1,1 0 0")
	flag.IntVar(&argv.Accuracy, "a", 100, "count of calculations")
	mm := strings.Split(*tmp, ",")
	for k, m := range mm {
		argv.Matrix[k] = make([]uint, 0)
		dd := strings.Split(strings.TrimSpace(m), " ")
		for _, d := range dd {
			digit, err := strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
			argv.Matrix[k] = append(argv.Matrix[k], uint(digit))
		}
	}
}

func main() {
	fmt.Println("input:", argv.Matrix)
	fmt.Println("count of calculations:", argv.Accuracy)
	var result uint
	for i := 0; i < argv.Accuracy; i++ {
		go func(){
			res := calc(argv.Matrix)
			mu.Lock()
			defer mu.Unlock()
			if result > res {
				result = res
			}
		}()
	}
	fmt.Println("result:", result)
}

func calc(m matrix) uint {

	cnt := len(m)
	for i := 0; i < cnt-2; i++ {
		fr := rand.Intn(len(m) - 1)
		f := m[fr]
		tmp := make([]int, 0)
		for k, v := range f {
			if v > 0 {
				tmp = append(tmp, k)
			}
		}

		// removing
		if len(tmp) == 0 {
			m = append(m[:fr], m[fr+1:]...)
			log.Errorln("point without any relations")
			continue
		}

		sr := tmp[rand.Intn(len(tmp)-1)]
		for k := 0; k < len(m); k++ {
			if m[k][sr] == 1 && k != fr {
				sr = k
				break
			}
		}
		s := m[sr]

		var edges = make([]int, 0)
		for k := range f {
			if s[k] > 0 && f[k] > 0 {
				f[k] = 0
				edges = append(edges, k)
			} else {
				f[k] = f[k] + s[k]
			}
		}

		m = m[:sr+copy(m[sr:], m[sr+1:])]
		c := len(m)
		for j := 0; j < c; j++ {
			for k := len(edges) - 1; k >= 0; k-- {
				v := edges[k]
				m[j] = append(m[j][:v], m[j][v+1:]...)
			}
		}
	}

	return uint(len(m[0]))
}
