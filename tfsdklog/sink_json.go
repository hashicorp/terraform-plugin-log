package tfsdklog

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"

	"github.com/hashicorp/go-hclog"
)

var stdErrBufferSize = 64 * 1024

// PipeJSONLogs processes go-hclog-formatted JSON output from `r` and writes it
// to `sink`, using any existing sub-loggers it can identify and falling back
// on the cli sub-logger if no sub-logger matches the module for the log entry.
// It will block until `ctx` is canceled.
//
// PipeJSONLogs is meant to be used to process output from a Terraform binary
// run via terraform-exec, merging the logs from the provider under test with
// the logs from Terraform and other providers.
//
// It is safe to run multiple instances of PipeJSONLogs concurrently without
// any coordination as long as each is the only reader of `r`.
func PipeJSONLogs(ctx context.Context, commandID string, r io.Reader) {
	reader := bufio.NewReaderSize(r, stdErrBufferSize)
	// continuation indicates the previous line was a prefix
	continuation := false

	for {
		select {
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				sink.Debug("context closed", "error", ctx.Err())
			}
			return
		default:
		}
		line, isPrefix, err := reader.ReadLine()
		switch {
		case err == io.EOF:
			return
		case err != nil:
			sink.Error("reading JSON logs", "error", err)
			return
		}

		// The line was longer than our max token size, so it's likely
		// incomplete and won't unmarshal.
		if isPrefix || continuation {
			sink.Error("log line larger than max log line size", "prefix", string(line))
			continuation = isPrefix
			continue
		}

		entry, err := parseJSON(line)
		if err != nil {
			sink.Error("parsing JSON logs", "input", line, "error", err)
			return
		}
		out := flattenKVPairs(entry.KVPairs)

		l := newCLILogger(ctx, commandID)
		if entry.Module != "" {
			l = l.Named(entry.Module)
		}

		out = append(out, "timestamp", entry.Timestamp.Format(hclog.TimeFormat))
		switch hclog.LevelFromString(entry.Level) {
		case hclog.Trace:
			l.Trace(entry.Message, out...)
		case hclog.Debug:
			l.Debug(entry.Message, out...)
		case hclog.Info:
			l.Info(entry.Message, out...)
		case hclog.Warn:
			l.Warn(entry.Message, out...)
		case hclog.Error:
			l.Error(entry.Message, out...)
		default:
			// if there was no log level, it's likely this is unexpected
			// json from something other than hclog, and we should output
			// it verbatim.
			sink.Error("parsing JSON logs", "input", string(line), "error", errors.New("no log level"))
		}
	}
}

// logEntry is the JSON payload that gets written by go-hclog.
type logEntry struct {
	Message   string        `json:"@message"`
	Level     string        `json:"@level"`
	Timestamp time.Time     `json:"@timestamp"`
	Module    string        `json:"@module"`
	KVPairs   []*logEntryKV `json:"kv_pairs"`
}

// logEntryKV is a key value pair within the go-hclog JSON payload.
type logEntryKV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// flattenKVPairs is used to flatten KVPair slice into []interface{}
// for hclog consumption.
func flattenKVPairs(kvs []*logEntryKV) []interface{} {
	var result []interface{}
	for _, kv := range kvs {
		result = append(result, kv.Key)
		result = append(result, kv.Value)
	}

	return result
}

// parseJSON handles parsing JSON output
func parseJSON(input []byte) (*logEntry, error) {
	var raw map[string]interface{}
	entry := &logEntry{}

	err := json.Unmarshal(input, &raw)
	if err != nil {
		return nil, err
	}

	// Parse hclog-specific objects
	if v, ok := raw["@message"]; ok {
		entry.Message = v.(string)
		delete(raw, "@message")
	}

	if v, ok := raw["@level"]; ok {
		entry.Level = v.(string)
		delete(raw, "@level")
	}

	if v, ok := raw["@timestamp"]; ok {
		t, err := time.Parse("2006-01-02T15:04:05.000000Z07:00", v.(string))
		if err != nil {
			return nil, err
		}
		entry.Timestamp = t
		delete(raw, "@timestamp")
	}

	if v, ok := raw["@module"]; ok {
		entry.Module = v.(string)
		delete(raw, "@module")
	}

	// Parse dynamic KV args from the hclog payload.
	for k, v := range raw {
		entry.KVPairs = append(entry.KVPairs, &logEntryKV{
			Key:   k,
			Value: v,
		})
	}

	sort.Slice(entry.KVPairs, func(i, j int) bool {
		return entry.KVPairs[i].Key < entry.KVPairs[j].Key
	})

	return entry, nil
}
