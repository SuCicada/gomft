package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

func GetLogicalDrives() []string {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _, _ := GetLogicalDrives.Call()
	s := strconv.FormatInt(int64(n), 2)
	fmt.Println(s)
	var res []string
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '1' {
			driver := string(rune('A'+len(s)-i-1)) + ":"
			res = append(res, driver)
		}
	}
	return res
}
func GetVolumeInformation() {
	// https://bbs.csdn.net/topics/350167157
	// https://stackoverflow.com/questions/35238933/golang-call-getvolumeinformation-winapi-function
	var RootPathName = `C:\`
	var VolumeNameBuffer = make([]uint16, syscall.MAX_PATH+1)
	var nVolumeNameSize = uint32(len(VolumeNameBuffer))
	var VolumeSerialNumber uint32
	var MaximumComponentLength uint32
	var FileSystemFlags uint32
	var FileSystemNameBuffer = make([]uint16, 255)
	var nFileSystemNameSize uint32 = syscall.MAX_PATH + 1

	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getVolume, _ := syscall.GetProcAddress(kernel32, "GetVolumeInformationW")

	var nargs uintptr = 8
	p, _ := syscall.UTF16PtrFromString(RootPathName)
	ret, _, callErr := syscall.Syscall9(getVolume,
		nargs,
		uintptr(unsafe.Pointer(p)),
		uintptr(unsafe.Pointer(&VolumeNameBuffer[0])),
		uintptr(nVolumeNameSize),
		uintptr(unsafe.Pointer(&VolumeSerialNumber)),
		uintptr(unsafe.Pointer(&MaximumComponentLength)),
		uintptr(unsafe.Pointer(&FileSystemFlags)),
		uintptr(unsafe.Pointer(&FileSystemNameBuffer[0])),
		uintptr(nFileSystemNameSize),
		0)
	fmt.Println(ret, callErr)
	fmt.Println(syscall.UTF16ToString(VolumeNameBuffer))
	fmt.Println(syscall.UTF16ToString(FileSystemNameBuffer))
}
func main() {
	//fmt.Println(GetLogicalDrives())
	GetVolumeInformation()
}
