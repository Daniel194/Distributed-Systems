package main

import (
	"net/http"; "io"
	"time"
)

func hello(res http.ResponseWriter, req *http.Request) {
	t := time.Now()

	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		`<DOCTYPE html>
	    <html>
	      <head>
		  <title>Time</title>
	      </head>
	      <body>
		  <h1>Hello !</h1>
		  <br />
		  <h3>The Current Time Is :
		  ` + t.Format("2006-01-02 15:04:05") + `</h3>
	      </body>
	    </html>`,
	)
}
func main() {
	http.HandleFunc("/time", hello)
	http.ListenAndServe(":9000", nil)
}