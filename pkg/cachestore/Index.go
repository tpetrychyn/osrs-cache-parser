package cachestore

import "sort"

type Index struct {
	Id          int
	Procotol    int8
	Named       bool
	Revision    int32
	Crc         uint32
	Compression int8
	Groups      map[uint16]*Group
}

func (i *Index) GetGroupsAsArray() []*Group {
	keys := make([]int, 0)
	for k := range i.Groups {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	groupArray := make([]*Group, 0)
	for _, k := range keys {
		groupArray = append(groupArray, i.Groups[uint16(k)])
	}

	return groupArray
}
