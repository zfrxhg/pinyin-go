package pinyin

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"sync"
	"unicode"
)

// Syllable represents a syllable in pinyin.
type Syllable struct {
	RawValue []rune
	Value    []rune
	Tone     rune
}

// Entry represents an entry in pinyin dictionary.
type Entry struct {
	Traditional string
	Simplified  string
	Syllables   []*Syllable
}

// Dictionary represents pinyin dictionary.
type Dictionary struct {
	EntryMap      *sync.Map
	EntryCount    int
	WordMaxLength int
}

func nextToken(array []rune) (element, remaining []rune) {
	if len(array) > 0 && array[0] == '[' {
		for i, r := range array {
			if r == ']' {
				return array[1:i], array[i+1:]
			}
		}
		return array[1:], nil
	}
	for i, r := range array {
		if r == ' ' {
			if i == 0 {
				return nextToken(array[1:])
			}
			return array[:i], array[i+1:]
		} else if r == '[' {
			return array[:i], array[i:]
		}
	}
	return array, nil
}

// LoadDictionary loads dictionary from file.
func LoadDictionary(filename string) (dict *Dictionary, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer gzfile.Close()
	dict = &Dictionary{
		EntryMap:      &sync.Map{},
		EntryCount:    0,
		WordMaxLength: 0,
	}
	scanner := bufio.NewScanner(gzfile)
	for scanner.Scan() {
		// The basic format of a CC-CEDICT entry is:
		// Traditional Simplified [pin1 yin1] /English equivalent 1/equivalent 2/
		// 中國 中国 [Zhong1 guo2] /China/Middle Kingdom/
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		array := []rune(line)
		var element []rune
		element, array = nextToken(array)
		if len(element) == 0 || element[0] == '#' {
			continue
		}
		if len(array) == 0 {
			err = fmt.Errorf("invalid entry: %s", line)
			return
		}
		entry := &Entry{
			Traditional: string(element),
		}
		element, array = nextToken(array)
		if len(element) == 0 || len(array) == 0 {
			err = fmt.Errorf("invalid entry: %s", line)
			return
		}
		entry.Simplified = string(element)
		element, array = nextToken(array)
		if len(element) == 0 {
			err = fmt.Errorf("invalid entry: %s", line)
			return
		}
		entry.Syllables = parseSyllables(element)
		dict.EntryCount++
		var entries interface{}
		entries, _ = dict.EntryMap.LoadOrStore(entry.Traditional, []*Entry{})
		dict.EntryMap.Store(entry.Traditional, append(entries.([]*Entry), entry))
		entries, _ = dict.EntryMap.LoadOrStore(entry.Simplified, []*Entry{})
		dict.EntryMap.Store(entry.Simplified, append(entries.([]*Entry), entry))
		if len(entry.Traditional) > dict.WordMaxLength {
			dict.WordMaxLength = len(entry.Traditional)
		}
		if len(entry.Simplified) > dict.WordMaxLength {
			dict.WordMaxLength = len(entry.Simplified)
		}
	}
	fmt.Printf("%d entries are loaded\n", dict.EntryCount)
	return
}

func parseSyllables(array []rune) (syllables []*Syllable) {
	var element []rune
	for {
		element, array = nextToken(array)
		if len(element) > 0 {
			syllable := &Syllable{
				RawValue: element,
			}
			last := element[len(element)-1]
			if unicode.IsNumber(last) {
				syllable.Tone = last
				syllable.Value = element[:len(element)-1]
			} else {
				syllable.Value = element
			}
			syllables = append(syllables, syllable)
		}
		if array == nil {
			break
		}
	}
	return
}

func (p *Dictionary) pinyinPartial(word []rune) (entries []*Entry, r rune, remaining []rune) {
	var partial []rune
	if len(word) <= p.WordMaxLength {
		partial = word
	} else {
		partial = word[:p.WordMaxLength]
	}
	for len(partial) > 0 {
		value, ok := p.EntryMap.Load(string(partial))
		if ok {
			entries = value.([]*Entry)
			remaining = word[len(partial):]
			return
		}
		partial = partial[:len(partial)-1]
	}
	r = word[0]
	remaining = word[1:]
	return
}

func (p *Dictionary) pinyin(word []rune, formatSyllable func(*Syllable) string, formatRune func(r rune) string, separator string) string {
	buf := &bytes.Buffer{}
	var (
		entries []*Entry
		r       rune
	)
	for len(word) > 0 {
		entries, r, word = p.pinyinPartial(word)
		if len(entries) > 0 {
			for _, syllable := range entries[0].Syllables {
				str := formatSyllable(syllable)
				if len(str) > 0 {
					if len(separator) > 0 && buf.Len() > 0 {
						buf.WriteString(separator)
					}
					buf.WriteString(str)
				}
			}
		} else {
			str := formatRune(r)
			if len(str) > 0 {
				if len(separator) > 0 && buf.Len() > 0 {
					buf.WriteString(separator)
				}
				buf.WriteString(str)
			}
		}
	}
	return buf.String()
}

// Pinyin looks up the word and returns the pinyin of the word.
func (p *Dictionary) Pinyin(word string) string {
	return p.pinyin([]rune(word),
		func(syllable *Syllable) string {
			return string(syllable.RawValue)
		},
		func(r rune) string {
			if unicode.IsSpace(r) {
				return ""
			}
			return string([]rune{r})
		},
		" ")
}

func (p *Dictionary) formatInitial(r rune) string {
	if unicode.IsLower(r) || unicode.IsNumber(r) {
		return string([]rune{r})
	}
	if unicode.IsUpper(r) {
		return string([]rune{r - 'A' + 'a'})
	}
	return ""
}

// PinyinInitials looks up the word and returns the pinyin initials of the word.
func (p *Dictionary) PinyinInitials(word string) string {
	return p.pinyin([]rune(word),
		func(syllable *Syllable) string {
			r := []rune(syllable.Value)[0]
			return p.formatInitial(r)
		},
		p.formatInitial,
		"")
}
