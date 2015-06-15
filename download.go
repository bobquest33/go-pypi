package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mcuadros/go-version"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ReleaseRequest struct {
	base_url string
	name     string
	version  string
}

type ReleaseRequests []ReleaseRequest

type Release struct {
	version   string
	url       string
	file_name string // file name to save as (from url)
}

type Releases []Release

func (slice Releases) Len() int {
	return len(slice)
}

func (slice Releases) Less(i, j int) bool {
	return version.Compare(slice[i].version, slice[j].version, "<")
}

func (slice Releases) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func PyPIGet(pkg ReleaseRequest, base_url string, extension string) {

	doc, err := goquery.NewDocument(base_url + pkg.name)
	if err != nil {
		log.Fatal(err)
	}

	releases := Releases{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		parsed_url, _ := s.Attr("href")
		if !strings.HasPrefix(parsed_url, "http") {
			// relative path. make it full again
			parsed_url = base_url + pkg.name + "/" + parsed_url
		}

		download_url, err := url.Parse(parsed_url)
		if err != nil {
			log.Fatal("could not parse %s url from %s", pkg.name, parsed_url)
		}
		download_url = download_url.ResolveReference(download_url)
		url_split := strings.Split(download_url.Path, "/")
		file_name := url_split[len(url_split)-1]

		releases = append(releases, Release{
			version:   s.Text(),
			url:       download_url.String(),
			file_name: file_name,
		})
	})

	for r := 0; r < len(releases); r++ {
		this_version := ""
		for c := len(pkg.name) + 1; c < len(releases[r].version); c++ {
			char := string(releases[r].version[c])
			if char == "." {
				this_version += char
			} else if _, err := strconv.Atoi(char); err == nil {
				this_version += char
			} else {
				break
			}
		}
		if last := len(this_version) - 1; last >= 0 && string(this_version[last]) == "." {
			this_version = this_version[:last]
		}
		releases[r].version = version.Normalize(this_version)
	}

	if len(releases) == 0 {
		fmt.Printf("No releases found for %s\n", pkg.name)
		return
	}

	sort.Sort(sort.Reverse(releases))

	var latest_version string
	if pkg.version != "" {
		latest_version = pkg.version
	} else {
		latest_version = releases[0].version
	}

	for r := 0; r < len(releases); r++ {
		if version.CompareSimple(latest_version, releases[r].version) == 0 {
			if extension == "" || strings.HasSuffix(releases[r].file_name, extension) {
				fmt.Printf(
					"%s downloaded (%d bytes)\n",
					releases[r].file_name,
					download_release(releases[r]),
				)
			}
		}
	}
}

func download_release(release Release) int64 {
	out, err := os.Create(release.file_name)
	defer out.Close()

	if err != nil {
		log.Fatal("error creating %s: %s", release.file_name, err)
	}

	resp, err := http.Get(release.url)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("error downloading %s: %s", release.url, err)
	}

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		log.Fatal("error writing to %s: %s", release.file_name, err)
	}

	return n
}

func main() {

	var urlPtr = flag.String(
		"url",
		"https://pypi.python.org/simple/",
		"The PyPI URL to use",
	)
	var extPtr = flag.String(
		"extension",
		"",
		"Only download files with the extension given",
	)

	flag.Parse()
	extension := *extPtr
	base_url := *urlPtr

	if !strings.HasPrefix(base_url, "http") {
		log.Fatal("url must start with http!")
	}

	if !strings.HasSuffix(base_url, "/") {
		base_url += "/"
	}

	requests := ReleaseRequests{}
	var args = flag.Args()
	for i := 0; i < len(args); i++ {
		if strings.Index(args[i], "=") > -1 && strings.Count(args[i], "=") == 1 {
			reqPkg := strings.Split(args[i], "=")[0]
			reqVer := strings.Split(args[i], "=")[1]
			requests = append(requests, ReleaseRequest{name: reqPkg, version: reqVer})
		} else {
			requests = append(requests, ReleaseRequest{name: args[i]})
		}
	}

	for r := 0; r < len(requests); r++ {
		PyPIGet(requests[r], base_url, extension)
	}
}
