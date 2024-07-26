package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type Scanner struct {
	Ports   []uint
	Threads uint
}

type Status struct {
	status bool
	host   string
	port   uint
}

func NewScanner(threads uint, ports ...uint) *Scanner {
	var scanner Scanner

	scanner.Ports = ports
	scanner.Threads = threads

	return &scanner
}

func (s *Scanner) Scan(targets ...string) {
	var wg sync.WaitGroup

	c := make(chan Status)
	semaphore := make(chan struct{}, s.Threads)

	for _, target := range targets {
		for _, port := range s.Ports {
			wg.Add(1)

			go func(target string, port uint, c chan Status, semaphore chan struct{}) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				s.checkPort(target, port, c)
			}(target, port, c, semaphore)
		}
	}

	result := make([]Status, len(targets)*len(s.Ports))

	for i := 0; i < len(result); i++ {
		result[i] = <-c

		if result[i].status {
			fmt.Printf("[+] %s %d\n", result[i].host, result[i].port)
		}
	}

	close(c)
	wg.Wait()
}

func (s *Scanner) checkPort(target string, port uint, c chan Status) {
	_port := strconv.Itoa(int(port))

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(target, _port), time.Second)
	if err != nil {
		c <- Status{status: false, host: target, port: port}
		return
	}

	if conn != nil {
		conn.Close()
	}

	c <- Status{status: true, host: target, port: port}
}
