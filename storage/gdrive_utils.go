package storage

import (
	"fmt"
	"io"
	"strings"

	"google.golang.org/api/drive/v3"
)

// const (
// 	Space = "appDataFolder"
// )

func EscapeName(name string) string {
	return strings.NewReplacer(
		"\\", "\\\\",
		"'", "\\'",
	).Replace(name)
}

func ListFiles(service *drive.Service, q string) ([]*drive.File, error) {
	result, err := service.Files.
		List().Q(q).Fields("nextPageToken, files(id, name)").
		Do()
	if err != nil {
		return nil, err
	}
	return result.Files, nil
}

func GetFid(service *drive.Service, q string) (string, int, error) {
	result, err := service.Files.
		List().PageSize(2).Q(q).
		Do()
	if err != nil {
		return "", -1, err
	}
	if len(result.Files) == 0 {
		return "", 0, fmt.Errorf("cannot find file with \"%s\" query", q)
	}
	return result.Files[0].Id, len(result.Files), nil
}

func CreateFolder(service *drive.Service, name string) (*drive.File, error) {
	fileMetadata := &drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
	return service.Files.
		Create(fileMetadata).
		Do()
}

func CreateFile(service *drive.Service, name string, folderId string, body io.Reader) (*drive.File, error) {
	fileMetadata := &drive.File{Name: name, Parents: []string{folderId}}
	return service.Files.
		Create(fileMetadata).Media(body).
		Do()
}

func DeleteFile(service *drive.Service, fid string) error {
	return service.Files.Delete(fid).Do()
}

func FetchFile(service *drive.Service, fid string) (io.ReadCloser, error) {
	response, err := service.Files.Get(fid).Download()
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func UpdateFile(service *drive.Service, fid string, body io.Reader) error {
	fileMetadata := &drive.File{}
	file, err := service.Files.
		Update(fid, fileMetadata).Media(body).
		Do()
	if (err == nil) && (file.Id != fid) {
		return fmt.Errorf("expected file id to be '%s', but was '%s'", fid, file.Id)
	}
	return err
}
