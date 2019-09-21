package http

import (
	"log"
	"time"
)

type Timings struct {
	DNSLookup,
	TCPConnection,
	TLSHandshake,
	ServerProcessing,
	ContentTransfer,
	NameLookup,
	Connect,
	PreTransfer,
	StartTransfer,
	Total time.Duration
}

func NewTimings() Timings {
	return Timings{}
}

func (t *Timings) WithDNSLookup(d time.Duration) *Timings {
	t.DNSLookup = d
	return t
}

func (t *Timings) WithTCPConnection(d time.Duration) *Timings {
	t.TCPConnection = d
	return t
}

func (t *Timings) WithTLSHandshake(d time.Duration) *Timings {
	t.TLSHandshake = d
	return t
}

func (t *Timings) WithServerProcessing(d time.Duration) *Timings {
	t.ServerProcessing = d
	return t
}

func (t *Timings) WithContentTransfer(d time.Duration) *Timings {
	t.ContentTransfer = d
	return t
}

func (t *Timings) WithNameLookup(d time.Duration) *Timings {
	t.NameLookup = d
	return t
}

func (t *Timings) WithConnect(d time.Duration) *Timings {
	t.Connect = d
	return t
}

func (t *Timings) WithPreTransfer(d time.Duration) *Timings {
	t.PreTransfer = d
	return t
}

func (t *Timings) WithStartTransfer(d time.Duration) *Timings {
	t.StartTransfer = d
	return t
}

func (t *Timings) WithTotal(d time.Duration) *Timings {
	t.Total = d
	return t
}

func (t *Timings) Print() {
	log.Printf("dns %d", t.DNSLookup/1000000)
	log.Printf("tcp %d", t.TCPConnection/1000000)
	log.Printf("tls %d", t.TLSHandshake/1000000)
	log.Printf("server %d", t.ServerProcessing/1000000)
	log.Printf("content %d", t.ContentTransfer/1000000)
	log.Printf("name %d", t.NameLookup/1000000)
	log.Printf("connect %d", t.Connect/1000000)
	log.Printf("pre %d", t.PreTransfer/1000000)
	log.Printf("start %d", t.StartTransfer/1000000)
	log.Printf("total %d", t.Total/1000000)
}
