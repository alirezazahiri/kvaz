package main

import (
	"fmt"
	"kvaz/internal/tcp"
)

func main() {
	ts := tcp.NewTCPServer(6379)
	defer ts.Close()

	for {
		tc, err := ts.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		go func(c *tcp.TCPConnection) {
			for {
				data := make(chan string)
				isWritten := make(chan bool)

				go func() {
					_, err := c.Write([]byte(<-data))
					if err != nil {
						fmt.Println("error occured while writing data")
						isWritten <- false
						return
					}
					isWritten <- true
				}()

				go func() {
					d := make([]byte, 1024)
					_, err := c.Read(d)
					if err != nil {
						data <- ""
						return
					}
					fmt.Println("read data in reader goroutine:", string(d))
					data <- string(d)
				}()

				<-isWritten
			}
		}(tc)
	}
}
