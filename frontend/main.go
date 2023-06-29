// main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	interval := 500
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := request("http://outlier-backend:8090")
		if err != nil {
			fmt.Fprintf(w, "backend fail: %v", err)
			return
		}
		values := r.URL.Query()
		val := values.Get("interval")
		fmt.Sscanf(val, "%v", &interval)
		if interval == 0 {
			interval = 500
		}
		w.Write(data)
	})
	go func() {
		for {
			data, err := request("http://outlier-backend:8090")
			if err != nil {
				fmt.Printf("backend fail: %v\n", err)
			} else {
				fmt.Printf("%v response: %v\n",
					time.Now().Format("2006-01-02 15:04:05"),
					string(data))
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()

	http.ListenAndServe(":8080", nil)
}

func request(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making the request: %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response body: %s\n", err.Error())
		return nil, err
	}
	return body, nil
}
