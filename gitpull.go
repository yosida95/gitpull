package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/smtp"
	"os/exec"
	"strings"
	"time"

	"github.com/yosida95/recvknocking"
)

func report(out string) {
	conn, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Quit()

	if err := conn.Mail("ssh-keys"); err != nil {
		log.Println(err)
		return
	}

	if err := conn.Rcpt("root"); err != nil {
		log.Println(err)
		return
	}

	if data, err := conn.Data(); err != nil {
		log.Println(err)
		return
	} else {
		io.WriteString(data, strings.Join([]string{
			"From: ssh-keys",
			"To: root",
			"Subject: Report of updating ssh-keys",
			"",
			out,
		}, "\r\n"))
		data.Close()
	}
}

func reportError(out string, err error) {
	body := "" +
		"An error occurred.  " +
		fmt.Sprintf("Details of this error is \"%s\"\n", err.Error()) +
		"The following lines is the output.\n" +
		out

	report(body)
}

func current() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	return string(out), err
}

func update(net.IP) {
	var result string

	prev, err := current()
	if err != nil {
		reportError(prev, err)
		return
	}

	out, err := exec.Command("git", "pull", repository, refspec).CombinedOutput()
	if err != nil {
		reportError(string(out), err)
		return
	}
	result += string(out)

	cur, err := current()
	if err != nil {
		reportError(cur, err)
		return
	}

	if prev != cur {
		out, err = exec.Command("make", "all").CombinedOutput()
		if err != nil {
			reportError(string(out), err)
		}
		result += string(out)
	}

	report(result)
	return
}

var (
	socket     string
	repository string
	refspec    string
)

func init() {
	flag.StringVar(&socket, "socket", "", "socket to listen")
	flag.StringVar(&repository, "repository", "origin", "remote repository to pull")
	flag.StringVar(&refspec, "refspec", "master", "refspec to pull")
}

func main() {
	flag.Parse()

	config := recvknocking.Config{
		Count:    3,
		Duration: 1 * time.Second,
		Factory: func() (l net.Listener, err error) {
			l, err = net.Listen("tcp4", socket)
			if err != nil {
				log.Fatalln(err)
				return
			}

			log.Printf("Listen on %s", l.Addr().String())
			return
		},
		Handler: update,
	}

	r := recvknocking.NewReceiver(config)
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
