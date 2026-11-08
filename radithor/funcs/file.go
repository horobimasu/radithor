package funcs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"radithor/utils"

	"github.com/horobimasu/gotils"
	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

const _CP_WINUNICODE uint16 = 0x04B0

const _LANGID_NEUTRAL uint16 = 0x0000

const _BLOCK_ALIGNMENT = 4

const (
	_BLOCK_TYPE_TEXT   uint16 = 1
	_BLOCK_TYPE_BINARY uint16 = 0
)

const ( // struct VS_FIXEDFILEINFO
	_DW_SIGNATURE      uint32 = 0xFEEF04BD
	_DW_STRUCT_VERSION uint32 = 0x00010000
	_DW_FILE_OS        uint32 = 0x00040004
	_DW_FILE_TYPE      uint32 = 0x00000001
)

const (
	GRPICONDIR_LEN      = 6
	GRPICONDIRENTRY_LEN = 14
)

type image struct {
	Width  byte
	Height byte
	Colors byte
	Planes uint16
	Bits   uint16
	Data   []byte
}

func SetExeFileFileDescription(filePath string, description string) error {
	filePath = gotils.ExpandPath(filePath)

	handle, err := windowsx.BeginUpdateResource(
		windows.StringToUTF16Ptr(filePath),
		false,
	)

	if err != nil {
		return err
	}

	var updated bool

	defer func() {
		if updated {
			return
		}

		windowsx.EndUpdateResource(handle, true)
	}()

	data := createVersionInfoDetail(description, _LANGID_NEUTRAL)

	err = windowsx.UpdateResource(
		handle,
		(*uint16)(unsafe.Pointer(uintptr(windowsx.RT_VERSION))), // windows moment
		(*uint16)(unsafe.Pointer(uintptr(1))),
		_LANGID_NEUTRAL,
		data,
	)

	if err != nil {
		return err
	}

	err = windowsx.EndUpdateResource(handle, false)
	if err != nil {
		return err
	}

	updated = true
	return nil
}

func SetExeFileFileIcon(fileToSetPath string, fileToGetFromPath string) error {
	fileToSetPath = gotils.ExpandPath(fileToSetPath)
	fileToGetFromPath = gotils.ExpandPath(fileToGetFromPath)

	icons, err := GetExeFileFileIcons(fileToGetFromPath)
	if err != nil {
		return err
	}

	handle, err := windowsx.BeginUpdateResource(
		windows.StringToUTF16Ptr(fileToSetPath),
		false,
	)

	if err != nil {
		return err
	}

	var updated bool

	defer func() {
		if updated {
			return
		}

		windowsx.EndUpdateResource(handle, true)
	}()

	iconGroup := make([]byte, GRPICONDIR_LEN+(len(icons)*GRPICONDIRENTRY_LEN))
	binary.LittleEndian.PutUint16(iconGroup[2:], 1)
	binary.LittleEndian.PutUint16(iconGroup[4:], uint16(len(icons)))

	for i, icon := range icons {
		id := uint16(i + 1)

		err := windowsx.UpdateResource(
			handle,
			(*uint16)(unsafe.Pointer(uintptr(windowsx.RT_ICON))),
			(*uint16)(unsafe.Pointer(uintptr(id))),
			_LANGID_NEUTRAL,
			icon.Data,
		)

		if err != nil {
			return err
		}

		groupResource := iconGroup[GRPICONDIR_LEN+(i*GRPICONDIRENTRY_LEN):]
		groupResource[0] = icon.Width
		groupResource[1] = icon.Height
		groupResource[2] = icon.Colors

		binary.LittleEndian.PutUint16(groupResource[4:], icon.Planes)
		binary.LittleEndian.PutUint16(groupResource[6:], icon.Bits)
		binary.LittleEndian.PutUint32(groupResource[8:], uint32(len(icon.Data)))
		binary.LittleEndian.PutUint16(groupResource[12:], id)
	}

	err = windowsx.UpdateResource(
		handle,
		(*uint16)(unsafe.Pointer(uintptr(windowsx.RT_GROUP_ICON))),
		(*uint16)(unsafe.Pointer(uintptr(1))),
		_LANGID_NEUTRAL,
		iconGroup,
	)

	if err != nil {
		return err
	}

	err = windowsx.EndUpdateResource(handle, false)
	if err != nil {
		return err
	}

	updated = true
	return nil
}

func GetExeFileFileIcons(filePath string) ([]image, error) {
	filePath = gotils.ExpandPath(filePath)

	theHandle, err := windows.LoadLibraryEx(filePath, 0, windows.LOAD_LIBRARY_AS_DATAFILE)
	if err != nil {
		return nil, err
	}

	defer windows.FreeLibrary(theHandle)

	var iconGroup uintptr
	var found bool

	callback := syscall.NewCallback(func(
		theHandle uintptr,
		kind uintptr,
		name uintptr,
		param uintptr,
	) uintptr {
		iconGroup = name
		found = true

		return 0
	})

	handle := uintptr(theHandle)

	windowsx.EnumResourceNames(uintptr(handle), windowsx.RT_GROUP_ICON, callback, 0)

	if !found {
		return nil, errors.New("failed to find icon group")
	}

	data, err := readResource(handle, iconGroup, windowsx.RT_GROUP_ICON)
	if err != nil {
		return nil, err
	}

	if len(data) < GRPICONDIR_LEN {
		return nil, errors.New("icon data is smaller than header")
	}

	count := int(binary.LittleEndian.Uint16(data[4:]))
	if len(data) < GRPICONDIR_LEN+(count*GRPICONDIRENTRY_LEN) {
		return nil, errors.New("icon data too small")
	}

	images := make([]image, 0, count)

	for i := range count {
		entry := data[GRPICONDIR_LEN+(i*GRPICONDIRENTRY_LEN):]
		id := uintptr(binary.LittleEndian.Uint16(entry[12:])) // nId

		data, err := readResource(handle, id, windowsx.RT_ICON)
		if err != nil {
			return nil, err
		}

		images = append(images, image{
			Width:  entry[0],
			Height: entry[1],
			Colors: entry[2],
			Planes: binary.LittleEndian.Uint16(entry[4:]),
			Bits:   binary.LittleEndian.Uint16(entry[6:]),
			Data:   data,
		})
	}

	return images, nil
}

func readResource(handle uintptr, name uintptr, kind uint32) ([]byte, error) {
	resource, err := windowsx.FindResource(
		handle,
		(*uint16)(unsafe.Pointer(name)),
		kind,
	)

	if err != nil {
		return nil, err
	}

	size, err := windowsx.SizeOfResource(handle, resource)
	if err != nil {
		return nil, err
	}

	data, err := windowsx.LoadResource(handle, resource)
	if err != nil {
		return nil, err
	}

	resourcePtr, err := windowsx.LockResource(data)
	if err != nil {
		return nil, err
	}

	realData := make([]byte, size)
	copy(realData, unsafe.Slice((*byte)(unsafe.Pointer(resourcePtr)), size))

	return realData, nil
}

func createVersionInfoDetail(description string, language uint16) []byte {
	fileInfo := createInfoBlock("FileDescription", _BLOCK_TYPE_TEXT, utils.StringToUTF16LeBytes(description))
	fileName := createInfoBlock("ProductName", _BLOCK_TYPE_TEXT, utils.StringToUTF16LeBytes(description))

	tableKey := fmt.Sprintf("%04X%04X", language, _CP_WINUNICODE)
	stringTable := createInfoBlock(tableKey, _BLOCK_TYPE_TEXT, nil, fileInfo, fileName)
	stringFileInfo := createInfoBlock("StringFileInfo", _BLOCK_TYPE_TEXT, nil, stringTable)

	translation := make([]byte, 4)
	binary.LittleEndian.PutUint16(translation[0:], language)
	binary.LittleEndian.PutUint16(translation[2:], _CP_WINUNICODE)

	translationBlock := createInfoBlock("Translation", _BLOCK_TYPE_BINARY, translation)
	varFileInfo := createInfoBlock("VarFileInfo", _BLOCK_TYPE_TEXT, nil, translationBlock)

	fixedInfo := make([]byte, 52)
	binary.LittleEndian.PutUint32(fixedInfo[0:], _DW_SIGNATURE)
	binary.LittleEndian.PutUint32(fixedInfo[4:], _DW_STRUCT_VERSION)
	binary.LittleEndian.PutUint32(fixedInfo[40:], _DW_FILE_OS)
	binary.LittleEndian.PutUint32(fixedInfo[44:], _DW_FILE_TYPE)

	return createInfoBlock("VS_VERSION_INFO", _BLOCK_TYPE_BINARY, fixedInfo, stringFileInfo, varFileInfo)
}

func createInfoBlock(key string, kind uint16, value []byte, children ...[]byte) []byte {
	realKey := utils.StringToUTF16LeBytes(key)

	buffer := make([]byte, 6) // the header

	buffer = append(buffer, realKey...)
	buffer = utils.AlignBytes(buffer, _BLOCK_ALIGNMENT)

	buffer = append(buffer, value...)
	buffer = utils.AlignBytes(buffer, _BLOCK_ALIGNMENT)

	for _, child := range children {
		buffer = append(buffer, child...)
	}

	var valueLen uint16

	if kind == _BLOCK_TYPE_TEXT {
		valueLen = uint16(len(value) / 2)
	} else {
		valueLen = uint16(len(value))
	}

	bufferLen := uint16(len(buffer))

	// header data
	binary.LittleEndian.PutUint16(buffer[0:], bufferLen)
	binary.LittleEndian.PutUint16(buffer[2:], valueLen)
	binary.LittleEndian.PutUint16(buffer[4:], kind)

	return buffer
}
