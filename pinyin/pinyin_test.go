package pinyin_test

import (
	"fmt"
	"testing"

	"github.com/zfrxhg/pinyin-go/pinyin"
)

const (
	dictFile          = `..\cedict_1_0_ts_utf-8_mdbg.txt.gz`
	str               = "体重两重天ABC 珍·奥斯汀 人为财死，鸟为食亡"
	strPinyin         = "ti3 zhong4 liang3 chong2 tian1 A bi1 C Zhen1 · Ao4 si1 ting1 ren2 wei4 cai2 si3 , niao3 wei4 shi2 wang2"
	strPinyinInitials = "tzlctabczastrwcsnwsw"
)

func Test_Pinyin(t *testing.T) {
	dict, err := pinyin.LoadDictionary(dictFile)
	if err != nil {
		t.Fatal(err)
	}
	if strPinyin != dict.Pinyin(str) {
		fmt.Println(dict.Pinyin(str))
		t.FailNow()
	}
	if strPinyinInitials != dict.PinyinInitials(str) {
		fmt.Println(dict.PinyinInitials(str))
		t.FailNow()
	}
}

func Benchmark_Pinyin(b *testing.B) {
	b.StopTimer()
	dict, err := pinyin.LoadDictionary(dictFile)
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
