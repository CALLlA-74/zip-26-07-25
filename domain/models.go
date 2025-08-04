package domain

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"strings"

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
	FileUid      string
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
	defer func() {
		if e := out.Close(); e != nil {
			logrus.Errorf("[Downloading] Error close [out] file \"%s\" error: %s", path, e.Error())
		}
	}()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			logrus.Errorf("[Downloading] Error close [resp.body] \"%s\" error: %s", url, e.Error())
		}
	}()

	logrus.Infof("Download file status: %s", resp.Status)
	if resp.StatusCode >= http.StatusBadRequest {
		return ErrLoadFile
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

func ValidateFile(path string) FTypes {
	logrus.Infof("[Validating] Try to validate file [%s]", path)
	f, e := os.Open(path)
	if e != nil {
		logrus.Errorf("[Validating] Open file \"%s\" error: %s", path, e)
		return UNKNOWN_TYPE
	}
	defer func() {
		if e := f.Close(); e != nil {
			logrus.Errorf("[Validating] Close file \"%s\" error: %s", path, e)
		}
	}()

	buf := make([]byte, 512)
	if _, e := f.Read(buf); e != nil {
		logrus.Errorf("[Validating] Read file \"%s\" error: %s", path, e)
		return UNKNOWN_TYPE
	}

	t := http.DetectContentType(buf)
	switch FTypes(t) {
	case "image/jpeg":
		return JPEG
	case "application/pdf":
		return PDF
	default:
		return UNKNOWN_TYPE
	}
}

func CopyFile(from, to string) error {
	logrus.Infof("[Coping] Try to copy from [%s] to [%s]", from, to)
	in, e1 := os.OpenFile(from, os.O_RDONLY, 0)
	if e1 != nil {
		logrus.Errorf("[Coping] Open file [from] \"%s\" error: %s", from, e1)
		return e1
	}
	defer func() {
		if e := in.Close(); e != nil {
			logrus.Errorf("[Coping] Close file [from] \"%s\" error: %s", from, e1)
		}
	}()

	out, e2 := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, 0)
	if e2 != nil {
		logrus.Errorf("[Coping] Open file [to] \"%s\" error: %s", to, e2)
		return e2
	}
	defer func() {
		if e := out.Close(); e != nil {
			logrus.Errorf("[Coping] Close file [to] \"%s\" error: %s", to, e2)
		}
	}()

	if _, e := io.Copy(out, in); e != nil {
		logrus.Errorf("[Coping] Copy file \"%s\" error: %s", to, e)
		return e
	}

	return nil
}

func PackToArchiver(name string, files []*File) error {
	archFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		logrus.Errorf("[Packing] Open archiver file \"%s\" error: %s", name, err)
		return err
	}
	defer func() {
		if e := archFile.Close(); e != nil {
			logrus.Errorf("[Packing] Close archiver file \"%s\" error: %s", name, err)
		}
	}()

	w := zip.NewWriter(archFile)
	defer func() {
		if e := w.Close(); e != nil {
			logrus.Errorf("[Packing] Close zipWriter error: %s", e)
		}
	}()

	for _, file := range files {
		fBytes, err := os.ReadFile(file.Path)
		if err != nil {
			logrus.Errorf("[Packing] Read file \"%s\" error: %s", file.Path, err)
			continue
		}
		splt := strings.Split(file.Path, "/")

		f, err := w.Create(splt[len(splt)-1])
		if err != nil {
			logrus.Errorf("[Packing] Create file \"%s\" in zip error: %s", file.Path, err)
			continue
		}

		_, e := f.Write(fBytes)
		if e != nil {
			logrus.Errorf("[Packing] Write file \"%s\" to zip error: %s", file.Path, e)
		}
	}

	return nil
}
