package utils

import (
	"errors"
	"math/rand"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcInfo struct {
	Name string
	Pid  uint32
}

func TerminateProc(pid uint32) error {
	handle, err := windows.OpenProcess(
		windows.PROCESS_TERMINATE|windows.SYNCHRONIZE,
		false,
		pid,
	)

	if err != nil {
		return err
	}

	defer windows.CloseHandle(handle)

	return windows.TerminateProcess(handle, 0)
}

func GetProcNameFromPid(pid uint32) (string, error) {
	proc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}

	defer windows.CloseHandle(proc)

	buffer := make([]uint16, 1204)
	bufferLen := uint32(len(buffer))

	err = windows.QueryFullProcessImageName(proc, 0, &buffer[0], &bufferLen)
	if err != nil {
		return "", err
	}

	path := windows.UTF16ToString(buffer)

	return filepath.Base(path), nil
}

func GetAllRunningProcs() ([]ProcInfo, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	defer windows.CloseHandle(snap)

	var proc windows.ProcessEntry32
	proc.Size = uint32(unsafe.Sizeof(proc))

	err = windows.Process32First(snap, &proc)
	if err != nil {
		return nil, err
	}

	var procs []ProcInfo

	for {
		procs = append(procs, ProcInfo{
			Name: windows.UTF16ToString(proc.ExeFile[:]),
			Pid:  proc.ProcessID,
		})

		err := windows.Process32Next(snap, &proc)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}

			return nil, err
		}
	}

	if len(procs) == 0 {
		return nil, errors.New("failed to get any processes")
	}

	return procs, nil
}

func GetRandomRunningProc() (ProcInfo, error) {
	procs, err := GetAllRunningProcs()
	if err != nil {
		return ProcInfo{}, err
	}

	return procs[rand.Intn(len(procs))], nil
}
