package main

import (
    "regexp"
    "testing"
)

func TestDownloads(t *testing.T) {
    request := ReleaseRequest{
        base_url:  "https://pypi.python.org/simple/",
        name:      "requests",
    }

    releases := get_releases(request)
    if len(releases) < 4 {
        t.Error("not finding anything for python-requests?")
    }

    normalized := normalize_versions(releases)

    pattern, _ := regexp.Compile(`^(\d+\.)*(\*|\d+)$`)
    for i := 0; i < len(normalized); i++ {
        if !pattern.Match([]byte(normalized[i].version)) {
            t.Errorf("%s failed to normalize version", normalized[i].version)
        }
    }
}

func TestDownloadReal(t *testing.T) {
    request := ReleaseRequest{
        base_url:  "https://pypi.python.org/simple/",
        name:      "requests",
        version:   "1.0.0",
    }
    PyPIGet(request)
    // Output: requests-1.0.0.tar.gz downloaded (335548 bytes)
}

func TestDownloadNotFound(t *testing.T) {
    request := ReleaseRequest{
        base_url:  "https://pypi.python.org/simple/",
        name:      "requests",
        extension: "whl",
        version:   "1.0.0",
    }
    PyPIGet(request)
    // Output: No releases found for requests=1.0.0 (whl)
}
