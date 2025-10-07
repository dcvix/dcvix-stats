//  SPDX-FileCopyrightText: 2025 Diego Cortassa
//  SPDX-License-Identifier: MIT

package logparser

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/dcvix/dcvix-stats/internal/globals"
	"github.com/dcvix/dcvix-stats/internal/logger"
)

type LogEntry struct {
	Timestamp string
	Metric    string
	LastValue float64
}

type LogParser struct {
	filename string
	metrics  []string
	entries  []LogEntry
	regex    *regexp.Regexp
}

func NewLogParser(filename string) *LogParser {

	// Regex to match the log line and extract timestamp, metric, and last value
	// 2025-09-26 10:39:33,895159 [  1139:1139  ] INFO  quictransport - Connection 3 - Stats (1): quic_lost_packets: [sum: 221, last: 221, max: 221, avg: 221.00]
	// regex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}),\d+ .* (quic_\w+|intermediates_rtt_nanos): \[.*last: ([0-9.]+),.*\]`)
	regex := regexp.MustCompile(`^(\S+\s+\S+),.*Stats \(\d+\): (\S+):.*last: ([0-9]+).*avg: ([0-9]+)`)

	return &LogParser{
		filename: filename,
		metrics:  globals.Metrics,
		entries:  make([]LogEntry, 0),
		regex:    regex,
	}
}

func (lp *LogParser) ReadLogFile() error {
	file, err := os.Open(lp.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var newEntries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		entry := lp.parseLine(line)
		if entry != nil {
			newEntries = append(newEntries, entry...)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	lp.entries = newEntries
	return nil
}

func (lp *LogParser) parseLine(line string) []LogEntry {
	matches := lp.regex.FindStringSubmatch(line)
	if len(matches) != 5 {
		return nil
	}

	timestampUTC := matches[1]
	metric := matches[2]
	lastValueStr := matches[3]
	avgValueStr := matches[4]

	// Check if this metric is one we're interested in
	found := false
	for _, m := range lp.metrics {
		if m == metric {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	// convert UTC timestamp to local time
	utcTime, err := time.ParseInLocation("2006-01-02 15:04:05", timestampUTC, time.UTC)
	if err != nil {
		logger.LogVerbose("Error parsing timestamp: %v", err)
		return nil
	}
	localTime := utcTime.Local()
	timestampLocalTime := localTime.Format("15:04:05")

	lastValue, err := strconv.ParseFloat(lastValueStr, 64)
	if err != nil {
		return nil
	}

	avgValue, err := strconv.ParseFloat(avgValueStr, 64)
	if err != nil {
		return nil
	}

	valEntry := LogEntry{
		Timestamp: timestampLocalTime,
		Metric:    metric,
		LastValue: lastValue,
	}

	avgEntry := LogEntry{
		Timestamp: timestampLocalTime,
		Metric:    metric + "_avg",
		LastValue: avgValue,
	}

	res := []LogEntry{valEntry, avgEntry}
	return res
}

func (lp *LogParser) GetEntriesByMetric(metric string) []LogEntry {
	var entries []LogEntry
	for _, entry := range lp.entries {
		if entry.Metric == metric {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (lp *LogParser) GetEntriesByMetricList(metrics []string) ([][]float64, []string) {
	var values [][]float64
	var timeStamps []string

	for i, metric := range metrics {
		entries := lp.GetEntriesByMetric(metric)

		start := 0
		if len(entries) > globals.LogEntriesQty {
			start = len(entries) - globals.LogEntriesQty
		}

		var metricValues []float64
		for j := start; j < len(entries); j++ {
			metricValues = append(metricValues, entries[j].LastValue)
			// Only add timestamps for the first metric to avoid duplicates
			if i == 0 {
				timeStamps = append(timeStamps, entries[j].Timestamp)
			}
		}
		if len(metricValues) > 0 {
			values = append(values, metricValues)
		}
	}

	return values, timeStamps
}
