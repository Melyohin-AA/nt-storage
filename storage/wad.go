package storage

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

type Wad struct {
	Config      *Config
	Service     *drive.Service
	DataFile    *DataFile
	rootFid     string
	dataFileFid string
}

func NewWad(config *Config) *Wad {
	return &Wad{Config: config}
}

func (wad *Wad) TryReadDataFile() (bool, error) {
	file, err := os.OpenFile("datafile", os.O_RDONLY, 0600)
	if err != nil {
		return false, nil
	}
	defer file.Close()
	df := new(DataFile)
	if err = df.Read(file); err != nil {
		return false, err
	}
	wad.DataFile = df
	return true, nil
}

func (wad *Wad) DoesRootExist() (bool, error) {
	if wad.rootFid != "" {
		return true, nil
	}
	fid, count, err := GetFid(wad.Service, fmt.Sprintf("name = '%s'", wad.Config.Root))
	if count == 1 {
		wad.rootFid = fid
		return true, nil
	}
	if count == 0 {
		return false, nil
	}
	if count > 1 {
		return false, fmt.Errorf("expected 1 root folder, %d found", count)
	}
	return false, err
}

func (wad *Wad) RootFid() (string, error) {
	if wad.rootFid != "" {
		return wad.rootFid, nil
	}
	fid, count, err := GetFid(wad.Service, fmt.Sprintf("name = '%s'", EscapeName(wad.Config.Root)))
	if err != nil {
		return "", err
	}
	if count > 1 {
		return "", fmt.Errorf("expected 1 root folder, %d found", count)
	}
	wad.rootFid = fid
	return fid, nil
}

func (wad *Wad) CreateRoot() error {
	if wad.rootFid != "" {
		return errors.New("root already exists")
	}
	file, err := CreateFolder(wad.Service, wad.Config.Root)
	if err != nil {
		return err
	}
	wad.rootFid = file.Id
	return nil
}

func (wad *Wad) DoesDataFileExist() (bool, error) {
	if wad.dataFileFid != "" {
		return true, nil
	}
	rootFid, err := wad.RootFid()
	if err != nil {
		return false, err
	}
	fid, count, err := GetFid(wad.Service, fmt.Sprintf("name = 'datafile' and '%s' in parents", rootFid))
	if count == 1 {
		wad.dataFileFid = fid
		return true, nil
	}
	if count == 0 {
		return false, nil
	}
	if count > 1 {
		return false, fmt.Errorf("expected 1 datafile, %d found", count)
	}
	return false, err
}

func (wad *Wad) DataFileFid() (string, error) {
	if wad.dataFileFid != "" {
		return wad.dataFileFid, nil
	}
	rootFid, err := wad.RootFid()
	if err != nil {
		return "", err
	}
	fid, count, err := GetFid(wad.Service, fmt.Sprintf("name = 'datafile' and '%s' in parents", rootFid))
	if err != nil {
		return "", err
	}
	if count > 1 {
		return "", fmt.Errorf("expected 1 datafile, %d found", count)
	}
	wad.dataFileFid = fid
	return fid, nil
}

func (wad *Wad) CreateDataFile() error {
	if wad.dataFileFid != "" {
		return errors.New("datafile already exists")
	}
	rootFid, err := wad.RootFid()
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(nil)
	if err = wad.DataFile.Write(body); err != nil {
		return err
	}
	encrypted, err := EncryptR(wad.Config.MasterKey, body)
	if err != nil {
		return err
	}
	remoteFile, err := CreateFile(wad.Service, "datafile", rootFid, encrypted)
	if err != nil {
		return err
	}
	wad.dataFileFid = remoteFile.Id
	return nil
}
