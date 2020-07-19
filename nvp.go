package rcmd

import (
	"fmt"
	"strings"
)

// NvpAppend a name and value pair to the list as an Nvp object
func NvpAppend(list NvpList, name, value string) NvpList {
	list.Pairs = append(list.Pairs, Nvp{Name: strings.Trim(name, " "), Value: strings.Trim(value, " ")})
	return list
}

// NvpAppendPair append a string of name=value pair to the list as an Nvp object
func NvpAppendPair(list NvpList, nvp string) NvpList {
	b := strings.Split(nvp, "=")
	return NvpAppend(list, b[0], b[1])
}

// Get a value by name
func (list NvpList) Get(name string) (value string, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair.Value, true
		}
	}
	return "", false
}

// GetPair an nvp by name
func (list NvpList) GetPair(name string) (nvp Nvp, exists bool) {
	for _, pair := range list.Pairs {
		if name == pair.Name {
			return pair, true
		}
	}
	return Nvp{}, false
}

// NvpRemove by name
func NvpRemove(list NvpList, name string) NvpList {
	n := -1
	for i, pair := range list.Pairs {
		if name == pair.Name {
			n = i
			break
		}
	}
	if n >= 0 {
		list.Pairs = append(list.Pairs[:n], list.Pairs[n+1:]...)
	}
	return list
}

// NvpUpdate a value by name and tell whether there was a value to update
func NvpUpdate(list NvpList, name string, value string) (NvpList, bool) {
	n := -1
	for i, pair := range list.Pairs {
		if name == pair.Name {
			n = i
		}
	}

	if n >= 0 {
		list.Pairs[n].Value = value
		return list, true
	}
	return list, false
}

// GetString returns a string as name=value
func (nvp Nvp) GetString(name string) (value string) {
	return fmt.Sprintf("%s=%s", nvp.Name, nvp.Value)
}
