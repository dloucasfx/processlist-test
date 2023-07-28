package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type procKey string

const version = "0.0.30"

var zlibCompressor = zlib.NewWriter(&bytes.Buffer{})
var now = time.Now

type TopProcess struct {
	ProcessID           int
	CreatedTime         time.Time
	Username            string
	Priority            int
	Nice                *int
	VirtualMemoryBytes  uint64
	WorkingSetSizeBytes uint64
	SharedMemBytes      uint64
	Status              string
	MemPercent          float64
	TotalCPUTime        time.Duration
	Command             string
}

func (tp *TopProcess) key() procKey {
	return procKey(fmt.Sprintf("%d|%s", tp.ProcessID, tp.Command))
}

// Monitor for Utilization
type Monitor struct {
	cancel        func()
	lastCPUCounts map[procKey]time.Duration
	nextPurge     time.Time
	logger        log.FieldLogger
}

func (m *Monitor) encodeEventMessage(procs []*TopProcess, sampleInterval time.Duration) (string, error) {
	if len(procs) == 0 {
		return "", errors.New("no processes to encode")
	}

	procsEncoded := []byte{'{'}
	for i := range procs {
		procsEncoded = append(procsEncoded, []byte(m.encodeProcess(procs[i], sampleInterval)+",")...)
	}
	procsEncoded[len(procsEncoded)-1] = '}'

	// escape and compress the process list
	escapedBytes := bytes.Replace(procsEncoded, []byte{byte('\\')}, []byte{byte('\\'), byte('\\')}, -1)
	compressedBytes, err := compressBytes(escapedBytes)
	if err != nil {
		return "", fmt.Errorf("couldn't compress process list: %v", err)
	}

	return fmt.Sprintf(
		"{\"t\":\"%s\",\"v\":\"%s\"}",
		base64.StdEncoding.EncodeToString(compressedBytes.Bytes()), version), nil
}

func (m *Monitor) encodeProcess(proc *TopProcess, sampleInterval time.Duration) string {
	key := proc.key()
	lastSampleInterval := sampleInterval
	lastCPUCount, ok := m.lastCPUCounts[key]
	if !ok {
		lastSampleInterval = time.Since(proc.CreatedTime)
	}
	m.lastCPUCounts[key] = proc.TotalCPUTime

	cpuPercent := float64(proc.TotalCPUTime-lastCPUCount) * 100.0 / float64(lastSampleInterval)

	nice := ""
	if proc.Nice == nil {
		nice = "unknown"
	} else {
		nice = strconv.Itoa(*proc.Nice)
	}

	return fmt.Sprintf(`"%d":["%s",%d,"%s",%d,%d,%d,"%s",%.2f,%.2f,"%s","%s"]`,
		proc.ProcessID,
		strings.ReplaceAll(proc.Username, `"`, "'"),
		proc.Priority,
		nice,
		proc.VirtualMemoryBytes/1024,
		proc.WorkingSetSizeBytes/1024,
		proc.SharedMemBytes/1024,
		proc.Status,
		cpuPercent,
		proc.MemPercent,
		toTime(proc.TotalCPUTime.Seconds()),
		strings.ReplaceAll(proc.Command, `"`, `'`),
	)
}

func (m *Monitor) purgeCPUCache(lastProcs []*TopProcess) {
	lastKeys := make(map[procKey]struct{}, len(lastProcs))
	for i := range lastProcs {
		lastKeys[lastProcs[i].key()] = struct{}{}
	}

	for k := range m.lastCPUCounts {
		if _, ok := lastKeys[k]; !ok {
			delete(m.lastCPUCounts, k)
		}
	}
}

// toTime returns the given seconds as a formatted string "min:sec.dec"
func toTime(secs float64) string {
	minutes := int(secs) / 60
	seconds := math.Mod(secs, 60.0)
	dec := math.Mod(seconds, 1.0) * 100
	return fmt.Sprintf("%02d:%02.f.%02.f", minutes, seconds, dec)
}

// compresses the given byte array
func compressBytes(in []byte) (buf bytes.Buffer, err error) {
	zlibCompressor.Reset(&buf)
	_, err = zlibCompressor.Write(in)
	_ = zlibCompressor.Close()
	return
}

func runOnInterval(ctx context.Context, fn func(), interval time.Duration) {
	timer := time.NewTicker(interval)

	go func() {
		fn()

		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				fn()
			}
		}
	}()
}

func main() {
	var ctx context.Context
	var quit = make(chan struct{})
	m := &Monitor{}

	logLevel := flag.String("log", "info", "debug : Display additional output")
	version := flag.Bool("version", false, "Display version")

	flag.Parse()
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Info("error")
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	m.logger = log.WithFields(log.Fields{"TestHarness": "Process_List"})
	m.logger.Debug("test")

	if *version {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			m.logger.Error("Getting build info failed (not in module mode?)!")
			return
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(bi); err != nil {
			panic(err)
		}
		return
	}
	ctx, m.cancel = context.WithCancel(context.Background())
	interval := time.Duration(10) * time.Second
	m.nextPurge = now().Add(3 * time.Minute)
	m.lastCPUCounts = make(map[procKey]time.Duration)
	osCache := initOSCache()
	procsOutput := ""

	runOnInterval(
		ctx,
		func() {
			procs, err := ProcessList(osCache, m.logger)
			if err != nil {
				m.logger.WithError(err).Error("Couldn't get process list")
				close(quit)
				return
			}

			if log.IsLevelEnabled(log.DebugLevel) {
				for i := range procs {
					procsOutput = procsOutput + "\n" + m.encodeProcess(procs[i], interval)
				}
				m.logger.Debugf("%s %s ----------------------------------------------------------------------------------------", procsOutput, "\n")
			}

			_, err = m.encodeEventMessage(procs, interval)
			if err != nil {
				m.logger.WithError(err).Error("Failed to encode process list")
				close(quit)
			}

			m.logger.Info("processes for this poll is collected and encoded")
		},
		interval)
	<-quit
	m.logger.Debug("exiting")
}
