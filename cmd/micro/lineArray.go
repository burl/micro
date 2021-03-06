package main

import (
	"bytes"
	"unicode/utf8"
)

func runeToByteIndex(n int, txt []byte) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
}

type LineArray struct {
	lines [][]byte
}

func NewLineArray(text []byte) *LineArray {
	la := new(LineArray)
	split := bytes.Split(text, []byte("\n"))
	la.lines = make([][]byte, len(split))
	for i := range split {
		la.lines[i] = make([]byte, len(split[i]))
		copy(la.lines[i], split[i])
	}

	return la
}

func (la *LineArray) String() string {
	return string(bytes.Join(la.lines, []byte("\n")))
}

func (la *LineArray) NewlineBelow(y int) {
	la.lines = append(la.lines, []byte(" "))
	copy(la.lines[y+2:], la.lines[y+1:])
	la.lines[y+1] = []byte("")
}

func (la *LineArray) insert(pos Loc, value []byte) {
	x, y := runeToByteIndex(pos.X, la.lines[pos.Y]), pos.Y
	// x, y := pos.x, pos.y
	for i := 0; i < len(value); i++ {
		if value[i] == '\n' {
			la.Split(Loc{x, y})
			x = 0
			y++
			continue
		}
		la.insertByte(Loc{x, y}, value[i])
		x++
	}
}

func (la *LineArray) insertByte(pos Loc, value byte) {
	la.lines[pos.Y] = append(la.lines[pos.Y], 0)
	copy(la.lines[pos.Y][pos.X+1:], la.lines[pos.Y][pos.X:])
	la.lines[pos.Y][pos.X] = value
}

func (la *LineArray) JoinLines(a, b int) {
	la.insert(Loc{len(la.lines[a]), a}, la.lines[b])
	la.DeleteLine(b)
}

func (la *LineArray) Split(pos Loc) {
	la.NewlineBelow(pos.Y)
	la.insert(Loc{0, pos.Y + 1}, la.lines[pos.Y][pos.X:])
	la.DeleteToEnd(Loc{pos.X, pos.Y})
}

func (la *LineArray) remove(start, end Loc) string {
	sub := la.Substr(start, end)
	startX := runeToByteIndex(start.X, la.lines[start.Y])
	endX := runeToByteIndex(end.X, la.lines[end.Y])
	if start.Y == end.Y {
		la.lines[start.Y] = append(la.lines[start.Y][:startX], la.lines[start.Y][endX:]...)
	} else {
		for i := start.Y + 1; i <= end.Y-1; i++ {
			la.DeleteLine(start.Y + 1)
		}
		la.DeleteToEnd(Loc{startX, start.Y})
		la.DeleteFromStart(Loc{endX - 1, start.Y + 1})
		la.JoinLines(start.Y, start.Y+1)
	}
	return sub
}

func (la *LineArray) DeleteToEnd(pos Loc) {
	la.lines[pos.Y] = la.lines[pos.Y][:pos.X]
}

func (la *LineArray) DeleteFromStart(pos Loc) {
	la.lines[pos.Y] = la.lines[pos.Y][pos.X+1:]
}

func (la *LineArray) DeleteLine(y int) {
	la.lines = la.lines[:y+copy(la.lines[y:], la.lines[y+1:])]
}

func (la *LineArray) DeleteByte(pos Loc) {
	la.lines[pos.Y] = la.lines[pos.Y][:pos.X+copy(la.lines[pos.Y][pos.X:], la.lines[pos.Y][pos.X+1:])]
}

func (la *LineArray) Substr(start, end Loc) string {
	startX := runeToByteIndex(start.X, la.lines[start.Y])
	endX := runeToByteIndex(end.X, la.lines[end.Y])
	if start.Y == end.Y {
		return string(la.lines[start.Y][startX:endX])
	}
	var str string
	str += string(la.lines[start.Y][startX:]) + "\n"
	for i := start.Y + 1; i <= end.Y-1; i++ {
		str += string(la.lines[i]) + "\n"
	}
	str += string(la.lines[end.Y][:endX])
	return str
}
