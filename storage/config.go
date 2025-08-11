package storage

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	TokenPath = "token.json"
)

type ConfigModel struct {
	Port      int    `json:"port"`
	CredPath  string `json:"credpath"`
	MasterKey string `json:"masterkey"`
	Root      string `json:"root"`
}

type Config struct {
	Port        int
	OauthConfig *oauth2.Config
	MasterKey   []byte
	Root        string
}

func ConfigFromFile(path string) (*Config, error) {
	bin, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var model ConfigModel
	if err = json.Unmarshal(bin, &model); err != nil {
		return nil, err
	}
	cfg, err := ConfigFromModel(model)
	if err != nil {
		return nil, err
	}
	if model.MasterKey == "" {
		model.MasterKey = hex.EncodeToString(cfg.MasterKey)
		if bin, err = json.Marshal(model); err != nil {
			return cfg, err
		}
		if err = os.WriteFile(path, bin, 0600); err != nil {
			return cfg, err
		}
	}
	cfg.Root = model.Root
	return cfg, nil
}

func ConfigFromModel(model ConfigModel) (*Config, error) {
	oauthConfig, err := getOauthConfig(model.CredPath)
	if err != nil {
		return nil, err
	}
	var masterkey []byte
	if model.MasterKey != "" {
		masterkey, err = hex.DecodeString(model.MasterKey)
		if err != nil {
			return nil, err
		}
		if len(masterkey) != aes.BlockSize {
			return nil, fmt.Errorf(
				"expected master-key to be %d bytes long but was %d", aes.BlockSize, len(masterkey),
			)
		}
	} else {
		masterkey = make([]byte, aes.BlockSize)
		rand.Read(masterkey)
	}
	return &Config{
		Port:        model.Port,
		OauthConfig: oauthConfig,
		MasterKey:   masterkey,
	}, nil
}

func getOauthConfig(credpath string) (*oauth2.Config, error) {
	binCreds, err := os.ReadFile(credpath)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(binCreds, drive.DriveFileScope)
}
