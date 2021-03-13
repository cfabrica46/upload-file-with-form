package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

func main() {

	http.HandleFunc("/", upload)

	fmt.Println("Listening on localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {

		log.Fatal(err)

	}

}

func upload(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		crutime := time.Now().Unix()

		h := md5.New()

		io.WriteString(h, strconv.FormatInt(crutime, 10))

		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("index.html")

		t.Execute(w, token)

	case "POST":

		r.ParseMultipartForm(4 << 20)

		file, handler, err := r.FormFile("uploadfile")

		if err != nil {
			w.Write([]byte("ERROR AL SUBIR ARCHIVO\n"))
			return
		}

		defer file.Close()

		fmt.Fprintln(w, handler.Header)
		fmt.Fprintln(w, handler.Filename)
		fmt.Fprintln(w, handler.Size)

		err = download(file, handler.Filename)

		if err != nil {
			if err == io.EOF {
				return
			}
			w.Write([]byte(err.Error()))
		}

	default:

		http.NotFound(w, r)

	}

}

func download(origin multipart.File, fileName string) (err error) {

	destino, err := os.Create(fileName)

	if err != nil {
		fmt.Println(2)
		return

	}

	defer destino.Close()

	b := make([]byte, 1024)

	var i int

	for {

		fmt.Println()
		_, err = origin.Read(b)

		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		fmt.Println(i)

		_, err = destino.Write(b)

		if err != nil {
			return
		}

		i++
	}

	return
}
