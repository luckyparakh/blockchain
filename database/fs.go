package database

import (
	"os"
	"path/filepath"
)

func initDataDirIfNotExists(dataDir string) error {
	if fileExists(getGenesisJsonFilePath(dataDir)) {
		return nil
	}
	if !dirExists(getDatabaseDirPath(dataDir)) {
		err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if err := writeGenesisToDisk(getGenesisJsonFilePath(dataDir)); err != nil {
		return err
	}
	if err := writeEmptyBlocksDbToDisk(getBlocksDbFilePath(dataDir)); err != nil {
		return err
	}
	return nil
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func dirExists(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func writeEmptyBlocksDbToDisk(path string) error {
	return os.WriteFile(path, []byte(""), os.ModePerm)
}
