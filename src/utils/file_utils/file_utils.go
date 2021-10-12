package file_utils

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/superbkibbles/bookstore_utils-go/rest_errors"
	"github.com/superbkibbles/realestate_users-api/src/utils/crypto_utils"
)

// Save
// If user not saved
// Delete Pic

func DeleteFile(fileName string) {
	os.Remove(filepath.Join("datasources/images", filepath.Base(fileName)))
}

func SaveFile(header *multipart.FileHeader, file multipart.File) (string, rest_errors.RestErr) {
	// Check if file is Pic Or Video

	splitter := strings.Split(header.Filename, ".")
	fileName := crypto_utils.GetMd5(header.Filename+strconv.FormatInt(time.Now().Unix(), 36)) + "." + splitter[len(splitter)-1]
	_, err := os.Stat(filepath.Join("datasources/images", filepath.Base(fileName)))
	if os.IsNotExist(err) {
		out, err := os.Create(filepath.Join("datasources/images", filepath.Base(fileName)))
		if err != nil {
			return "", rest_errors.NewInternalServerErr("Error while saving Pic", nil)
		}

		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return "", rest_errors.NewInternalServerErr("Error while saving Pic", nil)
		}
		return fileName, nil
	} else {
		return "", rest_errors.NewRestError("File Already exist", http.StatusAlreadyReported, "Already exist", nil)
	}
}

func UpdateFile(header *multipart.FileHeader, file multipart.File, path string) (string, rest_errors.RestErr) {
	splittedPath := strings.Split(path, "/")
	fileName := splittedPath[len(splittedPath)-1]
	DeleteFile(fileName)
	newFileName, err := SaveFile(header, file)
	if err != nil {
		return "", err
	}
	return newFileName, nil
}
