package main

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Url struct {
	Url string `xml:"url,attr"`
}

type Item struct {
	Title     string `xml:"title"`
	Enclosure Url    `xml:"enclosure"`
}

type List struct {
	XMLName xml.Name `xml:"rss"`
	Items   []Item   `xml:"channel>item"`
}

func main() {
	resp, err := http.Get("https://learningenglish.voanews.com/podcast/video.aspx/?zoneId=6042")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}
	var items List

	if err := xml.Unmarshal(data, &items); err != nil {
		log.Print(err)
		return
	}

	for _, item := range items.Items {
		fn := filepath.Join(item.Title, filepath.Ext(item.Enclosure.Url))
		fmt.Println(item.Title, item.Enclosure.Url, fn)
		if err := down(fn, item.Enclosure.Url); err != nil {
			log.Print(err)
		}
	}
}

func down(fn, url string) error {
	f, err := os.Create(sanitizePath(fn))
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func sanitizePath(path string) string {
	var b strings.Builder

	disallowed := `<>:"'/\|?*[]{};:!@$%&^#`
	prev := ' '
	for _, c := range path {
		if !unicode.IsPrint(c) || c == unicode.ReplacementChar ||
			strings.ContainsRune(disallowed, c) {
			c = ' '
		}

		if !(c == ' ' && prev == ' ') {
			b.WriteRune(c)
		}
		prev = c
	}

	path = strings.TrimSpace(b.String())
	return path
}
