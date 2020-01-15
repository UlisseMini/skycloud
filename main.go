package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/cheggaaa/pb"
)

const (
	userid     = "" // USERID HERE
	username   = "" // USERNAME HERE
	postURL    = "https://skycloud.pro/media/UploadToServer.php?userid=" + userid
	genfilekey = "https://my.skycloud.pro/MYCLOUD/generatefileKey.php?userid=" + userid + "&filename="
)

// multipart/form-data
func upload(path string, finfo os.FileInfo) error {
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
	req, err := http.NewRequest("POST", postURL, barReader)
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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = os.Stderr.Write(bytes.TrimSpace(b))
	return err
}

func main() {
	if userid == "" || username == "" {
		panic("invalid constants, recompile source with valid constants")
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	finfo, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := upload(os.Args[1], finfo); err != nil {
		fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("https://media.skycloud.pro/%s%s/SKYCLOUD/Files/%s",
		userid, username, url.PathEscape(finfo.Name()))
	fmt.Fprintf(os.Stderr, "\n\nUploaded: %s to %s\n", os.Args[1], url)
}
