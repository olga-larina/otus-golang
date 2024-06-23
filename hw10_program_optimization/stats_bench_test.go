//go:build bench
// +build bench

package hw10programoptimization

import (
	"archive/zip"
	"io"
	"testing"
)

// go test -v -count=10 -timeout=30s -tags bench -run '^$' -bench ^BenchmarkStats$ -benchmem | tee new.txt
// benchstat old.txt new.txt | tee benchstat.txt.
func BenchmarkStats(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Errorf("failed opening zip file %v", err)
	}
	if len(r.File) != 1 {
		b.Errorf("wrong number of files in zip: %d", len(r.File))
	}

	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			b.Errorf("failed closing zip file %v", err)
		}
	}(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		data, err := r.File[0].Open()
		if err != nil {
			b.Errorf("failed opening file: %v", err)
		}
		defer func(rc io.ReadCloser) {
			err := rc.Close()
			if err != nil {
				b.Errorf("failed closing file %v", err)
			}
		}(data)

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()

		if err != nil {
			b.Errorf("failed getting users: %v", err)
		}
	}
}
