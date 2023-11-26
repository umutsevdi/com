package index

import (
	"strings"
	"time"
)

type GEMTEXT int

const (
	GEMTEXT_TEXT = iota
	GEMTEXT_LINK
	GEMTEXT_LIST
	GEMTEXT_QUOTE
	GEMTEXT_H3
	GEMTEXT_H2
	GEMTEXT_H1
	GEMTEXT_PRE
)

type Gemtext struct {
	Title        string
	Author       string
	Created      time.Time
	LastModified time.Time
	Lines        *[]Line
}

type Line struct {
	Type   GEMTEXT
	Text   string
	Arg    string
	PreNum int
}

func parseGemini(v *FData) Gemtext {

	var g Gemtext
	linesRaw := strings.Split(string(v.Content), "\n")
	g.Created = v.Created
	g.LastModified = v.LastModified

	lines := make([]Line, 0, len(linesRaw))
	var pre bool = false
	var preLastHint string = ""
	for _, v := range linesRaw {
		if pre {
			if strings.HasPrefix(v, "```") {
				pre = false
				preLastHint = ""
				lines = append(lines, Line{Type: GEMTEXT_PRE, Arg: preLastHint, Text: "", PreNum: 2})
				continue
			}
			lines = append(lines, Line{Type: GEMTEXT_PRE, Arg: preLastHint, Text: v, PreNum: 1})
			continue
		}
		if strings.HasPrefix(v, "//") {
			g.parseHeader(v)
		} else if strings.HasPrefix(v, "```") {
			pre = true
			preLastHint = strings.Split(v, " ")[0][3:]
			lines = append(lines, Line{Type: GEMTEXT_PRE, Arg: preLastHint, Text: "", PreNum: 0})
		} else if strings.HasPrefix(v, ">") {
			lines = append(lines, parse(v[1:], ">"))
		} else if strings.HasPrefix(v, "*") {
			lines = append(lines, parse(v[1:], "*"))
		} else if strings.HasPrefix(v, "=>") {
			lines = append(lines, parseLink(v))
		} else if strings.HasPrefix(v, "#") {
			lines = append(lines, parseH(v))
		} else {
			lines = append(lines, Line{Text: strings.TrimSpace(v), Type: GEMTEXT_TEXT})
		}
	}
	for _, v := range lines {
		if v.Type != GEMTEXT_PRE {
			v.Text = strings.TrimSpace(v.Text)
			v.Arg = strings.TrimSpace(v.Arg)
		}
	}
	g.Lines = &lines

	return g
}

func (g *Gemtext) parseHeader(line string) {
	s := (line)[2:]
	sArray := strings.Split(s, ":")
	if len(sArray) != 2 {
		return
	}
	if strings.TrimSpace(sArray[0]) == "Title" {
		g.Title = strings.TrimSpace(sArray[1])
	} else if strings.TrimSpace(sArray[0]) == "Author" {
		g.Author = strings.TrimSpace(sArray[1])
	}
}

func (l *Line) IsImage() bool {
	s := strings.ToLower(ext(l.Arg)[1:])
	return strings.Compare(s, "png") == 0 || strings.Compare(s, "jpg") == 0 || strings.Compare(s, "jpeg") == 0
}

func (g *Gemtext) DateFmt(date time.Time) string {
	return date.Format("15:04 Mon, 2 Jan 06")
}

func parseH(line string) Line {
	l := Line{}
	jCount := 0
	for _, j := range line {
		if j == '#' {
			jCount++
		}
		if jCount > 2 {
			break
		}
	}
	switch jCount {
	case 1:
		l.Type = GEMTEXT_H1
	case 2:
		l.Type = GEMTEXT_H2
	case 3:
		l.Type = GEMTEXT_H3
	}
	l.Text = strings.TrimSpace(line[jCount:])
	return l
}

func parseLink(line string) Line {
	l := Line{Type: GEMTEXT_LINK}
	p := strings.Split(strings.TrimSpace(line[2:]), " ")
	if len(p) > 1 {
		l.Text = strings.Join(p[1:], " ")
		l.Arg = p[0]
	} else {
		l.Arg = p[0]
		l.Text = p[0]
	}
	return l
}

func parse(line, char string) Line {
	l := Line{}
	if char == "*" {
		l.Type = GEMTEXT_LIST
	} else {
		l.Type = GEMTEXT_QUOTE
	}
	l.Text = strings.TrimSpace(line[1:])
	return l
}
