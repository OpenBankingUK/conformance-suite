package report

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Importer - allows the importing of a `Report`.
type Importer interface {
	Import() (Report, error)
}

type zipImporter struct {
	reader io.Reader
}

// NewZipImporter - return new `Importer` that allows importing from a ZIP archive.
func NewZipImporter(reader io.Reader) Importer {
	return &zipImporter{
		reader: reader,
	}
}

// Import - import `report.json` from `reader`.
func (i *zipImporter) Import() (Report, error) {
	// We need to determine the length of `i.reader`, so read until EOF and use that as the size to the call to `zip.NewReader`.
	// Could possibly use one of these libraries:
	// * https://godoc.org/go4.org/readerutil#NewBufferingReaderAt
	// * https://github.com/go4org/go4/blob/master/readerutil/bufreaderat.go#L23
	// But `go4` makes no backwards compatibility promises.
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, i.reader)
	if err != nil {
		return Report{}, errors.Wrapf(err, "zipImporter.Import: io.Copy failed, copied=%d bytes", size)
	}

	reader := bytes.NewReader(buff.Bytes())

	// Open a zip archive for reading.
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return Report{}, errors.Wrapf(err, "zipImporter.Import: zip.NewReader failed, could not get new zip.Reader with size=%d", size)
	}

	for _, file := range zipReader.File {
		// Ignore anything that isn't `reportFilename`.
		if file.Name != reportFilename {
			continue
		}

		// Open `reportFilename` and read its contents, then Unmarshall to `Report` struct.
		readerCloser, err := file.Open()
		if err != nil {
			return Report{}, errors.Wrapf(err, "zipImporter.Import: file.Open failed, could not open %q", reportFilename)
		}

		buff := bytes.NewBuffer([]byte{})
		if size, err := io.Copy(buff, readerCloser); err != nil {
			readerCloser.Close()
			return Report{}, errors.Wrapf(err, "zipImporter.Import: io.Copy failed, copied=%d bytes from %q", size, reportFilename)
		}

		report := Report{}
		if err := json.Unmarshal(buff.Bytes(), &report); err != nil {
			readerCloser.Close()
			return Report{}, errors.Wrapf(err, "zipImporter.Import: json.Unmarshal failed, could not marshall %q to Report", reportFilename)
		}

		if err := readerCloser.Close(); err != nil {
			return Report{}, errors.Wrapf(err, "zipImporter.Import: file.Close failed, could not close %q", reportFilename)
		}

		return report, nil
	}

	return Report{}, fmt.Errorf("zipImporter.Import: could not find %q in ZIP archive", reportFilename)
}
