package main

import (
	"fmt"
	"store/store"
	"time"
)

// MetricsMiddleware tracks usage and statistics of the store
// Like LoggingMiddleware, it wraps a Storer and records how many times an operation is called

type MetricsMiddlewareCallCounts struct {
	get    int
	set    int
	delete int
	len    int

	getmisses int
	setmisses int
}
type MetricsMiddlewareLatency struct {
	get time.Duration
	set time.Duration
}
type MetricsMiddleware struct {
	Inner          store.Storer
	CallCounts     MetricsMiddlewareCallCounts
	TotalLatencies MetricsMiddlewareLatency
}

func NewMetricsMiddleware(inner store.Storer) *MetricsMiddleware {
	return &MetricsMiddleware{Inner: inner}
}

func (m *MetricsMiddleware) Get(key string) (string, error) {
	start := time.Now()
	val, err := m.Inner.Get(key)
	if err != nil {
		m.CallCounts.getmisses++
	}
	m.CallCounts.get++
	m.TotalLatencies.get += time.Since(start)
	return val, err
}

func (m *MetricsMiddleware) Set(key, value string) error {
	start := time.Now()
	err := m.Inner.Set(key, value)
	if err != nil {
		m.CallCounts.setmisses++
	}
	m.CallCounts.set++
	m.TotalLatencies.set += time.Since(start)
	return err
}

func (m *MetricsMiddleware) Delete(key string) {
	m.Inner.Delete(key)
	m.CallCounts.delete++
}

func (m *MetricsMiddleware) Len() int {
	got := m.Inner.Len()
	m.CallCounts.len++
	return got
}

func (m *MetricsMiddleware) Keys() []string {
	return m.Inner.Keys()
}

func (m *MetricsMiddleware) Report() {
	fmt.Printf("GET calls: %d (missed %d)\n", m.CallCounts.get, m.CallCounts.getmisses)
	fmt.Printf("SET calls: %d (missed %d)\n", m.CallCounts.get, m.CallCounts.setmisses)
	fmt.Printf("DELETE calls: %d\n", m.CallCounts.delete)
	fmt.Printf("LEN calls: %d\n", m.CallCounts.len)

	if m.CallCounts.get > 0 {
		avg := m.TotalLatencies.get / time.Duration(m.CallCounts.get)
		fmt.Printf("Average GET latency: %v\n", avg)
	}
	if m.CallCounts.set > 0 {
		avg := m.TotalLatencies.set / time.Duration(m.CallCounts.set)
		fmt.Printf("Average SET latency: %v\n", avg)
	}
}

/*
func CreateStore(maxSize int) store.Storer {
	plain := kv.NewStore(maxSize)
	logger := NewLoggingMiddleware(plain)
	metrics := NewMetricsMiddleware(logger)
	return metrics
}

this is how you can make your CreateStore now

If we did metrics = New..(plain) and then logger = New..(metrics) and returned logger
then we would be evaluating how long the store itself takes
But if we keep metrics outside (like in the shown func) then we are also taking time for logging included
*/
