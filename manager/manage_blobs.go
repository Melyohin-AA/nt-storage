package manager

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"nt-storage/storage"
	"os"
	"strings"
	"time"
)

func (m *Manager) ListBlobs(filter Filter) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	var i int
	m.wad.DataFile.EnumerateBlobs(func(blob *storage.Blob) bool {
		sb.WriteString(fmt.Sprintf("%d.", i))
		filter.apply(&sb, blob, m.wad.DataFile)
		sb.WriteRune('\n')
		i++
		return true
	})
	sb.WriteRune(']')
	return sb.String(), 200
}

func (m *Manager) AddBlob(name string) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	fi, err := os.Stat(name)
	if err != nil {
		return err.Error(), getFileStatErrorCode(err)
	}
	key := make([]byte, aes.BlockSize)
	rand.Read(key)
	file, err := os.OpenFile(name, os.O_RDONLY, 0700)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	encrypted, err := storage.EncryptR(key, file)
	if err != nil {
		return err.Error(), 500
	}
	rootFid, err := m.wad.RootFid()
	if err != nil {
		return err.Error(), 500
	}
	f, err := storage.CreateFile(m.wad.Service, "blob", rootFid, encrypted)
	if err != nil {
		return err.Error(), 500
	}
	blob := &storage.Blob{
		Fid:     f.Id,
		Key:     key,
		Name:    name,
		ModTime: fi.ModTime().UnixMilli(),
	}
	if err = m.wad.DataFile.AddBlob(blob); err != nil {
		return err.Error(), 500
	}
	if err = m.saveRepo(); err != nil {
		return err.Error(), 500
	}
	return fmt.Sprintf("blob successfully added and uploaded (fid: %s)", f.Id), 200
}

func (m *Manager) UpdateBlob(fid string) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	blob := m.wad.DataFile.GetBlob(fid)
	if blob == nil {
		return "blob not found", 404
	}
	fi, err := os.Stat(blob.Name)
	if err != nil {
		return err.Error(), getFileStatErrorCode(err)
	}
	blob.ModTime = fi.ModTime().UnixMilli()
	if err = m.saveRepo(); err != nil {
		return err.Error(), 500
	}
	// Updating remote blob
	file, err := os.OpenFile(blob.Name, os.O_RDONLY, 0700)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	encrypted, err := storage.EncryptR(blob.Key, file)
	if err != nil {
		return err.Error(), 500
	}
	if err = storage.UpdateFile(m.wad.Service, blob.Fid, encrypted); err != nil {
		return err.Error(), 500
	}
	return fmt.Sprintf("blob successfully updated (fid: %s)", fid), 200
}

func (m *Manager) FetchBlob(fid string) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	blob := m.wad.DataFile.GetBlob(fid)
	if blob == nil {
		return "blob not found", 404
	}
	body, err := storage.FetchFile(m.wad.Service, fid)
	if err != nil {
		return err.Error(), 500
	}
	file, err := os.OpenFile(escapeFsRestrictedChars(blob.Name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err.Error(), 500
	}
	defer file.Close()
	if err := storage.Decrypt(blob.Key, file, body); err != nil {
		return err.Error(), 500
	}
	os.Chtimes(blob.Name, time.Time{}, time.UnixMilli(blob.ModTime))
	return fmt.Sprintf("blob successfully fetched (fid: %s)", fid), 200
}

func (m *Manager) DeleteBlob(fid string) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	blob := m.wad.DataFile.GetBlob(fid)
	if blob == nil {
		return "blob not found", 404
	}
	if err := storage.DeleteFile(m.wad.Service, fid); err != nil {
		return err.Error(), 500
	}
	m.wad.DataFile.RemoveBlob(blob.Fid)
	if err := m.saveRepo(); err != nil {
		return err.Error(), 500
	}
	return fmt.Sprintf("blob successfully deleted (fid: %s)", fid), 200
}

func (m *Manager) GetBlobMeta(fid string) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	blob := m.wad.DataFile.GetBlob(fid)
	if blob == nil {
		return "blob not found", 404
	}
	encoded, err := json.Marshal(blob)
	if err != nil {
		return err.Error(), 500
	}
	return string(encoded), 200
}

func (m *Manager) SetBlobMeta(fid string, meta storage.Blob) (string, int) {
	if output, code := m.requireDataFile(); code != 0 {
		return output, code
	}
	blob := m.wad.DataFile.GetBlob(fid)
	if blob == nil {
		return "blob not found", 404
	}
	blob.Name = meta.Name
	blob.ModTime = meta.ModTime
	blob.Comment = meta.Comment
	blob.TagValues = meta.TagValues
	if err := m.saveRepo(); err != nil {
		return err.Error(), 500
	}
	return fmt.Sprintf("blob '%s' meta has been successfully updated", fid), 200
}
