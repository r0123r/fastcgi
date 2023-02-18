package fcgiclient

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	BA "github.com/r0123r/balance/backends"
)

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(us net.Conn, backend BA.Backend) {
	if backend == nil {
		log.Printf("no backend available for connection from %s", us.RemoteAddr())
		us.Close()
		return
	}

	ds, err := net.Dial("tcp", backend.String())
	if err != nil {
		log.Printf("failed to dial %s: %s", backend, err)
		us.Close()
		return
	}

	// Ignore errors
	go copy(ds, us)
	go copy(us, ds)
}

func tcpBalance(bind string, backends BA.Backends) error {
	log.Println("using tcp balancing")
	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("failed to bind: %s", err)
	}

	log.Printf("listening on %s, balancing %d backends", bind, backends.Len())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		go handleConnection(conn, backends.Choose())
	}

	return err
}

func php_fcgi(php_exe, adr string) {
	avg := []string{php_exe, "-b", adr}

	for {

		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
		procAttr.Sys = &syscall.SysProcAttr{HideWindow: true}

		php_process, err := os.StartProcess(php_exe, avg, procAttr)

		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}
		php_process.Wait()
		time.Sleep(time.Second * 1)
	}
}
func Start_balancer(port, count int) {
	php_exe, err := exec.LookPath("php-cgi.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	bind := fmt.Sprint("127.0.0.1:", port)
	hosts := make([]string, count)
	for i := 0; i < count; i++ {
		hosts[i] = fmt.Sprint("127.0.0.1:", port+i+1)
		go php_fcgi(php_exe, hosts[i])
	}
	ba := BA.Build("round-robin", hosts)
	go tcpBalance(bind, ba)
}
