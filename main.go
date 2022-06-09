package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const GOOGLE = "https://bbc.com"

var err error

func http2errno(v error) uintptr {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Uintptr {
		return uintptr(rv.Uint())
	}
	return 0
}

func isRSTError(errno syscall.Errno) bool {
	const WSAECONNRESET = 10054
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "windows" {
		if http2errno(errno) == WSAECONNRESET {
			return true
		}
		return false
	} else {
		if errno == syscall.ECONNRESET {
			return true
		}
		return false
	}
}

func main() {
	app := &cli.App{
		Name:  "kk",
		Usage: "knock the knock",
		Action: func(c *cli.Context) error {
			input := c.Args().Get(0)
			if input == "" {
				fmt.Println("no input")
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
			fmt.Println("Testing " + host + ":" + strconv.Itoa(port))

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
					fmt.Println("Unexpected error.\n" + err.Error())
				}
				if netErr.Timeout() {
					fmt.Println("Seems blocked.")
				} else {
					urlErr, ok := err.(*url.Error)
					if !ok {
						return err
					}
					opErr, ok := urlErr.Err.(*net.OpError)
					if !ok {
						return err
					}
					syscallErr, ok := opErr.Err.(*os.SyscallError)
					if !ok {
						return err
					}
					syscallErrNo, ok := syscallErr.Err.(syscall.Errno)
					if !ok {
						return err
					}
					if isRSTError(syscallErrNo) {
						fmt.Println("Seems OK.")
						return nil
					}
					return err
				}
			} else {
				fmt.Println("Seems you have ignored wall, not working.")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
