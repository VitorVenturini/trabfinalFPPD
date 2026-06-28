package metrics

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"trabfinalfppd/internal/matrix"
)

type RunResult struct {
	Mode         string
	N            int
	Seed         int64
	Nodes        int
	Processes    int
	ElapsedSec   float64
	Label        string
	Verification matrix.Verification
}

func (r RunResult) CSVHeader() string {
	return "mode,n,seed,nodes,processes,elapsed_sec,label,c00,c0n,cn0,cnn,checksum\n"
}

func (r RunResult) CSVRow() string {
	return fmt.Sprintf(
		"%s,%d,%d,%d,%d,%.9f,%s,%.15f,%.15f,%.15f,%.15f,%.15f\n",
		r.Mode,
		r.N,
		r.Seed,
		r.Nodes,
		r.Processes,
		r.ElapsedSec,
		escapeCSV(r.Label),
		r.Verification.TopLeft,
		r.Verification.TopRight,
		r.Verification.BottomLeft,
		r.Verification.BottomRight,
		r.Verification.Checksum,
	)
}

func AppendCSV(path string, result RunResult) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	fileInfo, err := os.Stat(path)
	writeHeader := false
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		writeHeader = true
	} else if fileInfo.Size() == 0 {
		writeHeader = true
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	if writeHeader {
		if _, err := file.WriteString(result.CSVHeader()); err != nil {
			return err
		}
	}

	_, err = file.WriteString(result.CSVRow())
	return err
}

func escapeCSV(value string) string {
	if value == "" {
		return ""
	}

	escaped := strings.ReplaceAll(value, "\"", "\"\"")
	return fmt.Sprintf("\"%s\"", escaped)
}
