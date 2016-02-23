// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/golang/glog"
)

func mmapFile(f *os.File) mmapData {
	st, err := f.Stat()
	if err != nil {
		glog.Fatal(err)
	}
	size := st.Size()
	if int64(int(size+4095)) != size+4095 {
		glog.Fatalf("%s: too large for mmap", f.Name())
	}
	if size == 0 {
		return mmapData{f, nil, nil}
	}
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		glog.Fatalf("CreateFileMapping %s: %v", f.Name(), err)
	}
	defer syscall.CloseHandle(syscall.Handle(h))

	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		glog.Fatalf("MapViewOfFile %s: %v", f.Name(), err)
	}

	data := (*[1 << 30]byte)(unsafe.Pointer(addr))
	return mmapData{f, data[:size], data[:]}
}

func unmmapFile(m *mmapData) error {
	err := syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&m.d[0])))
	if err != nil {
		return err
	}

	return m.f.Close()
}

func unmmap(d []byte) error {
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&d)))
}
