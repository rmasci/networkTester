// networkTester project DbGen.go

/*
dbGen document
*/
package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//var layout string = "1/2/2006 03:04:05pm"

func getRand(i int32, seed int64) int32 {
	r := rand.New(rand.NewSource(seed))
	return r.Int31n(i)
}

func getLastName(seed int64) string {
	rand.Seed(seed)
	last := []string{
		"SMITH",
		"JOHNSON",
		"WILLIAMS",
		"JONES",
		"BROWN",
		"DAVIS",
		"MILLER",
		"WILSON",
		"MOORE",
		"TAYLOR",
		"ANDERSON",
		"THOMAS",
		"JACKSON",
		"WHITE",
		"HARRIS",
		"MARTIN",
		"THOMPSON",
		"GARCIA",
		"MARTINEZ",
		"ROBINSON",
		"CLARK",
		"RODRIGUEZ",
		"LEWIS",
		"LEE",
		"WALKER",
		"HALL",
		"ALLEN",
		"YOUNG",
		"HERNANDEZ",
		"KING",
		"WRIGHT",
		"LOPEZ",
		"HILL",
		"SCOTT",
		"GREEN",
		"ADAMS",
		"BAKER",
		"GONZALEZ",
		"NELSON",
		"CARTER",
		"MITCHELL",
		"PEREZ",
		"ROBERTS",
		"TURNER",
		"PHILLIPS",
		"CAMPBELL",
		"PARKER",
		"EVANS",
		"EDWARDS",
		"COLLINS",
	}
	return fmt.Sprintf(last[rand.Intn(len(last))])
}

func getFirstName(seed int64) string {
	rand.Seed(seed)
	first := []string{
		"JAMES",
		"JOHN",
		"ROBERT",
		"MICHAEL",
		"WILLIAM",
		"DAVID",
		"RICHARD",
		"CHARLES",
		"JOSEPH",
		"THOMAS",
		"CHRISTOPHER",
		"DANIEL",
		"PAUL",
		"MARK",
		"DONALD",
		"GEORGE",
		"KENNETH",
		"STEVEN",
		"EDWARD",
		"BRIAN",
		"RONALD",
		"ANTHONY",
		"KEVIN",
		"JASON",
		"MATTHEW",
		"GARY",
		"TIMOTHY",
		"JOSE",
		"LARRY",
		"JEFFREY",
		"FRANK",
		"SCOTT",
		"ERIC",
		"STEPHEN",
		"ANDREW",
		"RAYMOND",
		"GREGORY",
		"JOSHUA",
		"JERRY",
		"DENNIS",
		"WALTER",
		"PATRICK",
		"PETER",
		"HAROLD",
		"DOUGLAS",
		"HENRY",
		"CARL",
		"ARTHUR",
		"RYAN",
		"ROGER",
	}
	return fmt.Sprintf(first[rand.Intn(len(first))])
}

func DbGen(w http.ResponseWriter, req *http.Request) {
	date := time.Now()
	var outStr string
	host, _ := os.Hostname()
	outStr += fmt.Sprintln("<html>")
	outStr += fmt.Sprintln("<div class=\"ttester\">")
	outStr += fmt.Sprintln("<style type=\"text/css\">")
	outStr += fmt.Sprintln("#ttester {")
	outStr += fmt.Sprintln("    font-family: \"Trebuchet MS\", Arial, Helvetica, sans-serif;")
	outStr += fmt.Sprintln("    width: 100%;")
	outStr += fmt.Sprintln("    border-collapse: collapse;")
	outStr += fmt.Sprintln("}")
	outStr += fmt.Sprintln("#ttester td, #ttester th {")
	outStr += fmt.Sprintln("    font-size: 1em;")
	outStr += fmt.Sprintln("    border: 1px solid #000000;")
	outStr += fmt.Sprintln("    padding: 3px 7px 2px 7px;")
	outStr += fmt.Sprintln("}")
	outStr += fmt.Sprintln("#ttester th {")
	outStr += fmt.Sprintln("    text-align: left;")
	outStr += fmt.Sprintln("    padding-top: 5px;")
	outStr += fmt.Sprintln("    padding-bottom: 4px;")
	outStr += fmt.Sprintln("    background-color: #C0C0C0;")
	outStr += fmt.Sprintln("    color: #000000;")
	outStr += fmt.Sprintln("}")
	outStr += fmt.Sprintln("#ttester tr.alt td {")
	outStr += fmt.Sprintln("    color: #000000;")
	outStr += fmt.Sprintln("    background-color: #EAF2D3;")
	outStr += fmt.Sprintln("}")
	outStr += fmt.Sprintln("</style>")
	outStr += fmt.Sprintln("</div>")
	outStr += fmt.Sprintf("<h1>Data Testing %v</h1><hr><pre>%v, HTTP 1.1/ 200 OK\n", host, date.Local().Format(layout))
	outStr += fmt.Sprintf("<table id=\"ttester\"><tr><th>First Name</th><th>Last Name</th><th>Sold Today</th></tr>\n")
	for I := 1; I <= 50; I++ {
		seed := time.Now().UnixNano() + int64(I)
		i := getRand(100, seed)
		d := getRand(99, seed)
		fname := getFirstName(seed)
		lname := getLastName(seed)
		outStr += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%v.%v</td></tr>\n", fname, lname, i, d)
	}
	outStr += fmt.Sprintf("<table><body><html>")
	io.WriteString(w, outStr)
}
