// +build !windows !nix

package uptime

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"time"
	"unsafe"
)

// Uptime returns the number of seconds since the system has booted
func Uptime() (int64, error) {
	tv := syscall.Timeval32{}
	if err := sysctlbyname("kern.boottime", &tv); err != nil {
		return 0, err
	}

	newT := time.Since(time.Unix(int64(tv.Sec), int64(tv.Usec)*1000)).Seconds()

	return int64(newT), nil
}

// generic Sysctl buffer unmarshalling
func sysctlbyname(name string, data interface{}) (err error) {
	val, err := syscall.Sysctl(name)
	if err != nil {
		return err
	}

	buf := []byte(val)

	switch v := data.(type) {
	case *uint64:
		*v = *(*uint64)(unsafe.Pointer(&buf[0]))
		return
	}

	bbuf := bytes.NewBuffer([]byte(val))
	return binary.Read(bbuf, binary.LittleEndian, data)
}
