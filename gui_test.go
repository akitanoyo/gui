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

func TestFindWindow(t *testing.T) {
    wnd, err := FindWindow("", "lockmouse")
    if err != nil {
        t.Errorf("got find window...")
    }

    if wnd > 0 {
        t.Logf("found lockmouse window\n")
        cn := GetClassName(wnd)
        t.Logf("classname %s\n", cn)
        CloseWindow(wnd)
        t.Logf("close window\n")
    } else {
        t.Logf("not found lockmouse window\n")
    }
}

