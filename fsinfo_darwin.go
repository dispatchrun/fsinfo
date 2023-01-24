package fsinfo

import "syscall"

func setMode(stat *syscall.Stat_t, mode uint32) {
	stat.Mode = uint16(mode)
}

func setMtime(stat *syscall.Stat_t, sec, nsec int64) {
	stat.Mtimespec = syscall.Timespec{Sec: sec, Nsec: nsec}
}

func setSize(stat *syscall.Stat_t, size int64) {
	stat.Size = size
}

func statMode(stat *syscall.Stat_t) uint32 {
	return uint32(stat.Mode)
}

func statUid(stat *syscall.Stat_t) uint32 {
	return stat.Uid
}

func statGid(stat *syscall.Stat_t) uint32 {
	return stat.Gid
}

func statIno(stat *syscall.Stat_t) uint64 {
	return stat.Ino
}

func statNlink(stat *syscall.Stat_t) uint64 {
	return uint64(stat.Nlink)
}

func statDev(stat *syscall.Stat_t) uint64 {
	return uint64(stat.Dev)
}

func statMtime(stat *syscall.Stat_t) (int64, int64) {
	return stat.Mtimespec.Unix()
}

func statAtime(stat *syscall.Stat_t) (int64, int64) {
	return stat.Atimespec.Unix()
}

func statCtime(stat *syscall.Stat_t) (int64, int64) {
	return stat.Ctimespec.Unix()
}
