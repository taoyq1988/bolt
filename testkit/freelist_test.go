package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
	"unsafe"
)

type pgid uint64
type page struct {
	id       pgid
	flags    uint16
	count    uint16
	overflow uint32
	ptr      uintptr
}
const maxAllocSize = 0x7FFFFFFF

func TestFreelist_A(t *testing.T) {
	//buf := make([]byte, 16+8*12)
	//p := (*page)(unsafe.Pointer(&buf))
	p := &page{}
	fmt.Println(unsafe.Sizeof(*p), unsafe.Sizeof(p.ptr), unsafe.Sizeof(p.id))
	fmt.Println(int(unsafe.Offsetof(((*page)(nil)).ptr)))
	fmt.Println(p.id, p.flags, p.count, p.overflow)
	ids := (*[maxAllocSize]pgid)(unsafe.Pointer(&p.ptr))[:]
	fmt.Println(maxAllocSize)
	fmt.Println(len(ids))
	fmt.Println(cap(ids))
	fmt.Println(ids[:10])
	//fmt.Println(ids)
	//fmt.Println(ids[maxAllocSize/2])
	p.ptr = math.MaxUint64
	fmt.Println(p.ptr)
}

func TestFile(t *testing.T) {
	p := &page{}
	ptr := (*[maxAllocSize]byte)(unsafe.Pointer(p))
	sz := 4096
	offset := int64(10086)
	buf := ptr[:sz]
	f, err := os.Create("test.bin")
	defer os.Remove("test.bin")
	if err != nil {
		t.Fatal(err)
	}
	writeAt := f.WriteAt
	t.Log(len(buf), buf)
	if n, err := writeAt(buf, offset); err != nil {
		t.Fatal()
	} else {
		t.Log(n)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(b), b)
}
