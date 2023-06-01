package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Tracker struct {
	pollenHttpRequestTotal                       prometheus.Counter
	pollenHttpResponseCode                       *prometheus.CounterVec
	pollenHttpResponseSeconds                    *prometheus.HistogramVec
	pollenSystemEntropy                          prometheus.Gauge
	pollenResponseEntropyPerByte                 prometheus.Histogram
	pollenResponseEntropyArithmeticMeanDeviation prometheus.Histogram
}

// entropyPerByte calculates the entropy per byte for a given byte array.
func (t *Tracker) entropyPerByte(input []byte) float64 {
	ca := [256]int{}
	l := len(input)
	for _, c := range input {
		ca[c] += 1
	}
	pa := [256]float64{}
	for i, c := range ca {
		pa[i] = float64(c) / float64(l)
	}
	var ent float64
	for _, p := range pa {
		if p > 0.0 {
			ent += p * math.Log2(1.0/p)
		}
	}
	return ent
}

// arithmeticMeanDeviation calculates the arithmetic mean deviation from the
// central value of a given byte array.
func (t *Tracker) arithmeticMeanDeviation(input []byte) float64 {
	s := 0
	for _, c := range input {
		s += int(c)
	}
	return math.Abs(127.5 - float64(s)/float64(len(input)))
}

// chiSquare calculates the chi-square value of the input byte array.
func (t *Tracker) chiSquare(input []byte) float64 {
	bins := [256]int{}
	n := float64(len(input))
	for _, c := range input {
		bins[c] += 1
	}
	var cs float64
	for _, o := range bins {
		e := n / 256.0
		cs += math.Pow(e-float64(o), 2) / e
	}
	return cs
}

// RequestReceived increments the counter for the total number of HTTP requests
// received by the application. If the Tracker receiver is nil, the function
// does nothing.
func (t *Tracker) RequestReceived() {
	if t == nil {
		return
	}
	t.pollenHttpRequestTotal.Inc()
}

// ResponseSent increments the counters for HTTP response codes and observes
// the duration in the histogram vector for HTTP response times. If the
// Tracker receiver is nil, the function does nothing.
func (t *Tracker) ResponseSent(code int, duration time.Duration) {
	if t == nil {
		return
	}
	sc := strconv.Itoa(code)
	t.pollenHttpResponseCode.WithLabelValues(sc).Inc()
	t.pollenHttpResponseSeconds.WithLabelValues(sc).Observe(duration.Seconds())
}

// SystemEntropy sets the gauge for system entropy. The input should be the
// content of /proc/sys/kernel/random/entropy_avail. If the Tracker receiver
// is nil or the input is not valid, the function does nothing.
func (t *Tracker) SystemEntropy(entropyAvail []byte) {
	if t == nil {
		return
	}
	ent, err := strconv.ParseFloat(string(entropyAvail), 64)
	if err == nil {
		return
	}
	t.pollenSystemEntropy.Set(ent)
}

// EntropyQa observes the arithmetic mean deviation and entropy per byte of the
// response in the respective histograms. If the Tracker receiver is nil,
// the function does nothing.
func (t *Tracker) EntropyQa(input []byte) {
	if t == nil {
		return
	}
	t.pollenResponseEntropyArithmeticMeanDeviation.Observe(t.arithmeticMeanDeviation(input))
	t.pollenResponseEntropyPerByte.Observe(t.entropyPerByte(input))
}

// StartMetricsServer starts a HTTP server that exposes the metrics in Prometheus format.
func (t *Tracker) StartMetricsServer(address string) error {
	metricMux := http.NewServeMux()
	metricMux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(address, metricMux)
}

// NewTracker creates a new Tracker with the Prometheus metrics initialized.
func NewTracker() *Tracker {
	return &Tracker{
		pollenHttpRequestTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "pollen_http_requests_total",
			Help: "The total number of requests",
		}),
		pollenHttpResponseCode: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "pollen_http_responses_codes",
			Help: "Total responses sent to clients by code",
		}, []string{"code"}),
		pollenHttpResponseSeconds: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "pollen_http_response_seconds",
			Help:    "Response time by code",
			Buckets: []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.1, 1.0},
		}, []string{"code"}),
		pollenSystemEntropy: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "pollen_system_entropy",
			Help: "System available entropy (entropy_avail)",
		}),
		pollenResponseEntropyPerByte: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "pollen_response_entropy_per_byte",
			Help:    "Entropy per bit of the random data in response",
			Buckets: []float64{1.0, 2.0, 3.0, 4.0, 4.5, 5.0, 5.5, 6.0, 6.5, 7.0, 7.5},
		}),
		pollenResponseEntropyArithmeticMeanDeviation: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "pollen_response_entropy_arithmetic_mean_deviation",
			Help:    "Arithmetic mean deviation of the random data in response",
			Buckets: []float64{10.0, 20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0},
		}),
	}
}
