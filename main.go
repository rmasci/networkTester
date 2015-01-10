package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var runAs = filepath.Base(os.Args[0])
var exitStat int = 0
var bits float64
var srvAddr, port, logDur, logFileName string
var cl, srv, dae, v bool
var timeOut time.Duration
var layout string = "1/2/2006 03:04:05pm"

//var logFile log.Logger

func LoadTest(w http.ResponseWriter, req *http.Request) {
	tn := time.Now()
	c := 1024
	b := make([]byte, c)
	var count float64
	//url := req.URL.Path
	log.Printf("Connect From: %v, to: %v \n", req.RemoteAddr, req.Host)
	for {
		if (count / 1024.0 / 1024.0) >= bits {
			break
		}
		if close, ok := w.(http.CloseNotifier); ok {
		} else {
			<-close.CloseNotify()
		}
		n, err := io.ReadFull(rand.Reader, b)
		if n != len(b) || err != nil {
			fmt.Println("error:", err)
			return
		}
		str := base64.StdEncoding.EncodeToString(b)
		//io.WriteString(w, string(b))
		_, er := fmt.Fprintf(w, str)
		if er != nil {
			log.Printf("Disconnect\n")
			break
		}
		count = count + float64(len(str))
	}
	dur := time.Since(tn)
	tsTmp := strings.Split(dur.String(), ".")
	hms := tsTmp[0]
	ms := tsTmp[1]
	ts := hms + "." + ms[0:2] + "s"
	sec := dur.Seconds()
	bps := ((count * 8) / 1024 / 1024) / sec
	log.Printf("Sent: %.2fm From: %v to %v in %s %.2fMbps\n", count/1024.0/1024.0, req.Host, req.RemoteAddr, ts, bps)
}

func HelloLT(w http.ResponseWriter, req *http.Request) {
	date := time.Now()
	host, _ := os.Hostname()
	outStr := fmt.Sprintf("<h1>Web Testing</h1><hr><pre>HTTP 1.1/ 200 OK\nDate: %v\nFrom: %v\nTo: %v - %v<br>\n", date.Local().Format(layout), req.RemoteAddr, host, req.Host)
	cStr := fmt.Sprintf("Web Testing\n-----------------\nHTTP 1.1/ 200 OK\nDate: %v\nFrom: %v\nTo: %v - %v\n", date, req.RemoteAddr, host, req.Host)
	io.WriteString(w, outStr)
	log.Printf(cStr)
}
func timeoutDialer(ns time.Time) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		c, err := net.Dial(netw, addr)
		if err != nil {
			return nil, err
		}
		c.SetDeadline(ns)
		return c, nil
	}
}

func httpTest(url string) {
	tn := time.Now()
	timeOut := tn.Add(90 * time.Second)
	var res *http.Response
	c := http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(timeOut),
		},
	}
	res, err := c.Get(url)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		time.Sleep(30 * time.Second)
		return
	}
	log.Printf("Connected to %s on port %s ", srvAddr, port)
	httpBod, err := ioutil.ReadAll(res.Body)
	n := len(httpBod)
	res.Body.Close()
	if err != nil {
		log.Printf("ERROR: %v", err)
		time.Sleep(30 * time.Second)
		return
	}
	dur := time.Since(tn)
	tsTmp := strings.Split(dur.String(), ".")
	hms := tsTmp[0]
	ms := tsTmp[1]
	ts := hms + "." + ms[0:2] + "s"
	sec := dur.Seconds()
	bps := (float64(n*8) / 1024 / 1024) / sec
	log.Printf("%s, downloaded %.2fm in %s, %.2fMbps\n", res.Status, float64(n/1024.0/1024.0), ts, bps)
	return
}

func Usage() {
	fmt.Printf("Usage:\n------\nClient: %v -c [-p 8888][-t 1000] [-l 5] -s <Server IP> [-b Megabytes to send]\n", runAs)
	fmt.Printf("Server: %v -s [-p 8888] [-l 5] [-t 1]\n", runAs)
	flag.PrintDefaults()
	fmt.Printf("Disclaimer:\n%v has been written by Richard Masci (rx7322), however use at your own risk", runAs)
	fmt.Printf(", writer is not to be held liable for any damage to system or network.\n")
	fmt.Printf("Shake well. Batteries not included. Void where prohibited. Use only as directed. Do not use this program while operating heavy equipment.\n")
	//fmt.Printf(" No MSG, but may contain peanuts. If conditions persist consult medical advice. ")
	//fmt.Printf(". No animals were harmed in the creation of this program.\nAnything you mail, can and will be used against you.\n")
}

func init() {
	var t string
	// define flags passed at runtime, and assign them to the variables defined above
	flag.BoolVar(&cl, "c", false, "Run as a Client.")
	flag.BoolVar(&srv, "s", false, "Run as a Server.")
	flag.Float64Var(&bits, "b", 100, "Megabytes the server should send to the client.")
	flag.StringVar(&port, "p", "8888", "TCP Port. Defaults to 8888")
	flag.StringVar(&logDur, "l", "5", "Log every xx minute(s)")
	flag.StringVar(&logFileName, "L", "", "Log File")
	flag.StringVar(&srvAddr, "a", "127.0.0.1", "Server IP Address")
	flag.StringVar(&t, "t", "60", "Time out in Seconds")
	flag.BoolVar(&v, "v", false, "Verbose Output -- don't log to file")
	flag.Parse()
	timeOut, _ = time.ParseDuration(t + "s")
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if (cl && srv) || (!cl && !srv) {
		Usage()
		os.Exit(1)
	}
	if runtime.GOOS == "windows" || logFileName == "" {
		v = true
		log.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Could not open logfile ", logFileName, err)
		}
		log.SetOutput(file)
		fmt.Printf("Log file is in %s\n", logFileName)
		defer file.Close()
	}
	if cl {
		timeOut = timeOut * time.Millisecond
		url := "http://" + srvAddr + ":" + port + "/loadtest"
		//var totalT time.Duration
		for {
			httpTest(url)
		}
	} else {
		if bits >= 250 {
			fmt.Printf("Server output too large, scaling down to 250m\n")
			bits = 250
		}
		http.HandleFunc("/loadtest", LoadTest)
		http.HandleFunc("/", HelloLT)
		http.HandleFunc("/dbgen", DbGen)
		s := &http.Server{
			Addr:           ":" + port,
			Handler:        nil,
			ReadTimeout:    timeOut * time.Second,
			WriteTimeout:   timeOut * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Fatal(s.ListenAndServe())
		<-make(chan bool)
	}
}
