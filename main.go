package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const GOOGLE = "https://bbc.com"

var err error

func main() {
	app := &cli.App{
		Name:  "kk",
		Usage: "knock the knock",
		Action: func(c *cli.Context) error {
			input := c.Args().Get(0)
			if input == "" {
				log.Fatal("no input")
			}
			var host string
			var port = 443
			if strings.Contains(input, ":") {
				host = strings.Split(input, ":")[0]
				port, err = strconv.Atoi(strings.Split(input, ":")[1])
				if err != nil {
					return err
				}
			} else {
				host = input
			}
			log.Println("Testing " + host + ":" + strconv.Itoa(port))

			dialer := &net.Dialer{
				Timeout: 5 * time.Second,
			}

			client := http.Client{
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = host + ":" + strconv.Itoa(port)
						return dialer.DialContext(ctx, network, addr)
					},
				},
			}
			_, err = client.Get(GOOGLE)
			if err != nil {
				netErr, ok := err.(net.Error)
				if !ok {
					log.Println("Unexpected error.\n" + err.Error())
				}
				if netErr.Timeout() {
					log.Println("Seems blocked.")
				} else {
					urlErr, ok := err.(*url.Error)
					if !ok {
						return err
					}
					opErr, ok := urlErr.Err.(*net.OpError)
					if !ok {
						return err
					}
					switch t := opErr.Err.(type) {
					case *os.SyscallError:
						if errno, ok := t.Err.(syscall.Errno); ok {
							switch errno {
							case syscall.WSAECONNRESET:
								fallthrough
							case syscall.ECONNRESET:
								log.Println("Seems OK.")
								return nil
							}
						}
					}
					return err
				}
			} else {
				log.Println("Seems you have ignored wall, not working.")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
