package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func CalcSha256(data []byte) (string, error) {

	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	// fmt.Println("digest:", digest)
	return digest, nil
}

func DownloadPage(image Image, client *MdClient) error {

	res, err := client.Get(image.Url)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer res.Body.Close()

	ext := filepath.Ext(image.Url)

	filename := fmt.Sprintf("%s_p%04d%s", image.Path, image.Idx+1, ext)

	img := filepath.Base(image.Url)
	withoutExt := strings.Split(img, ".")[0]
	origHash := strings.Split(withoutExt, "-")[1]

	if _, err := os.Stat(filename); os.IsNotExist(err) {

		err := os.MkdirAll(filepath.Dir(filename), 0775)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

	} else {

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer file.Close()

		fileContents, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fileHash, err := CalcSha256(fileContents)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		if origHash == fileHash {
			return nil
		} else {
			err := os.Rename(filename, filename+".bak")
			if err != nil {
				return err
			}
			fmt.Println("Hash mismatch. Old file is renamed with .bak extension.")
		}

	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	var resbuf bytes.Buffer
	_, err = resbuf.ReadFrom(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	responseHash, err := CalcSha256(resbuf.Bytes())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if origHash == responseHash {

		_, err := io.Copy(file, &resbuf)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(filename)

		return nil
	}

	return errors.New("invalid response received")
}

func Worker(id int, jobs <-chan Image, results chan<- error, client *MdClient, wg *sync.WaitGroup) {
	defer wg.Done()

	delay := [3]int{2, 5, 10}

	for job := range jobs {
		var err error
		for retry := 0; retry < 3; retry++ {
			err = DownloadPage(job, client)
			if err == nil {
				break
			}
			time.Sleep(time.Second * time.Duration(delay[retry]))
		}
		results <- err
	}

}
