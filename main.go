package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Latters string
	Rows    int
	Columns []int
}

func main() {
	var cp = flag.String("config", "./config.toml", "Set path to config.toml")
	// Decode Toml Config File
	var conf Config
	if _, err := toml.DecodeFile(*cp, &conf); err != nil {
		panic(err)
	}
	// Generate Random Seed
	rand.Seed(time.Now().UnixNano())
	columns_count := conf.Columns
	latters := conf.Latters
	c := make(chan string)
	wg := &sync.WaitGroup{}
	for i := 0; i < conf.Rows; i++ {
		wg.Add(1)
		go GenRow(columns_count, latters, c, wg)
	}
	wgout := &sync.WaitGroup{}
	go PrintRow(c, wgout)
	wg.Wait()
	close(c)
	wgout.Wait()
}

func RandString(n int, latters string) string {
	// Generate Random Strings
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = latters[rand.Intn(len(latters))]
	}
	return string(bytes)
}

func GenRow(columns []int, latters string, c chan string, wg *sync.WaitGroup) {
	// Create 1 row
	defer func() {
		wg.Done()
	}()
	row := ""
	for i, n := range columns {
		r := RandString(n, latters)
		if i != len(columns)-1 {
			row += r + ","
		} else {
			row += r + "\n"
			c <- row
		}
	}
	return
}

func PrintRow(c chan string, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case r, ok := <-c:
			if !ok {
				return
			}
			fmt.Printf("%s", r)
		}
	}
}
