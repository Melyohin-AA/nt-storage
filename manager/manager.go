package manager

import (
	"nt-storage/storage"
	"os"
)

type Manager struct {
	wad *storage.Wad
}

func NewManager(wad *storage.Wad) Manager {
	return Manager{wad: wad}
}

func (m *Manager) saveRepo() error {
	file, err := os.OpenFile("datafile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return m.wad.DataFile.Write(file)
}

func (m *Manager) requireDataFile() (string, int) {
	if m.wad.DataFile == nil {
		return "local datafile is missing", 404
	}
	return "", 0
}

func getFileStatErrorCode(err error) int {
	if os.IsNotExist(err) {
		return 404
	}
	if os.IsPermission(err) {
		return 403
	}
	return 400
}
