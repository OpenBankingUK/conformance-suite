package report

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
)

const (
	marshalIndentPrefix = ""
	marshalIndent       = "  "
	reportFilename      = "report.json"
)

// Exporter - allows the exporting of a `Report`.
type Exporter interface {
	Export() error
}

type zipExporter struct {
	report Report
	writer io.Writer
}

// NewZipExporter - return new `Exporter` that exports to a ZIP archive to `writer`.
// The caller should close `writer` after calling `Export`.
//
// For example:
//     writer, err := os.Create("report.zip")
//     defer writer.Close()
//     exporter := NewZipExporter(Report{}, writer)
//     exporter.Export()
func NewZipExporter(report Report, writer io.Writer) Exporter {
	return &zipExporter{
		report: report,
		writer: writer,
	}
}

// Export - export `report` as a `.zip` to file named `filename`.
func (e *zipExporter) Export() error {
	// Create a new zip archive.
	zipWriter := zip.NewWriter(e.writer)
	defer zipWriter.Close()

	reportJSON, err := json.MarshalIndent(e.report, marshalIndentPrefix, marshalIndent)
	if err != nil {
		return fmt.Errorf("zip exporter cannot MarshalIndent report: %+v", err)
	}

	// Create file within ZIP archive
	reportFile, err := zipWriter.Create(reportFilename)
	if err != nil {
		return fmt.Errorf("zip exporter cannot Create %q file: %+v", reportFilename, err)
	}
	// Create report contents to zip
	if _, err := reportFile.Write(reportJSON); err != nil {
		// Only print the first 20 bytes of what we failed to write.
		return fmt.Errorf("zip exporter cannot Write %q: %+v", string(reportJSON), err)
	}

	return nil
}
