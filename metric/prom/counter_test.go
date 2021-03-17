package prom_test

import (
	"github.com/magiconair/properties/assert"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xmchz/go-one/metric/prom"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
)

func TestCounter(t *testing.T) {
	s := httptest.NewServer(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	defer s.Close()

	scrape := func() string {
		resp, _ := http.Get(s.URL)
		buf, _ := ioutil.ReadAll(resp.Body)
		return string(buf)
	}

	namespace, subsystem, name := "ns", "ss", "foo"
	re := regexp.MustCompile(namespace + `_` + subsystem + `_` + name + `{alpha="alpha-value",beta="beta-value"} ([0-9\.]+)`)

	counter := prom.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      "This is the help string.",
	}, []string{"alpha", "beta"}).With("beta", "beta-value", "alpha", "alpha-value") // order shouldn't matter


	want := 1.11
	counter.Add(want)

	got := func() float64 {
		matches := re.FindStringSubmatch(scrape())
		f, _ := strconv.ParseFloat(matches[1], 64)
		return f
	}
	assert.Equal(t, got(), want)
}


