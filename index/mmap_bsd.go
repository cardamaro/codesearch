// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd openbsd netbsd

package index

import (
	"os"
	"syscall"

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
	n := int(size)
	if n == 0 {
		return mmapData{f, nil, nil}
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, (n+4095)&^4095, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		glog.Fatalf("mmap %s: %v", f.Name(), err)
	}
	return mmapData{f, data[:n], data}
}

func unmmapFile(m *mmapData) error {
	if err := syscall.Munmap(m.o); err != nil {
		return err
	}

	return m.f.Close()
}

func unmmap(d []byte) error {
	return syscall.Munmap(d)
}
