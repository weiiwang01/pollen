/*

pollen: Entropy-as-a-Server web server

  Copyright (C) 2012-2013 Dustin Kirkland <dustin.kirkland@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, version 3 of the License.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
	"crypto/sha512"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	httpPort    = flag.String("http-port", "80", "The HTTP port on which to listen")
	httpsPort   = flag.String("https-port", "443", "The HTTPS port on which to listen")
	metricsPort = flag.String("metrics-port", "", "The Prometheus metrics HTTP endpoint port")
	device      = flag.String("device", "/dev/random", "The device to use for reading and writing random data")
	size        = flag.Int("bytes", 64, "The size in bytes to read from the random device")
	cert        = flag.String("cert", "/etc/pollen/cert.pem", "The full path to cert.pem")
	key         = flag.String("key", "/etc/pollen/key.pem", "The full path to key.pem")
)

// this matches the syslog.Writer functions
type logger interface {
	Close() error
	Info(string) error
	Err(string) error
	Crit(string) error
	Emerg(string) error
}

type PollenServer struct {
	// randomSource is usually /dev/random or /dev/urandom
	randomSource io.ReadWriter
	log          logger
	readSize     int
	tracker      *Tracker
}

const usePollinateError = "Please use the pollinate client.  'sudo apt-get install pollinate' or download from: https://bazaar.launchpad.net/~pollinate/pollinate/trunk/view/head:/pollinate"

func (p *PollenServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	p.tracker.RequestReceived()
	var avail []byte
	challenge := r.FormValue("challenge")
	if challenge == "" {
		http.Error(w, usePollinateError, http.StatusBadRequest)
		p.tracker.ResponseSent(http.StatusBadRequest, time.Since(startTime))
		return
	}
	checksum := sha512.New()
	io.WriteString(checksum, challenge)
	challengeResponse := checksum.Sum(nil)
	var err error
	_, err = p.randomSource.Write(challengeResponse)
	if err != nil {
		/* Non-fatal error, but let's log this to syslog */
		p.log.Err(fmt.Sprintf("Cannot write to random device at [%v]", time.Now().UnixNano()))
	}
	/* Record entropy bits before */
	avail, err = ioutil.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		/* Non-fatal error */
		p.log.Err(fmt.Sprintf("Cannot record entropy bits at [%v]", time.Now().UnixNano()))
		avail = []byte{'?'}
	}
	p.log.Info(fmt.Sprintf("Server received challenge from [%s, %s] at [%v] with [e%s] available", r.RemoteAddr, r.UserAgent(), time.Now().UnixNano(), strings.Split(string(avail), "\n")[0]))
	data := make([]byte, p.readSize)
	_, err = io.ReadFull(p.randomSource, data)
	if err != nil {
		/* Fatal error for this connection, if we can't read from device */
		p.log.Err(fmt.Sprintf("Cannot read from random device at [%v]", time.Now().UnixNano()))
		http.Error(w, "Failed to read from random device", http.StatusInternalServerError)
		p.tracker.ResponseSent(http.StatusInternalServerError, time.Since(startTime))
		return
	}
	p.tracker.EntropyQa(data)
	checksum.Write(data)
	/* The checksum of the bytes from /dev/random is simply for print-ability, when debugging */
	seed := checksum.Sum(nil)
	fmt.Fprintf(w, "%x\n%x\n", challengeResponse, seed)
	p.tracker.ResponseSent(200, time.Since(startTime))
	/* Record entropy bits after */
	avail, err = ioutil.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		/* Non-fatal error */
		p.log.Err(fmt.Sprintf("Cannot record entropy bits at [%v]", time.Now().UnixNano()))
		avail = []byte{'?'}
	} else {
		p.tracker.SystemEntropy(avail)
	}
	p.log.Info(fmt.Sprintf("Server sent response to [%s, %s] at [%v] in [%.6fs] with [e%s] available",
		r.RemoteAddr, r.UserAgent(), time.Now().UnixNano(), time.Since(startTime).Seconds(), strings.Split(string(avail), "\n")[0]))
}

func main() {
	flag.Parse()
	if *httpPort == "" && *httpsPort == "" {
		fatal("Nothing to do if http and https are both disabled")
	}
	log, err := syslog.New(syslog.LOG_ERR, "pollen")
	if err != nil {
		fatalf("Cannot open syslog: %s\n", err)
	}
	defer log.Close()
	log.Info(fmt.Sprintf("pollen starting at [%v]", time.Now().UnixNano()))
	dev, err := os.OpenFile(*device, os.O_RDWR, 0)
	if err != nil {
		fatalf("Cannot open device: %s\n", err)
	}
	defer dev.Close()
	var tracker *Tracker

	if *metricsPort != "" {
		tracker = NewTracker()
	}
	handler := &PollenServer{randomSource: dev, log: log, readSize: *size, tracker: tracker}
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	var httpListeners sync.WaitGroup
	if *httpPort != "" {
		httpAddr := fmt.Sprintf(":%s", *httpPort)
		httpListeners.Add(1)
		go func() {
			handler.fatal(http.ListenAndServe(httpAddr, mux))
			httpListeners.Done()
		}()
	}
	if *httpsPort != "" {
		httpsAddr := fmt.Sprintf(":%s", *httpsPort)
		httpListeners.Add(1)
		go func() {
			config := &tls.Config{MinVersion: tls.VersionTLS10}
			server := &http.Server{Addr: httpsAddr, Handler: handler, TLSConfig: config}
			handler.fatal(server.ListenAndServeTLS(*cert, *key))
			httpListeners.Done()
		}()
	}
	if *metricsPort != "" {
		metricsAddr := fmt.Sprintf(":%s", *metricsPort)
		httpListeners.Add(1)
		go func() {
			handler.fatal(tracker.StartMetricsServer(metricsAddr))
			httpListeners.Done()
		}()
	}
	httpListeners.Wait()
}

func (p *PollenServer) fatal(args ...interface{}) {
	p.log.Crit(fmt.Sprint(args...))
	fatal(args...)
}

func (p *PollenServer) fatalf(format string, args ...interface{}) {
	p.log.Emerg(fmt.Sprintf(format, args...))
	fatalf(format, args...)
}

func fatal(args ...interface{}) {
	args = append(args, "\n")
	fmt.Fprint(os.Stderr, args...)
	os.Exit(1)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
