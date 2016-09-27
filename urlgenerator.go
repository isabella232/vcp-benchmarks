package main

import (
	"flag"
	"fmt"
	"math/rand"
)

const (
	factor            = 1000
	default_maxage    = 600
	default_hitpct    = 90
	default_rndseed   = 4242
	default_minsize   = 1000
	default_maxsize   = 200000
	default_hitdelay  = 100
	default_missdelay = 100
	default_numurls   = 10000
)

func main() {
	var maxage int
	var hitpct int
	var misspct = 100 - hitpct
	var rndseed int64
	var minsize int
	var maxsize int
	var numurls int
	var hitdelay int
	var missdelay int
	var urlprefix string

	flag.IntVar(&maxage, "max-age", default_maxage, "Max-age for cacheable urls.")
	flag.IntVar(&hitpct, "hit-percent", default_hitpct, "Percentage of urls that are hits.")
	flag.IntVar(&minsize, "min-size", default_minsize, "Minimum request size.")
	flag.IntVar(&maxsize, "max-size", default_maxsize, "Maximum request size.")
	flag.IntVar(&hitdelay, "hit-delay", default_hitdelay, "Delay for backend-requests that are hits.")
	flag.IntVar(&missdelay, "miss-delay", default_missdelay, "Delay for backend-requests that are misses.")
	flag.IntVar(&numurls, "urls", default_numurls, "Number of urls to generate.")
	flag.StringVar(&urlprefix, "url-prefix", "", "Prefix for url string.")
	flag.Int64Var(&rndseed, "rndseed", default_rndseed, "Seed for random generator.")

	flag.Parse()

	hiturl := fmt.Sprintf("/hit/?content-length&max-age=%d&header-delay=%dpredictable-content=", maxage, hitdelay)
	missurl := fmt.Sprintf("/miss/?content-length&max-age=%d&header-delay=%dpredictable-content=", 0, missdelay)

	urls := make([]string, numurls)
	size := int64(minsize * factor)
	delta := int64(((maxsize - minsize) * factor) / ((numurls * hitpct / 100) + 1))

	for i := 0; i < (numurls * hitpct / 100); i++ {
		urls[i] = fmt.Sprintf("%s/%d%s%d", urlprefix, (size / factor), hiturl, (size / factor))
		size = size + delta
	}

	size = int64(minsize * factor)
	delta = int64(((maxsize - minsize) * factor) / ((numurls * misspct / 100) + 1))

	for i := (numurls * hitpct / 100); i < numurls; i++ {
		urls[i] = fmt.Sprintf("%s/%d%s%d", urlprefix, (size / factor), missurl, (size / factor))
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
