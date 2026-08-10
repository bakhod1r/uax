package main

import (
	"testing"

	"github.com/avct/uasurfer"
	"github.com/bakhod1r/uax"
	"github.com/mileusna/useragent"
	"github.com/mssola/user_agent"
	"github.com/ua-parser/uap-go/uaparser"
)

const ua = "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Mobile Safari/537.36"

var uapParser = uaparser.NewFromSaved()

func BenchmarkUax(b *testing.B) {
	var c uax.Client
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		uax.ParseTo(&c, ua)
	}
}

func BenchmarkMileusna(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = useragent.Parse(ua)
	}
}

func BenchmarkMssola(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = user_agent.New(ua)
	}
}

func BenchmarkUasurfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		uasurfer.Parse(ua)
	}
}

func BenchmarkUapGo(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uapParser.Parse(ua)
	}
}
