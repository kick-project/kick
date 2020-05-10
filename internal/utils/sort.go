package utils

import (
	"regexp"
	"sort"
	"strconv"
)

// matchVersion takes the version string splits into major, minor & patch integers
func matchVersion(version string) (major, minor, patch *int) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	match := re.FindStringSubmatch(version)
	if len(match) != 4 {
		return nil, nil, nil
	}
	maj, err := strconv.Atoi(match[1])
	if err != nil {
		panic(err)
	}
	min, err := strconv.Atoi(match[2])
	if err != nil {
		panic(err)
	}
	p, err := strconv.Atoi(match[3])
	if err != nil {
		panic(err)
	}
	return &maj, &min, &p
}

type ByVersion []string

func (v ByVersion) Len() int {
	return len(v)
}

func (v ByVersion) Swap(i, j int) {
	v[i], v[j] = v[j], v[i] 
}

func (v ByVersion) Less(i, j int) bool {
	return lessVersion(v[i], v[j])
}

func lessVersion(i, j string) bool {
	imajor, iminor, ipatch := matchVersion(i)
	jmajor, jminor, jpatch := matchVersion(j)

	if *imajor < *jmajor {
		return true
	} else if *imajor > *jmajor {
		return false
	}
	
	if *iminor < *jminor {
		return true
	} else if *iminor > *jminor {
		return false
	}

	if *ipatch < *jpatch {
		return true
	}

	return false
}

func LatestVersion (versions ...string) (string) {
	sort.Sort(sort.Reverse(ByVersion(versions)))
	return versions[0]
}