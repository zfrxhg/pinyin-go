package main

import "C"

import (
	"unsafe"

	py "github.com/zfrxhg/pinyin-go/pinyin"
)

var dict *py.Dictionary

func main() {}

// LoadPinyinDictionary creates a pinyin dictionary and returns the handle of it.
//export LoadPinyinDictionary
func LoadPinyinDictionary(filename *C.char) {
	var err error
	dict, err = py.LoadDictionary(GoString(filename))
	if err != nil {
		panic(err)
	}
}

// GetPinyin looks up the word in the pinyin dictionary.
// If the word has been found, It puts the pinyin of the word into buf and returns the length of the pinyin.
// Otherwise it returns 0.
//export GetPinyin
func GetPinyin(buf *C.char, bufLen C.size_t, word *C.char) C.size_t {
	pinyin := dict.Pinyin(GoString(word))
	return MemoryCopy(unsafe.Pointer(buf), bufLen, []byte(pinyin))
}

// GetPinyinInitials looks up the word in the pinyin dictionary.
// If the word has been found, It puts the pinyin initials of the word into buf and returns the length of the pinyin initials.
// Otherwise it returns 0.
//export GetPinyinInitials
func GetPinyinInitials(buf *C.char, bufLen C.size_t, word *C.char) C.size_t {
	pinyin := dict.PinyinInitials(GoString(word))
	return MemoryCopy(unsafe.Pointer(buf), bufLen, []byte(pinyin))
}
