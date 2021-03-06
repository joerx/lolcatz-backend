package mp

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/joerx/lolcatz-backend/http/errors"
)

// UploadedFile represents the result of a file upload processed via GetUploadedFile
type UploadedFile struct {
	TmpFile     *os.File
	Filename    string
	ContentType string
}

func downloadTempFile(f multipart.File, h *multipart.FileHeader) (*os.File, error) {
	log.Printf("Received file '%s', size is %d bytes", h.Filename, h.Size)

	defer f.Close()

	tmpFile, err := ioutil.TempFile("", "lolcatz-upload-")
	if err != nil {
		return nil, err
	}

	defer tmpFile.Close()

	numBytes, err := io.Copy(tmpFile, f)
	if err != nil {
		return nil, err
	}

	log.Printf("%d bytes written to %s", numBytes, tmpFile.Name())
	return tmpFile, nil
}

// GetUploadedFile extracts the uploaded file from a multipart file upload, copies it into a
// temporary file and returns the tmp file location along with some metadata
func GetUploadedFile(r *http.Request, fieldName string) (*UploadedFile, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, errors.BadRequest("No file in upload")
	}

	tmpFile, err := downloadTempFile(file, header)
	if err != nil {
		return nil, err
	}

	uf := &UploadedFile{
		TmpFile:     tmpFile,
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
	}

	return uf, nil
}
