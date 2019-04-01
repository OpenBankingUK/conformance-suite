package report

import (
	"archive/zip"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "zipExporter.Export: json.MarshalIndent failed, report=%+v", e.report)
	}

	// Create file within ZIP archive
	reportFile, err := zipWriter.Create(reportFilename)
	if err != nil {
		return errors.Wrapf(err, "zipExporter.Export: zip.Writer.Create failed, could not create file %q", reportFilename)
	}
	// Create report contents to zip
	if _, err := reportFile.Write(reportJSON); err != nil {
		return errors.Wrapf(err, "zipExporter.Export: zip.Writer.Write failed, could write to %q, reportJSON=%+v", reportFilename, string(reportJSON))
	}

	return nil
}
