package gui

import (
    "testing"
    "sort"
)

var (
)

func init() {
}

func TestKeyList(t *testing.T) {
    list := KeyList()
    if len(list) == 0 {
       t.Errorf("got KeyList()\n")
    }
    sort.Slice(list, func(i, j int) bool {
        return list[i] < list[j]
    })
    for _, v := range list {
        t.Logf("%s\n", v)
    }
}

