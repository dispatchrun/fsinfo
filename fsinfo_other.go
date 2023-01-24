//go:build !darwin && !linux

package fsinfo

import (
	"fmt"
	"io/fs"
	"time"
)

func newFileInfo(name string, mode fs.FileMode, modTime time.Time, size int64, sys any) *fileInfo {
	return &fileInfo{name, mode, modTime, size, sys}
}

type fileInfo struct {
	name    string
	mode    fs.FileMode
	modTime time.Time
	size    int64
	sys     any
}

func (info *fileInfo) String() string {
	mode := info.Mode()
	size := info.Size()
	mtime := info.ModTime().Format("Jan _2 15:04")
	return fmt.Sprintf("%s 1 %d %s %s", mode, size, mtime, info.name)
}

func (info *fileInfo) Name() string       { return info.name }
func (info *fileInfo) Mode() fs.FileMode  { return info.mode }
func (info *fileInfo) ModTime() time.Time { return info.modTime }
func (info *fileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info *fileInfo) Size() int64        { return info.size }
func (info *fileInfo) Sys() any           { return info.sys }

func makeMode(mode fs.FileMode) uint32 { return uint32(mode.Perm()) }

func mode(info fs.FileInfo) uint32 { return FileMode(info.Mode()) }

func uid(info fs.FileInfo) uint32 { return 0 }

func gid(info fs.FileInfo) uint32 { return 0 }

func ino(info fs.FileInfo) uint64 { return 0 }

func nlink(info fs.FileInfo) uint64 { return 1 }

func device(info fs.FileInfo) uint64 { return 0 }

func mtime(info fs.FileInfo) time.Time { return info.ModTime() }

func atime(info fs.FileInfo) time.Time { return time.Time{} }

func ctime(info fs.FileInfo) time.Time { return time.Time{} }
