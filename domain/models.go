package domain

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type TStatuses string

const (
	WAITING_LINKS TStatuses = "WAITING_LINKS"
	PROCESSING    TStatuses = "PROCESSING"
	FINISHED      TStatuses = "FINISHED"
)

type Task struct {
	Uuid        string
	TaskStatus  TStatuses
	NumOfLinks  int64
	ArchiveName string
	ArchiveLink string
	Version     int64
}

type FTypes string

const (
	UNKNOWN_TYPE FTypes = "UNKNOWN_TYPE"
	PDF          FTypes = "pdf"
	JPEG         FTypes = "jpeg"
)

type File struct {
	TaskUuid     string
	Link         string
	Path         string
	ErrorMessage string
	Version      int64
}

func DownloadFile(url, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ValidateFile(path string) FTypes {
	f, e := os.Open(path)
	if e != nil {
		logrus.Errorf("Open file \"%s\" error: %s", path, e)
		return UNKNOWN_TYPE
	}
	defer func() {
		if e := f.Close(); e != nil {
			logrus.Errorf("Close file \"%s\" error: %s", path, e)
		}
	}()

	buf := make([]byte, 512)
	if _, e := f.Read(buf); e != nil {
		logrus.Errorf("Read file \"%s\" error: %s", path, e)
		return UNKNOWN_TYPE
	}

	t := http.DetectContentType(buf)
	switch FTypes(t) {
	case "image/jpeg", "application/pdf":
		return FTypes(t)
	default:
		return UNKNOWN_TYPE
	}
}

func CopyFile(from, to string) error {
	in, e1 := os.OpenFile(from, os.O_RDONLY, 0)
	if e1 != nil {
		logrus.Errorf("Open file [from] \"%s\" error: %s", from, e1)
		return e1
	}
	defer in.Close()

	out, e2 := os.OpenFile(from, os.O_WRONLY, 0)
	if e2 != nil {
		logrus.Errorf("Open file [to] \"%s\" error: %s", to, e2)
		return e2
	}
	defer out.Close()

	if _, e := io.Copy(out, in); e != nil {
		logrus.Errorf("Copy file \"%s\" error: %s", to, e)
		return e
	}

	return nil
}

func PackToArchiver(name string, files []*File) error {
	w := zip.NewWriter(new(bytes.Buffer))
	defer w.Close()

	for _, file := range files {
		fBytes, err := os.ReadFile(file.Path)
		if err != nil {
			logrus.Errorf("Read file \"%s\" error: %s", file.Path, err)
			continue
		}

		f, err := w.Create(file.Path)
		if err != nil {
			logrus.Errorf("Create file \"%s\" in zip error: %s", file.Path, err)
			continue
		}

		_, e := f.Write(fBytes)
		if e != nil {
			logrus.Errorf("Write file \"%s\" to zip error: %s", file.Path, e)
		}
	}

	return nil
}
