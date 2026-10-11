package config

import "encoding/binary"

func GetConfigAnyRaw[T any](rawConfig []byte, id uint8) T {
	var thing T

	data, found := GetConfigRaw(rawConfig, id)
	if !found {
		return thing
	}

	value, _ := data.(T)
	return value
}

func GetConfigStringRaw(rawConfig []byte, id uint8) string {
	data, found := GetConfigRaw(rawConfig, id)
	if !found {
		return ""
	}

	bytes, correct := data.([]byte)
	if !correct {
		return ""
	}

	return string(bytes)
}

func GetConfigRaw(rawConfig []byte, id uint8) (any, bool) {
	kind, data, found := findConfigId(rawConfig, id)
	if !found {
		return nil, false
	}

	switch kind {
	case _ITEM_TYPE_U8:
		return data[0], true

	case _ITEM_TYPE_U16:
		return binary.LittleEndian.Uint16(data), true

	case _ITEM_TYPE_U32:
		return binary.LittleEndian.Uint32(data), true

	case _ITEM_TYPE_U64:
		return binary.LittleEndian.Uint64(data), true

	case _ITEM_TYPE_BYTES:
		return data, true

	case _ITEM_TYPE_BOOL_GROUP:
		return data[0] != 0, true
	}

	return nil, false
}

func findConfigId(rawConfig []byte, target uint8) (byte, []byte, bool) { // type, data, found
	buffer := rawConfig

	for len(buffer) > 0 {
		if buffer[0] == _ITEM_TYPE_BOOL_GROUP {
			if len(buffer) < 2 {
				return _ITEM_TYPE_INVALID, nil, false
			}

			count := uint(buffer[1])
			if 2+count+1 > uint(len(buffer)) {
				return _ITEM_TYPE_INVALID, nil, false
			}

			ids := buffer[2 : 2+count]
			bits := buffer[2+count]
			len := 2 + count + 1

			for i, id := range ids {
				if id != target {
					continue
				}

				demasked := []byte{(bits >> uint(i) & 1)}
				return _ITEM_TYPE_BOOL_GROUP, demasked, true
			}

			buffer = buffer[len:]
			continue
		}

		if len(buffer) < 2 {
			return _ITEM_TYPE_INVALID, nil, false
		}

		id := buffer[0]
		kind := buffer[1]
		data := buffer[2:]

		var size uint

		switch kind {
		case _ITEM_TYPE_U8:
			size = _ITEM_SIZE_U8

		case _ITEM_TYPE_U16:
			size = _ITEM_SIZE_U16

		case _ITEM_TYPE_U32:
			size = _ITEM_SIZE_U32

		case _ITEM_TYPE_U64:
			size = _ITEM_SIZE_U64

		case _ITEM_TYPE_BYTES:
			if len(data) < 2 {
				return _ITEM_TYPE_INVALID, nil, false
			}

			len := binary.LittleEndian.Uint16(data)
			size = 2 + uint(len)

		default:
			return _ITEM_TYPE_INVALID, nil, false
		}

		if uint(len(data)) < size {
			return _ITEM_TYPE_INVALID, nil, false
		}

		if id == target {
			if kind == _ITEM_TYPE_BYTES {
				len := binary.LittleEndian.Uint16(data)
				return _ITEM_TYPE_BYTES, data[2 : 2+len], true
			}

			return kind, data[:size], true
		}

		buffer = data[size:]
	}

	return _ITEM_TYPE_INVALID, nil, false
}
