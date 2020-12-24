package pinyin

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Syllable represents a syllable in pinyin.
type Syllable struct {
	Value string
	Tone  int
}

// Entry represents an entry in pinyin dictionary.
type Entry struct {
	Traditional string
	Simplified  string
	Syllables   []*Syllable
}

// Dictionary represents pinyin dictionary.
type Dictionary struct {
	EntryMap *sync.Map
}

// NewDictionary creates a pinyin dictionary with specified dict file and returns it.
func NewDictionary(filename string) (dict *Dictionary, err error) {
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
		EntryMap: &sync.Map{},
	}
	scanner := bufio.NewScanner(gzfile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		index := strings.IndexRune(line, ' ')
		if index < 0 {
			continue
		}
		entry := &Entry{}
		entry.Traditional = line[:index]
		line = line[index:]
		index = strings.IndexRune(line, '[')
		if index < 0 {
			continue
		}
		entry.Simplified = strings.TrimSpace(line[:index])
		line = line[index+1:]
		index = strings.IndexRune(line, ']')
		if index < 0 {
			continue
		}
		pinyin := strings.TrimSpace(line[:index])
		syllables := strings.Split(pinyin, " ")
		for _, syllable := range syllables {
			syllableLen := len(syllable)
			if syllableLen == 0 {
				continue
			}
			var tone int
			tone, err = strconv.Atoi(syllable[syllableLen-1:])
			if err != nil {
				continue
			}
			entry.Syllables = append(entry.Syllables, &Syllable{
				Value: syllable[:syllableLen-1],
				Tone:  tone,
			})
		}
		var entries interface{}
		entries, _ = dict.EntryMap.LoadOrStore(entry.Traditional, []*Entry{})
		dict.EntryMap.Store(entry.Traditional, append(entries.([]*Entry), entry))
		entries, _ = dict.EntryMap.LoadOrStore(entry.Simplified, []*Entry{})
		dict.EntryMap.Store(entry.Simplified, append(entries.([]*Entry), entry))
	}
	return
}

func (p *Dictionary) pinyin(buf *bytes.Buffer, word []rune, join func(*bytes.Buffer, []*Syllable)) {
	for i := len(word); i >= 0; i-- {
		temp := word[:i]
		entries, ok := p.EntryMap.Load(string(temp))
		if ok {
			join(buf, entries.([]*Entry)[0].Syllables)
			word = word[i:]
			p.pinyin(buf, word, join)
			break
		}
	}
}

// Pinyin looks up the word and returns the pinyin of the word.
func (p *Dictionary) Pinyin(word string) string {
	buf := &bytes.Buffer{}
	p.pinyin(buf, []rune(word), joinSyllables)
	return buf.String()
}

// PinyinInitials looks up the word and returns the pinyin initials of the word.
func (p *Dictionary) PinyinInitials(word string) string {
	buf := &bytes.Buffer{}
	p.pinyin(buf, []rune(word), joinSyllableInitials)
	return buf.String()
}

func joinSyllableInitials(buf *bytes.Buffer, syllables []*Syllable) {
	for _, syllable := range syllables {
		buf.WriteRune(([]rune(syllable.Value))[0])
	}
	return
}

func joinSyllables(buf *bytes.Buffer, syllables []*Syllable) {
	for _, syllable := range syllables {
		if buf.Len() > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(syllable.Value)
	}
	return
}
