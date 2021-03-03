package app

import (
	"bufio"
	"fmt"
	"github.com/urfin-djuice/ratelimit/pkg/params"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Application application struct
type Application struct {
	Rate        int
	Inflight    int
	Command     string
	stopSig     chan bool
	inputStream chan string
}

// NewApplication create new application
func NewApplication() (*Application, error) {
	var a Application
	var err error
	a.Rate, a.Inflight, a.Command, err = params.Get(os.Args[1:])
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Run run application
func (a *Application) Run() chan bool {
	a.stopSig = make(chan bool)
	go a.reader()
	go a.run()
	return a.stopSig
}

func (a *Application) reader() {
	scanner := bufio.NewScanner(os.Stdin)
	a.inputStream = make(chan string, a.Inflight)
	for scanner.Scan() {
		a.inputStream <- scanner.Text()
	}
	close(a.inputStream)
}

func (a *Application) run() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGQUIT)
	var wg sync.WaitGroup
	wg.Add(a.Inflight)
	for i := 1; i <= a.Inflight; i++ {
		go func() {
			defer wg.Done()
			t := time.NewTimer(0)
			for {
				select {
				case <-sigc:
					return
				case <-t.C:
					for i := 1; i <= a.Rate; i++ {
						s, isOpen := <- a.inputStream
						if !isOpen {
							return
						}
						go a.process(s)
					}
					t = time.NewTimer(time.Second)
				}
			}
		}()
	}
	wg.Wait()
	close(a.stopSig)
}

func (a *Application) process(param string) {
	c := strings.Split(fmt.Sprintf(a.Command, param), " ")
	cmd := exec.Command(c[0], c[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(string(stdout))
}
