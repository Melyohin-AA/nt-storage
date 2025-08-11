package manager

import (
	"nt-storage/storage"
	"os"
)

func (m *Manager) InitRepo() (string, int) {
	// State validation
	rootExists, err := m.wad.DoesRootExist()
	if err != nil {
		return err.Error(), 500
	}
	if rootExists {
		dfExists, err := m.wad.DoesDataFileExist()
		if err != nil {
			return err.Error(), 500
		}
		if dfExists {
			return "cannot initialize: datafile already exists", 409
		}
	} else {
		if err = m.wad.CreateRoot(); err != nil {
			return err.Error(), 500
		}
	}
	m.wad.DataFile = storage.NewDataFile()
	// Creating local datafile
	file, err := os.OpenFile("datafile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	if err = m.wad.DataFile.Write(file); err != nil {
		return err.Error(), 500
	}
	// Creating remote datafile
	if err = m.wad.CreateDataFile(); err != nil {
		return err.Error(), 500
	}
	return "remote datafile has been successfully updated from local", 200
}

func (m *Manager) PushRepo() (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	fid, err := m.wad.DataFileFid()
	if err != nil {
		return err.Error(), 404
	}
	file, err := os.OpenFile("datafile", os.O_RDONLY, 0600)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	encrypted, err := storage.EncryptR(m.wad.Config.MasterKey, file)
	if err != nil {
		return err.Error(), 500
	}
	err = storage.UpdateFile(m.wad.Service, fid, encrypted)
	if err != nil {
		return err.Error(), 500
	}
	return "remote datafile has been successfully updated from local", 200
}

func (m *Manager) PullRepo() (string, int) {
	// Fetching data file
	fid, err := m.wad.DataFileFid()
	if err != nil {
		return err.Error(), 404
	}
	body, err := storage.FetchFile(m.wad.Service, fid)
	if err != nil {
		return err.Error(), 500
	}
	file, err := os.OpenFile("datafile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	if err := storage.Decrypt(m.wad.Config.MasterKey, file, body); err != nil {
		return err.Error(), 500
	}
	// Loading data file
	if _, err = file.Seek(0, 0); err != nil {
		return err.Error(), 500
	}
	var df storage.DataFile
	if err = df.Read(file); err != nil {
		return err.Error(), 500
	}
	m.wad.DataFile = &df
	return "local datafile has been successfully updated from remote", 200
}
