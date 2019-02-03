package vlc

import (
	"bufio"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

type Version struct {
	version         *version.Version
	shutdownCommand string
}

type VersionFactory struct {
	versions        []*Version
	playerVersionRE *regexp.Regexp
}

func NewVersionFactory() *VersionFactory {
	vBase := &Version{version.Must(version.NewVersion("0.0.0")), "quit"}
	v4 := &Version{version.Must(version.NewVersion("4.0.0")), "shutdown"}

	versions := []*Version{vBase, v4}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].version.Compare(versions[j].version) > 0
	})

	return &VersionFactory{
		versions:        versions,
		playerVersionRE: regexp.MustCompile(`VLC media player (\w+.\w+.\w+)`),
	}
}

func (f *VersionFactory) Get(versionStr string) *Version {
	v, err := version.NewVersion(versionStr)
	if err != nil {
		return f.versions[len(f.versions)-1]
	}

	for _, vs := range f.versions {
		if v.Compare(vs.version) >= 0 {
			return vs
		}
	}

	return f.versions[len(f.versions)-1]
}

func (f *VersionFactory) Find(text string) *Version {
	s := bufio.NewScanner(strings.NewReader(text))
	for s.Scan() {
		line := s.Bytes()
		matches := f.playerVersionRE.FindSubmatch(line)
		if matches != nil {
			return f.Get(string(matches[1]))
		}
	}
	return f.Get("")
}
