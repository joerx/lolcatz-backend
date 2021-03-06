package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/joerx/lolcatz-backend/test/dbtest"
	"github.com/joerx/lolcatz-backend/upload"
)

func Test_listUploads(t *testing.T) {
	// truncate database
	if err := dbtest.Truncate(dbc); err != nil {
		t.Fatal(err)
	}

	// seed uploads
	uploads, err := dbtest.SeedUploads(dbc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/list"), nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected statuscode %d but got %d", e, a)
	}

	body, _ := ioutil.ReadAll(w.Body)
	items := make([]map[string]interface{}, 0)

	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatal(err)
	}

	if e, a := len(uploads), len(items); e != a {
		t.Errorf("Expected %d uploads, got %d", e, a)
	}
}

func Test_postUpload(t *testing.T) {
	if err := dbtest.Truncate(dbc); err != nil {
		t.Fatal(err)
	}

	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)

	// Attach test file
	fileName := "test/files/cat.jpg"
	mediaData, err := ioutil.ReadFile("./files/cat.jpg")
	if err != nil {
		t.Fatal(err)
	}

	mediaPart, err := form.CreateFormFile("file", path.Base(fileName))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(mediaPart, bytes.NewReader(mediaData)); err != nil {
		t.Fatal(err)
	}
	form.Close()

	req, err := http.NewRequest(http.MethodPost, "/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", fmt.Sprintf("multipart/form-data; boundary=%s", form.Boundary()))

	// Make request
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if e, a := http.StatusOK, w.Result().StatusCode; e != a {
		t.Errorf("Expected status %d, got %d", e, a)
	}

	// Parse response
	data, _ := ioutil.ReadAll(w.Result().Body)
	var upload upload.Upload

	if err := json.Unmarshal(data, &upload); err != nil {
		t.Fatal(err)
	}

	if id := upload.ID; id == 0 {
		t.Errorf("Expected returned record to contain a non-zero id, got %d", id)
	}
}
