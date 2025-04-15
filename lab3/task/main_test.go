package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

// запускаем перед основными функциями по разу чтобы файл остался в памяти в файловом кеше
// ioutil.Discard - это ioutil.Writer который никуда не пишет
func init() {
	SlowSearch(ioutil.Discard)
	FastSearch(ioutil.Discard)
}

// -----
// go test -v

func TestSearch(t *testing.T) {
	slowOut := new(bytes.Buffer)
	SlowSearch(slowOut)
	slowResult := slowOut.String()

	fastOut := new(bytes.Buffer)
	FastSearch(fastOut)
	fastResult := fastOut.String()

	if slowResult != fastResult {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", fastResult, slowResult)
	}
}

// -----
// go test -bench=. -benchmem

func BenchmarkSlow(b *testing.B) {
	cpuFile, err := os.Create("cpu_slow.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuFile.Close()

	memFile, err := os.Create("mem_slow.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer memFile.Close()

	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	for i := 0; i < b.N; i++ {
		SlowSearch(io.Discard)
	}

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func BenchmarkFast(b *testing.B) {
	cpuFile, err := os.Create("cpu_fast.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuFile.Close()

	memFile, err := os.Create("mem_fast.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer memFile.Close()

	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	for i := 0; i < b.N; i++ {
		FastSearch(io.Discard)
	}

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func BenchmarkSlowExecutionTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		start := time.Now()
		SlowSearch(io.Discard)
		b.ReportMetric(time.Since(start).Seconds(), "seconds")
	}
}

func BenchmarkFastExecutionTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		start := time.Now()
		FastSearch(io.Discard)
		b.ReportMetric(time.Since(start).Seconds(), "seconds")
	}
}
