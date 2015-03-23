package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/kasworld/idgen"
	"github.com/kasworld/idpos/idpos1m"
	"github.com/kasworld/idpos/idpos1s"
	"github.com/kasworld/idpos/idpos2m"
	"github.com/kasworld/idpos/idpos2s"
	"github.com/kasworld/idpos/idposi"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func writeHeapProfile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("mem profile %v", err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

func startCPUProfile(filename string) func() {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("cpu profile %v", err)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

func main() {
	log.Printf("go # %v", runtime.NumGoroutine())

	var cpuprofilename = flag.String("cpuprofilename", "", "cpu profile filename")
	var memprofilename = flag.String("memprofilename", "", "memory profile filename")
	flag.Parse()
	args := flag.Args()

	if *cpuprofilename != "" {
		fn := startCPUProfile(*cpuprofilename)
		defer fn()
	}

	sttime := time.Now()
	doMain(args)
	fmt.Printf("%v\n", time.Now().Sub(sttime))
	log.Printf("go # %v", runtime.NumGoroutine())

	if *memprofilename != "" {
		writeHeapProfile(*memprofilename)
	}
}

///

type posobj struct {
	ID  int64
	Pos [2]int
}

func (o *posobj) GetID() int64 {
	return o.ID
}
func (o *posobj) GetPos() [2]int {
	return o.Pos
}

func doMain(args []string) {
	xlen, ylen := 1024, 1024

	objs := make([]*posobj, 0, xlen*ylen)
	sttime := time.Now()
	for x := 0; x < xlen; x++ {
		for y := 0; y < ylen; y++ {
			objs = append(objs, &posobj{<-idgen.GenCh(), [2]int{x, y}})
		}
	}
	fmt.Printf("init %v\n", time.Now().Sub(sttime))

	idp1m := idpos1m.New(xlen, ylen)
	bench(objs, xlen, ylen, idp1m, "1d map")

	idp1s := idpos1s.New(xlen, ylen)
	bench(objs, xlen, ylen, idp1s, "1d slice")

	idp2m := idpos2m.New(xlen, ylen)
	bench(objs, xlen, ylen, idp2m, "2d map")

	idp2s := idpos2s.New(xlen, ylen)
	bench(objs, xlen, ylen, idp2s, "2d slice")

}

func bench(objs []*posobj, xlen, ylen int, idp idposi.IDPosManI, name string) {
	sttime2 := time.Now()
	sttime := time.Now()
	for _, v := range objs {
		idp.Add(v)
	}
	fmt.Printf("%v add %v\n", name, time.Now().Sub(sttime))

	sttime = time.Now()
	for i := 0; i < 10; i++ {
		for _, v := range objs {
			newpos := [2]int{v.Pos[1], v.Pos[0]}
			idp.UpdateToPos(v, newpos)
			v.Pos = newpos
		}
	}
	fmt.Printf("%v move %v\n", name, time.Now().Sub(sttime))

	sttime = time.Now()
	for _, v := range objs {
		idp.Del(v)
	}
	fmt.Printf("%v del %v\n", name, time.Now().Sub(sttime))

	fmt.Printf("%v %v\n\n", name, time.Now().Sub(sttime2))
}
