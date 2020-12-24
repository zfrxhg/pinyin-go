package pinyin_test

import (
	"testing"

	"github.com/zfrxhg/pinyin-go/pinyin"
)

const (
	dictFile          = `..\cedict_1_0_ts_utf-8_mdbg.txt.gz`
	str               = "体重两重天"
	strPinyin         = "ti zhong liang chong tian"
	strPinyinInitials = "tzlct"
)

func Test_Pinyin(t *testing.T) {
	dict, err := pinyin.NewDictionary(dictFile)
	if err != nil {
		t.Fatal(err)
	}
	if strPinyin != dict.Pinyin(str) {
		t.FailNow()
	}
	if strPinyinInitials != dict.PinyinInitials(str) {
		t.FailNow()
	}
}

func Benchmark_Pinyin(b *testing.B) {
	b.StopTimer()
	dict, err := pinyin.NewDictionary(dictFile)
	if err != nil {
		b.Fatal(err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if strPinyin != dict.Pinyin(str) {
			b.FailNow()
		}
	}
}
