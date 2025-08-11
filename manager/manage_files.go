package manager

import (
	"fmt"
	"nt-storage/storage"
	"strings"
)

func (m *Manager) ListFiles() (string, int) {
	rootFid, err := m.wad.RootFid()
	if err != nil {
		return err.Error(), 404
	}
	files, err := storage.ListFiles(m.wad.Service, fmt.Sprintf("'%s' in parents", rootFid))
	if err != nil {
		return err.Error(), 500
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, file := range files {
		sb.WriteString(fmt.Sprintf("%d. %s '%s'\n", i, file.Id, file.Name))
	}
	sb.WriteRune(']')
	return sb.String(), 200
}

func (m *Manager) DeleteFile(fid string) (string, int) {
	err := storage.DeleteFile(m.wad.Service, fid)
	if err != nil {
		return err.Error(), 500
	}
	return fmt.Sprintf("Successfully deleted file '%s'", fid), 200
}
