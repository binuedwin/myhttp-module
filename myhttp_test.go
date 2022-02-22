package main

import (
	"sync"
	"testing"
)

func TestSheduleWorkers(t *testing.T) {
	reqChan := make(chan string, 100)
	for i := 0; i < 100; i++ {
		reqChan <- ""
	}

	close(reqChan)

	respChan := make(chan string, 100)

	var wg sync.WaitGroup

	// lets schedule 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(&wg, func(string) []byte { return []byte("") }, reqChan, respChan)
	}

	wg.Wait()

	close(respChan)

	if len(respChan) != 100 {
		t.Errorf("Error in sheduling workers: expected 100 responses but got %d", len(respChan))
	}
}

func TestGenerateMD5HashString(t *testing.T) {
	suite := []struct {
		got  string
		want string
	}{
		{"Hello, World!", "65a8e27d8879283831b664bd8b7f0ad4"},
		{"hi mom", "92959a96fd69146c5fe7cbde6e5720f2"},
		{"Go doesn't ship your tests", "294002c9bed37bc15d4c0cf5431e3c60"},
	}

	for _, entry := range suite {
		result := generateMD5HashString([]byte(entry.got))
		if entry.want != result {
			t.Errorf("MD5 Hash of '%s' was incorrect, got: '%s', want: '%s'", entry.got, result, entry.want)
		}
	}
}

func TestNormalizeToURL(t *testing.T) {
	suite := []struct {
		got  string
		want string
	}{
		{"adjust.com", "http://adjust.com"},
		{"google.com", "http://google.com"},
		{"https://facebook.com", "https://facebook.com"},
		{"http://yahoo.com", "http://yahoo.com"},
		{"yandex.com", "http://yandex.com"},
	}

	for _, entry := range suite {
		result := normalizeToURL(entry.got)
		if entry.want != result {
			t.Errorf("Normalization of '%s' was incorrect, got: '%s', want: '%s'", entry.got, result, entry.want)
		}
	}
}
