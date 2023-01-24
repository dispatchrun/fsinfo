//go:build darwin || linux

package fsinfo

import (
	"fmt"
	"io/fs"
	"syscall"
	"time"
)

func newFileInfo(name string, mode fs.FileMode, modTime time.Time, size int64, sys any) *fileInfo {
	info := &fileInfo{name: name}
	if s := sys.(*syscall.Stat_t); s != nil {
		info.stat = *s
	}
	if mode != 0 {
		setMode(&info.stat, makeMode(mode))
	}
	if !modTime.IsZero() {
		setMtime(&info.stat, modTime.Unix(), int64(modTime.Nanosecond()))
	}
	if size >= 0 {
		setSize(&info.stat, size)
	}
	return info
}

type fileInfo struct {
	name string
	stat syscall.Stat_t
}

func (info *fileInfo) String() string {
	mode := info.Mode()
	size := info.Size()
	mtime := info.ModTime().Format("Jan _2 15:04")
	nlink := statNlink(&info.stat)
	return fmt.Sprintf("%s %d %d %s %s", mode, nlink, size, mtime, info.name)
}

func (info *fileInfo) Name() string {
	return info.name
}

func (info *fileInfo) Mode() fs.FileMode {
	mode := fs.FileMode(info.stat.Mode).Perm()
	switch info.stat.Mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		mode |= fs.ModeDevice
	case syscall.S_IFCHR:
		mode |= fs.ModeDevice | fs.ModeCharDevice
	case syscall.S_IFDIR:
		mode |= fs.ModeDir
	case syscall.S_IFIFO:
		mode |= fs.ModeNamedPipe
	case syscall.S_IFLNK:
		mode |= fs.ModeSymlink
	case syscall.S_IFSOCK:
		mode |= fs.ModeSocket
	}
	if (info.stat.Mode & syscall.S_ISGID) != 0 {
		mode |= fs.ModeSetgid
	}
	if (info.stat.Mode & syscall.S_ISUID) != 0 {
		mode |= fs.ModeSetuid
	}
	if (info.stat.Mode & syscall.S_ISVTX) != 0 {
		mode |= fs.ModeSticky
	}
	return mode
}

func (info *fileInfo) ModTime() time.Time {
	return time.Unix(statMtime(&info.stat))
}

func (info *fileInfo) IsDir() bool {
	return (info.stat.Mode & syscall.S_IFDIR) != 0
}

func (info *fileInfo) Size() int64 {
	return info.stat.Size
}

func (info *fileInfo) Sys() any {
	return &info.stat
}

func makeMode(fileMode fs.FileMode) (mode uint32) {
	mode = uint32(fileMode.Perm())
	switch fileMode.Type() {
	case fs.ModeDevice:
		mode |= syscall.S_IFBLK
	case fs.ModeDevice | fs.ModeCharDevice:
		mode |= syscall.S_IFCHR
	case fs.ModeDir:
		mode |= syscall.S_IFDIR
	case fs.ModeNamedPipe:
		mode |= syscall.S_IFIFO
	case fs.ModeSymlink:
		mode |= syscall.S_IFLNK
	case fs.ModeSocket:
		mode |= syscall.S_IFSOCK
	default:
		mode |= syscall.S_IFREG
	}
	if (fileMode & fs.ModeSetgid) != 0 {
		mode |= syscall.S_ISGID
	}
	if (fileMode & fs.ModeSetuid) != 0 {
		mode |= syscall.S_ISUID
	}
	if (fileMode & fs.ModeSticky) != 0 {
		mode |= syscall.S_ISVTX
	}
	return mode
}

func mode(info fs.FileInfo) uint32 {
	if stat := stat(info); stat != nil {
		return statMode(stat)
	}
	return makeMode(info.Mode())
}

func uid(info fs.FileInfo) uint32 {
	if stat := stat(info); stat != nil {
		return statUid(stat)
	}
	return 0
}

func gid(info fs.FileInfo) uint32 {
	if stat := stat(info); stat != nil {
		return statGid(stat)
	}
	return 0
}

func ino(info fs.FileInfo) uint64 {
	if stat := stat(info); stat != nil {
		return statIno(stat)
	}
	return 0
}

func nlink(info fs.FileInfo) uint64 {
	if stat := stat(info); stat != nil {
		return statNlink(stat)
	}
	return 1
}

func device(info fs.FileInfo) uint64 {
	if stat := stat(info); stat != nil {
		return statDev(stat)
	}
	return 0
}

func mtime(info fs.FileInfo) time.Time {
	if stat := stat(info); stat != nil {
		return time.Unix(statMtime(stat))
	}
	return time.Time{}
}

func atime(info fs.FileInfo) time.Time {
	if stat := stat(info); stat != nil {
		return time.Unix(statAtime(stat))
	}
	return time.Time{}
}

func ctime(info fs.FileInfo) time.Time {
	if stat := stat(info); stat != nil {
		return time.Unix(statCtime(stat))
	}
	return time.Time{}
}

func stat(info fs.FileInfo) *syscall.Stat_t {
	stat, _ := info.Sys().(*syscall.Stat_t)
	return stat
}
