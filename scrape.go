package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func get_options(root *html.Node) []*html.Node {
	var result []*html.Node
	if root.FirstChild != nil {
		// log.Println("descending")
		result = get_options(root.FirstChild)
	}
	if root.DataAtom == atom.Option {
		result = append(result, root)
	}
	start := root.NextSibling
	for start != nil {
		if start.DataAtom == atom.Option {
			result = append(result, start)
		} else if start.Type == html.ElementNode {
			// log.Println(start.Data)
			result = append(result, get_options(start)...)
		}

		start = start.NextSibling
	}
	// log.Println("got", len(result), "results")
	return result
}

func get_value(option *html.Node) string {
	for _, a := range option.Attr {
		if a.Key == "value" {
			return a.Val
		}
	}
	return ""
}

func get_embed(root *html.Node) *html.Node {

	for root != nil {
		if root.Type == html.ElementNode && root.DataAtom == atom.Embed {
			// log.Println("got one")
			return root
		}
		if root.FirstChild != nil {
			// log.Println("trying kid")
			root = root.FirstChild
			continue
		}
		if root.NextSibling != nil {
			// log.Println("trying sibling")
			root = root.NextSibling
			continue
		}
		if root.Parent == nil {
			return nil
		}
		for {
			root = root.Parent
			if root == nil {
				return nil
			}
			if root.NextSibling != nil {
				root = root.NextSibling
				break
			}
		}

	}
	return nil
}

func download_file(cliaoke_dir, url string) error {
	_, file_name := path.Split(url)
	f, err := os.Create(filepath.Join(cliaoke_dir, file_name))
	if err != nil {
		return errors.New("Couldn't create file! " + err.Error())
	}

	defer f.Close() // double close -> error swallowed.

	resp, err := http.Get(url)
	if err != nil {
		return errors.New("Couldn't download " + url + ": " + err.Error())
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return errors.New("Couldn't download file " + file_name + ": " + err.Error())
	}

	err = f.Close()
	if err != nil {
		return errors.New("Error writing downloaded content " + file_name + ": " + err.Error())
	}

	fmt.Println("Downloaded", url)

	return nil
}

func scrape(cliaoke_dir string) error {
	base_uri := "http://www.albinoblacksheep.com/audio/midi/"
	response, err := http.Get(base_uri)
	if err != nil {
		return errors.New("Couldn't fetch base content: " + err.Error())
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return errors.New("Couldn't understand document body: " + err.Error())
	}

	for _, option := range get_options(doc) {
		slug := get_value(option)
		// log.Println(option)
		if slug != "" {
			var embed *html.Node
			page := base_uri + "/" + slug
			// log.Println(page)
			err := func() error {
				response, err := http.Get(page)
				if err != nil {
					return errors.New("Error fetching page " + page + ": " + err.Error())
				}
				defer response.Body.Close()
				slugdoc, err := html.Parse(response.Body)
				if err != nil {
					return errors.New("Error parsing page " + page + ": " + err.Error())
				}
				embed = get_embed(slugdoc)
				return nil
			}()
			if err != nil {
				return err
			}
			var file_url string
			for _, r := range embed.Attr {
				if r.Key == "src" {
					err = nil
					file_url = r.Val
					goto SKIP_ERR
				}
			}
			return errors.New("No src attribute found in embed in page " + page)

		SKIP_ERR:
			if err = download_file(cliaoke_dir, file_url); err != nil {
				return err
			}
		}
	}
	return nil
}
