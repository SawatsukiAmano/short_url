﻿package utils

import "encoding/binary"

type MapKY interface {
	string | int
}

// 获取map所有key
func GetMapAllKeys[T MapKY](m map[T]T) []T {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 获取map所有value
func GetMapAllValues[T MapKY](m map[T]T) []T {
	values := make([]T, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// 无符号int转byte数组
func UInt32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}
