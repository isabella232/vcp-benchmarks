package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/Iwark/spreadsheet"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

func createBenchSheet(secret string, sid string, name string) (*spreadsheet.Worksheet, error) {
	data, _ := ioutil.ReadFile(secret)
	conf, _ := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	client := conf.Client(context.TODO())

	service := &spreadsheet.Service{Client: client}
	service.ReturnEmpty = true

	sheets, _ := service.Get(sid)

	_, err := sheets.NewWorksheet(name, 21, 10)
	if err != nil {
		panic(err)
	}

	wsheet, err := sheets.FindByTitle(name)
	if err != nil {
		panic(err)
	}

	return wsheet, err
}

func populateWrkSheet(ws *spreadsheet.Worksheet) {
	ws.Rows[0][0].Update("VHA Status")
	ws.Rows[1][0].Update("Backend Requests")
	ws.Rows[2][0].Update("wrk Threads")
	ws.Rows[3][0].Update("wrk Connections")
	ws.Rows[4][0].Update("Latency")
	ws.Rows[5][0].Update("Latency StdDev")
	ws.Rows[6][0].Update("Latency Max")
	ws.Rows[7][0].Update("Latency +-StdDev")
	ws.Rows[8][0].Update("Req/sec")
	ws.Rows[9][0].Update("Req StdDev")
	ws.Rows[10][0].Update("Req Max")
	ws.Rows[11][0].Update("Req +-StdDev")
	ws.Rows[12][0].Update("Total Requests")
	ws.Rows[13][0].Update("Benchmark Runtime")
	ws.Rows[14][0].Update("Total Transfered")
	ws.Rows[15][0].Update("Connect Errors")
	ws.Rows[16][0].Update("Read Errors")
	ws.Rows[17][0].Update("Write Errors")
	ws.Rows[18][0].Update("Timeout Errors")
	ws.Rows[19][0].Update("Requests/sec")
	ws.Rows[20][0].Update("Transfer/sec")
}

func populateSiegeSheet(ws *spreadsheet.Worksheet) {
	ws.Rows[0][0].Update("VHA Status")
	ws.Rows[1][0].Update("Backend Requests")
	ws.Rows[2][0].Update("Siege Concurrent Users")
	ws.Rows[3][0].Update("Transactions")
	ws.Rows[4][0].Update("Availability")
	ws.Rows[5][0].Update("Elapsed time")
	ws.Rows[6][0].Update("Data transferred")
	ws.Rows[7][0].Update("Response time")
	ws.Rows[8][0].Update("Transaction rate")
	ws.Rows[9][0].Update("Throughput")
	ws.Rows[10][0].Update("Concurrency")
	ws.Rows[11][0].Update("Successful transactions")
	ws.Rows[12][0].Update("Failed transactions")
	ws.Rows[13][0].Update("Longest transaction")
	ws.Rows[14][0].Update("Shortest transaction")
}

func parseWrkBenchmark(ws *spreadsheet.Worksheet, filepath string, offset int, benchrun int, vhastatus string) {
	var line string

	file, err := os.Open(filepath + "/backend_requests" + strconv.Itoa(benchrun) + "_total_vha" + vhastatus + ".log")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		ws.Rows[1][offset].Update(line)
	}
	file.Close()

	file, err = os.Open(filepath + "/benchmark_test" + strconv.Itoa(benchrun) + "_vha" + vhastatus + ".log")
	if err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(file)

	ws.Rows[0][offset].Update(vhastatus)

	wrkp_re := regexp.MustCompile(`^\s*(\S*) threads and (\S*) connections`)
	tlat_re := regexp.MustCompile(`^\s*Latency\s*(\S*s)\s*(\S*s)\s*(\S*s)\s*(\S*%)`)
	treq_re := regexp.MustCompile(`^\s*Req/Sec\s*(\S*)\s*(\S*)\s*(\S*)\s*(\S*)`)
	req_re := regexp.MustCompile(`^\s*(\S*)\srequests in (\S*[s,m,h]), (\S*) read`)
	err_re := regexp.MustCompile(`^\s*Socket errors: connect (\S*), read (\S*), write (\S*), timeout (\S*)`)
	rps_re := regexp.MustCompile(`^\s*Requests/sec:\s*(\S*)`)
	xfer_re := regexp.MustCompile(`^\s*Transfer/sec:\s*(\S*)`)

	for scanner.Scan() {
		line = scanner.Text()
		if wrkp_re.MatchString(line) {
			match := wrkp_re.FindStringSubmatch(line)
			ws.Rows[2][offset].Update(match[1])
			ws.Rows[3][offset].Update(match[2])
		}
		if tlat_re.MatchString(line) {
			match := tlat_re.FindStringSubmatch(line)
			ws.Rows[4][offset].Update(match[1])
			ws.Rows[5][offset].Update(match[2])
			ws.Rows[6][offset].Update(match[3])
			ws.Rows[7][offset].Update(match[4])
		}
		if treq_re.MatchString(line) {
			match := treq_re.FindStringSubmatch(line)
			ws.Rows[8][offset].Update(match[1])
			ws.Rows[9][offset].Update(match[2])
			ws.Rows[10][offset].Update(match[3])
			ws.Rows[11][offset].Update(match[4])
		}
		if req_re.MatchString(line) {
			match := req_re.FindStringSubmatch(line)
			ws.Rows[12][offset].Update(match[1])
			ws.Rows[13][offset].Update(match[2])
			ws.Rows[14][offset].Update(match[3])
		}
		if err_re.MatchString(line) {
			match := err_re.FindStringSubmatch(line)
			ws.Rows[15][offset].Update(match[1])
			ws.Rows[16][offset].Update(match[2])
			ws.Rows[17][offset].Update(match[3])
			ws.Rows[18][offset].Update(match[4])
		}
		if rps_re.MatchString(line) {
			match := rps_re.FindStringSubmatch(line)
			ws.Rows[19][offset].Update(match[1])
		}
		if xfer_re.MatchString(line) {
			match := xfer_re.FindStringSubmatch(line)
			ws.Rows[20][offset].Update(match[1])
		}
	}

	if scanerr := scanner.Err(); scanerr != nil {
		fmt.Println(scanerr)
	}

	// Make sure call Synchronize to reflect the changes
	ws.Synchronize()
	file.Close()
}

func parseSiegeBenchmark(ws *spreadsheet.Worksheet, filepath string, offset int, benchrun int, vhastatus string) {
	var line string

	file, err := os.Open(filepath + "/backend_requests" + strconv.Itoa(benchrun) + "_siege_total_vha" + vhastatus + ".log")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		ws.Rows[1][offset].Update(line)
	}
	file.Close()

	file, err = os.Open(filepath + "/benchmark_test" + strconv.Itoa(benchrun) + "_siege_vha" + vhastatus + ".log")
	if err != nil {
		panic(err)
	}
	scanner = bufio.NewScanner(file)

	ws.Rows[0][offset].Update(vhastatus)

	siegelog_re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2},\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*),\s*(\S*)`)

	for scanner.Scan() {
		line = scanner.Text()
		if siegelog_re.MatchString(line) {
			match := siegelog_re.FindStringSubmatch(line)
			ws.Rows[3][offset].Update(match[1])
			ws.Rows[5][offset].Update(match[2])
			ws.Rows[6][offset].Update(match[3])
			ws.Rows[7][offset].Update(match[4])
			ws.Rows[8][offset].Update(match[5])
			ws.Rows[9][offset].Update(match[6])
			ws.Rows[10][offset].Update(match[7])
			ws.Rows[11][offset].Update(match[8])
			ws.Rows[12][offset].Update(match[9])
		}
	}

	if scanerr := scanner.Err(); scanerr != nil {
		fmt.Println(scanerr)
	}

	// Make sure call Synchronize to reflect the changes
	ws.Synchronize()
	file.Close()
}

func main() {
	var sheetid string
	var clientsecret string
	var benchdir string
	var benchtype string
	var reps int

	flag.StringVar(&sheetid, "sheet-id", "", "The ID of the Google spreadsheet to store results in.")
	flag.StringVar(&clientsecret, "client-secret", "client_secret.json", "The filename of the client secret.")
	flag.StringVar(&benchdir, "benchmark-dir", "", "The directory containing the benchmark data.")
	flag.StringVar(&benchtype, "benchmark-type", "wrk", "The type of benchmark data to store [wrk, siege].")
	flag.IntVar(&reps, "repetitions", 3, "Number of test repetitions in set")

	flag.Parse()
	ws, _ := createBenchSheet(clientsecret, sheetid, benchdir)

	offset := 1

	if benchtype == "wrk" {
		populateWrkSheet(ws)
		for i := 1; i <= reps; i++ {
			parseWrkBenchmark(ws, benchdir, offset, i, "enabled")
			offset++
		}
		for i := 1; i <= reps; i++ {
			parseWrkBenchmark(ws, benchdir, offset, i, "disabled")
			offset++
		}
	}
	if benchtype == "siege" {
		populateSiegeSheet(ws)
		for i := 1; i <= reps; i++ {
			parseSiegeBenchmark(ws, benchdir, offset, i, "enabled")
			offset++
		}
		for i := 1; i <= reps; i++ {
			parseSiegeBenchmark(ws, benchdir, offset, i, "disabled")
			offset++
		}
	}
}
