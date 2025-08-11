package main

import (
	"fmt"
	"nt-storage/server"
	"nt-storage/storage"
)

// https://developers.google.com/workspace/drive/api/quickstart/go
// https://github.com/douglasmakey/oauth2-example

func main() {
	// key := []byte{2, 6, 2, 79, 36, 123, 95, 34, 8, 3, 78, 234, 54, 67, 34, 1}
	config, err := storage.ConfigFromFile("config.json")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		return
	}
	wad := storage.NewWad(config)
	dfRead, err := wad.TryReadDataFile()
	if err != nil {
		fmt.Printf("Local datafile is corrupted: %v\n", err)
		return
	}
	if dfRead {
		fmt.Println("Datafile loaded from file")
	}
	s, err := server.NewServer(wad)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting server")
	if err = s.Run(); err != nil {
		panic(err)
	}
	// listFiles(wad)
	// getIndexFid(config)
	// createFile(config, "example")
	// deleteFile(config, "")
	// fetchFile(config, "", "example2", key)
	// updateFile(config, "", "example", key)
}

// func listFiles(wad *storage.Wad) {
// 	files, err := storage.ListFiles(wad.Service)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Files:")
// 	for i, file := range files {
// 		fmt.Printf("%d. %s - %s\n", i, file.Id, file.Name)
// 	}
// }

// func getIndexFid(config storage.Config) {
// 	fid, err := storage.GetFid(config.Service, "index")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Id of the index file is %s\n", fid)
// }

// func createFile(config storage.Config, name string, key []byte) {
// 	body, err := os.ReadFile(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cbody, err := storage.EncryptR(key, bytes.NewReader(body))
// 	if err != nil {
// 		panic(err)
// 	}
// 	file, err := storage.CreateFile(config.Service, name, cbody)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Id of the created file is %s\n", file.Id)
// }

// func deleteFile(config storage.Config, fid string) {
// 	err := storage.DeleteFile(config.Service, fid)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Successfully deleted file '%s'\n", fid)
// }

// func fetchFile(config storage.Config, fid string, name string, key []byte) {
// 	body, err := storage.FetchFile(config.Service, fid)
// 	if err != nil {
// 		panic(err)
// 	}
// 	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	if err := storage.Decrypt(key, file, body); err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Successfully wrote local file '%s'\n", name)
// }

// func updateFile(config storage.Config, fid string, name string, key []byte) {
// 	body, err := os.ReadFile(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	cbody, err := storage.EncryptR(key, bytes.NewReader(body))
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = storage.UpdateFile(config.Service, fid, cbody)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Successfully updated file '%s'\n", fid)
// }
