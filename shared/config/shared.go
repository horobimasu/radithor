package config

const (
	_ITEM_TYPE_INVALID    uint8 = 0
	_ITEM_TYPE_U8         uint8 = 1
	_ITEM_TYPE_U16        uint8 = 2
	_ITEM_TYPE_U32        uint8 = 3
	_ITEM_TYPE_U64        uint8 = 4
	_ITEM_TYPE_BYTES      uint8 = 5
	_ITEM_TYPE_BOOL_GROUP uint8 = 6
)

const (
	_ITEM_SIZE_U8  uint = 1
	_ITEM_SIZE_U16 uint = 2
	_ITEM_SIZE_U32 uint = 4
	_ITEM_SIZE_U64 uint = 8
)

const _MAX_ITEM_SIZE_STRING = 0xFF_FF
