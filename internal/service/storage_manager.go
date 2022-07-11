package service

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"sort"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StorageManagerService struct {
	file_storage_path string
}

func NewStorageManagerService(file_storage_path string) *StorageManagerService {
	if stat, err := os.Stat(file_storage_path); err != nil && stat.IsDir() {
		//path is not a directory or directory not exists
		os.Mkdir(file_storage_path, 0777)
	}
	return &StorageManagerService{file_storage_path}
}

func (obj *StorageManagerService) CreateUnicFile(file_extension string) (*os.File, string, error) {
	unic_filename, err := obj.CreateUnicFileNameHashSum64()
	if err != nil {
		return &os.File{}, "", errors.Wrap(err, "cant create unic filename")
	}
	logrus.Info(unic_filename)
	file, file_path, err := obj.SaveFile(unic_filename, file_extension)
	if err != nil {
		return &os.File{}, "", errors.Wrap(err, "cant save file to local storage")
	}
	return file, file_path, nil

}

func (obj *StorageManagerService) CreateUnicFileNameHashSum64() (string, error) {
	filesInfo, err := ioutil.ReadDir(obj.file_storage_path)
	if err != nil {
		return "", errors.Wrap(err, "ioutil: cant ReadDir")
	}
	h := fnv.New64a()
	last_file_name := uuid.New().String()
	if len(filesInfo) != 0 {
		sort.Slice(filesInfo, func(i, j int) bool {
			return filesInfo[i].ModTime().Before(filesInfo[j].ModTime())
		})
		last_file_name = filesInfo[len(filesInfo)-1].Name()
		logrus.Info(last_file_name)
	}
	h.Write([]byte(last_file_name))
	result := fmt.Sprintf("%v", h.Sum64())
	return result, nil
}
func (obj *StorageManagerService) SaveFile(filename, file_extension string) (*os.File, string, error) {
	file_path := obj.file_storage_path + filename + file_extension
	file, err := os.Create(file_path)
	if err != nil {
		file.Close()
		return &os.File{}, file_path, errors.Wrap(err, "cant save file to local storage")
	}
	return file, file_path, nil
}

func (obj *StorageManagerService) GetFileStoragePath() string {
	return obj.file_storage_path
}
