package config

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

type boolConfig struct {
	Id    uint8
	Value bool
}

func EncodeConfigs(values map[string]any, ids map[string]uint8, ignores []string, printTypes bool) ([]byte, error) {
	ignore := make(map[string]bool, len(ignores))
	for _, v := range ignores {
		ignore[v] = true
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var buffer bytes.Buffer
	var bools []boolConfig
	var ignored []string

	doneIds := make(map[uint8]string)

	for _, key := range keys {
		if ignore[key] {
			ignored = append(ignored, key)
			continue
		}

		goName := formatName(key)

		id, found := ids[goName]
		if !found {
			return nil, fmt.Errorf("failed to find mapping for key %q", key)
		}

		done, seen := doneIds[id]
		if seen {
			return nil, fmt.Errorf("key %q and %q are mapped to the same Id (%d)", done, key, id)
		}

		doneIds[id] = key

		val := values[key]

		switch theVal := val.(type) {
		case bool:
			bools = append(bools, boolConfig{
				Id:    id,
				Value: theVal,
			})

			if printTypes {
				fmt.Printf("%s -> bool\n", goName)
			}

		case float64:
			if theVal < 0 {
				return nil, fmt.Errorf("key %q cannot be a negative number", key)
			}

			if theVal != math.Trunc(theVal) {
				return nil, fmt.Errorf("key %q needs to be a whole number", key)
			}

			writeConfigSmallestNumber(&buffer, id, uint(theVal), key, printTypes)

		case string:
			data := []byte(theVal)
			if len(data) > _MAX_ITEM_SIZE_STRING {
				return nil, fmt.Errorf("key %q exceeds %d bytes (%d)", key, _MAX_ITEM_SIZE_STRING, len(data))
			}

			buffer.WriteByte(id)
			buffer.WriteByte(_ITEM_TYPE_BYTES)
			binary.Write(&buffer, binary.LittleEndian, uint16(len(data)))
			buffer.Write(data)

			if printTypes {
				fmt.Printf("%s -> string\n", goName)
			}

		default:
			return nil, fmt.Errorf("key %q is an unsupported type (%T)", key, val)
		}
	}

	for len(bools) > 0 {
		boolsCount := len(bools)
		if boolsCount > 8 {
			boolsCount = 8
		}

		group := bools[:boolsCount]
		bools = bools[boolsCount:]

		buffer.WriteByte(_ITEM_TYPE_BOOL_GROUP)
		buffer.WriteByte(byte(len(group)))

		var bits uint8

		for i, config := range group {
			buffer.WriteByte(config.Id)

			if config.Value {
				bits |= 1 << uint(i)
			}
		}

		buffer.WriteByte(bits)
	}

	return buffer.Bytes(), nil
}

func formatName(name string) string {
	var result strings.Builder

	for _, char := range name {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(unicode.ToUpper(char))
		} else {
			result.WriteByte('_')
		}
	}

	return result.String()
}

func getSmallestNumberType(num uint) (byte, uint) { // type, size
	switch {
	case num <= 0xFF:
		return _ITEM_TYPE_U8, _ITEM_SIZE_U8

	case num <= 0xFF_FF:
		return _ITEM_TYPE_U16, _ITEM_SIZE_U16

	case num <= 0xFF_FF_FF_FF:
		return _ITEM_TYPE_U32, _ITEM_SIZE_U32

	default:
		return _ITEM_TYPE_U64, _ITEM_SIZE_U64
	}
}

func writeConfigSmallestNumber(buffer *bytes.Buffer, id uint8, num uint, key string, printTypes bool) {
	kind, _ := getSmallestNumberType(num)

	buffer.WriteByte(id)
	buffer.WriteByte(kind)

	switch kind {
	case _ITEM_TYPE_U8:
		buffer.WriteByte(byte(num))

		if printTypes {
			fmt.Printf("%s -> u8\n", key)
		}

	case _ITEM_TYPE_U16:
		binary.Write(buffer, binary.LittleEndian, uint16(num))

		if printTypes {
			fmt.Printf("%s -> u16\n", key)
		}

	case _ITEM_TYPE_U32:
		binary.Write(buffer, binary.LittleEndian, uint32(num))

		if printTypes {
			fmt.Printf("%s -> u32\n", key)
		}

	case _ITEM_TYPE_U64:
		binary.Write(buffer, binary.LittleEndian, num)

		if printTypes {
			fmt.Printf("%s -> u64\n", key)
		}
	}
}
