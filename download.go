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
	base_url 	string
	extension 	string
	name     	string
	version  	string
}

type ReleaseRequests []ReleaseRequest

type Release struct {
	name      string
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

func get_releases(request ReleaseRequest) Releases {
	doc, err := goquery.NewDocument(request.base_url + request.name)
	if err != nil {
		log.Fatal(err)
	}

	releases := Releases{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		parsed_url, _ := s.Attr("href")
		if !strings.HasPrefix(parsed_url, "http") {
			// relative path. make it full again
			parsed_url = request.base_url + request.name + "/" + parsed_url
		}

		download_url, err := url.Parse(parsed_url)
		if err != nil {
			log.Fatal("could not parse download url from %s", parsed_url)
		}

		download_url = download_url.ResolveReference(download_url)
		url_split := strings.Split(download_url.Path, "/")

		releases = append(releases, Release{
			name:	   request.name,
			version:   s.Text(),
			url:       download_url.String(),
			file_name: url_split[len(url_split)-1],
		})
	})
	return releases
}

func normalize_versions(releases Releases) Releases {
	for i := 0; i < len(releases); i++ {
		this_version := ""
		for c := len(releases[i].name) + 1; c < len(releases[i].version); c++ {
			char := string(releases[i].version[c])
			if char == "." {
				this_version += char
			} else if _, err := strconv.Atoi(char); err == nil {
				this_version += char
			} else {
				break
			}
		}
		releases[i].version = version.Normalize(strings.TrimSuffix(this_version, "."))
	}
	return releases
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

func get_user_requests() ReleaseRequests {
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
		if strings.Count(args[i], "=") == 1 {
			reqPkg := strings.Split(args[i], "=")[0]
			reqVer := strings.Split(args[i], "=")[1]
			requests = append(requests, ReleaseRequest{
				name: reqPkg,
				version: reqVer,
				base_url: base_url,
				extension: extension,
			})
		} else {
			requests = append(requests, ReleaseRequest{
				name: args[i],
				base_url: base_url,
				extension: extension,
			})
		}
	}
	return requests
}

func PyPIGet(request ReleaseRequest) {
	releases := normalize_versions(get_releases(request))

	if len(releases) == 0 {
		fmt.Printf("No releases found for %s\n", request.name)
		return
	}

	sort.Sort(sort.Reverse(releases))

	var requested_version string
	if request.version == "" {
		requested_version = releases[0].version
	} else {
		requested_version = request.version
	}

	downloads := 0
	for i := 0; i < len(releases); i++ {
		if version.CompareSimple(requested_version, releases[i].version) == 0 {
			if request.extension == "" || strings.HasSuffix(releases[i].file_name, request.extension) {
				fmt.Printf(
					"%s downloaded (%d bytes)\n",
					releases[i].file_name,
					download_release(releases[i]),
				)
				downloads += 1
			}
		}
	}
	if downloads == 0 {
		err_txt := "No releases found for " + request.name
		if request.version != "" {
			err_txt += "=" + request.version
		}
		if request.extension != "" {
			err_txt += " (" + request.extension + ")"
		}
		fmt.Println(err_txt)
	}
}

func main() {
	requests := get_user_requests()
	for i := 0; i < len(requests); i++ {
		PyPIGet(requests[i])
	}
}
