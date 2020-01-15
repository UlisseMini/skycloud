package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cheggaaa/pb"
)

const url = "https://skycloud.pro/media/UploadToServer.php?userid=USERID_HERE"

// multipart/form-data
func upload(path string) error {
	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	mw := multipart.NewWriter(buf)

	mwf, err := mw.CreateFormFile("uploaded_file", f.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(mwf, f)
	if err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}

	// Create progressbar
	bar := pb.Full.Start64(finfo.Size())
	barReader := bar.NewProxyReader(buf)
	defer bar.Finish()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, barReader)
	if err != nil {
		return err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", mw.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	if err := upload(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\n\nUploaded: %s\n", os.Args[1])
}
