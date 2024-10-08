package storage

import (
	"encoding/json"
	"os"
)

func SaveBlockchain(fileName string, blockchain interface{}) error {
    data, err := json.Marshal(blockchain)
    if err != nil {
        return err
    }
    return os.WriteFile(fileName, data, 0644)
}

func LoadBlockchain(fileName string, blockchain interface{}) error {
    data, err := os.ReadFile(fileName)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, blockchain)
}
