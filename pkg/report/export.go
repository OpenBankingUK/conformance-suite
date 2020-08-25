package report

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	marshalIndentPrefix    = ""
	marshalIndent          = "  "
	reportFilename         = "report.json"
	discoveryFilename      = "discovery.json"
	responseFieldsFilename = "responseFields.json"
)

var (
	// ErrExportFailure is the common error type returned on all export errors
	ErrExportFailure = errors.New("export failed")

	// Not a real secret; it's used when calculating a checksum for the exported report files.
	// The main purpose is to protect against accidental edits in the exported files.
	exportSecret = []byte(os.Getenv("EXPORT_SECRET"))
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

	toExport, err := readManifestFiles(&e.report)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrExportFailure, err)
	}

	reportJSON, err := json.MarshalIndent(e.report, marshalIndentPrefix, marshalIndent)
	if err != nil {
		return fmt.Errorf("%w: json.MarshalIndent failed: %s, report=%+v", ErrExportFailure, err.Error(), e.report)
	}

	discoveryJSON, err := json.MarshalIndent(e.report.Discovery, marshalIndentPrefix, marshalIndent)
	if err != nil {
		return fmt.Errorf("%w: json.MarshalIndent failed: %s, discovery=%+v", ErrExportFailure, err.Error(), e.report.Discovery)
	}

	toExport[reportFilename] = reportJSON
	toExport[discoveryFilename] = discoveryJSON
	toExport[responseFieldsFilename] = []byte(e.report.ResponseFields)
	toExport["report.checksum"] = createChecksum(exportSecret, reportJSON)

	return writeFiles(zipWriter, toExport)
}

func createChecksum(data, secret []byte) []byte {
	h := sha256.New()
	h.Write(data)   // never returns non-nil error
	h.Write(secret) //

	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func writeFiles(zipWriter *zip.Writer, files map[string][]byte) error {
	for fileName, contents := range files {
		header := &zip.FileHeader{
			Name:     fileName,
			Method:   zip.Deflate,
			Modified: time.Now(),
		}

		f, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("error writing file: '%s': %s", fileName, err.Error())
		}

		if _, err := f.Write(contents); err != nil {
			return fmt.Errorf("error writing file: '%s': %s", fileName, err.Error())
		}
	}
	return nil
}

func readManifestFiles(report *Report) (map[string][]byte, error) {
	manifests := map[string][]byte{}
	for _, manifestPath := range report.manifestFilePaths() {
		fileContents, err := readManifestFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read manifest file: %w", err)
		}

		fileName := path.Base(manifestPath)
		manifests[fileName] = fileContents
	}
	return manifests, nil
}

func readManifestFile(path string) ([]byte, error) {
	searchPaths := []string{".", "../.."}
	var err error
	for _, searchPath := range searchPaths {
		filePath := strings.Join([]string{searchPath, path}, "/")
		var fileContents []byte
		fileContents, err = ioutil.ReadFile(filePath)
		if err == nil {
			return fileContents, nil
		}
		if os.IsNotExist(err) {
			continue
		}
		break
	}
	return nil, err
}
