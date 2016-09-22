package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const (
	rndseed = 4242
	numurls = 100000
	hitpct  = 90
	misspct = 10
	minsize = 1000
	maxsize = 600000
	factor  = 1000
)

func main() {
	hiturl := string("/hit/?content-length&max-age=6000&header-delay=100predictable-content=")
	missurl := string("/miss/?content-length&max-age=0&header-delay=100predictable-content=")

	urls := make([]string, numurls)
	size := int64(minsize * factor)
	delta := int64(((maxsize - minsize) * factor) / ((numurls * hitpct / 100) + 1))

	for i := 0; i < (numurls * hitpct / 100); i++ {
		urls[i] = "/" + strconv.FormatInt((size/factor), 10) + hiturl + strconv.FormatInt((size/factor), 10)
		size = size + delta
	}

	size = int64(minsize * factor)
	delta = int64(((maxsize - minsize) * factor) / ((numurls * misspct / 100) + 1))

	for i := (numurls * hitpct / 100); i < numurls; i++ {
		urls[i] = "/" + strconv.FormatInt((size/factor), 10) + missurl + strconv.FormatInt((size/factor), 10)
		size = size + delta
	}

	rand.Seed(rndseed)
	for i := len(urls) - 1; i > 0; i-- {
		j := rand.Intn(i)
		urls[i], urls[j] = urls[j], urls[i]
	}

	for i := 0; i < len(urls); i++ {
		fmt.Println(urls[i])
	}
}
