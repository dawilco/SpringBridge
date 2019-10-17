package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const numberImages = 6

var configuration struct {
	MySQLUser      string
	MySQLPassoword string
	MySQLHost      string
	MySQLDB        string
	MySQLQuery     string
}

func main() {
	//filename is the path to the json config file
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	db, _ := sql.Open("mysql", configuration.MySQLUser+":"+configuration.MySQLPassoword+"@("+configuration.MySQLHost+")/"+configuration.MySQLDB)

	rows, _ := db.Query(configuration.MySQLQuery)
	defer rows.Close()

	var wg sync.WaitGroup

	for rows.Next() {
		var ip string
		err := rows.Scan(&ip)
		if err != nil {
			log.Panicln(err)
		}
		wg.Add(1)
		go setImages(ip, &wg)
		time.Sleep(time.Millisecond * 50)
	}
	wg.Wait()
}

func setImages(ip string, thisWait *sync.WaitGroup) {
	fmt.Println("posting ", ip)
	client := &http.Client{}
	var files []string
	root := "img"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	files = slicePop(files, 0)

	var wg sync.WaitGroup

	for i := 0; i < numberImages; i++ {
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(len(files))
		fmt.Println(files[random])
		index := strconv.Itoa(i)
		wg.Add(1)
		go writeImage(client, "http://"+ip+"/api/config/splashbackground/"+index, files[random], &wg)
		files = slicePop(files, random)
	}
	wg.Wait()
	thisWait.Done()
}

func writeImage(client *http.Client, url string, filePath string, wg *sync.WaitGroup) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var requestBody bytes.Buffer

	multiPartWriter := multipart.NewWriter(&requestBody)

	fileName := strings.Split(filePath, "/")

	fieldWriter, err := multiPartWriter.CreateFormFile("file", fileName[1])
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fieldWriter, file)
	if err != nil {
		panic(err)
	}
	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		panic(err)
	}

	multiPartWriter.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("****ERROR ENCOUNTERED CONTINUING****", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)

	wg.Done()
}

func slicePop(slice []string, index int) []string {
	slice[index] = slice[len(slice)-1] // Copy last element to index i.
	slice[len(slice)-1] = ""           // Erase last element (write zero value).
	return slice[:len(slice)-1]        // Truncate slice.
}
