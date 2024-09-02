package calendar

import (
  "bufio"
  "bytes"
  _ "embed"
  "fmt"
  "strings"
  "time"
)

//go:embed calendar.csv
var calendarData []byte

type Passage string

type Calendar [365][4]Passage

var MCheyne Calendar

func (c *Calendar) At(day int) (passages [4]Passage, ok bool) {
  if day < 0 || day > 365 {
    return
  }

  passages = c[day]
  ok = true
  return
}

func (c *Calendar) On(startedAt time.Time) (passage [4]Passage, ok bool) {
  daysSince := int(time.Since(startedAt) / (time.Hour*24))
  return c.At(int(daysSince))
}

func init() {
  scanner := bufio.NewScanner(bytes.NewReader(calendarData))
  idx := 0
  for scanner.Scan() {
    var parts []string
    txt := scanner.Text()
    if len(txt) == 0 {
      goto loopEnd
    }

    parts = strings.Split(scanner.Text(), "|")
    if len(parts) != 4 {
      panic(fmt.Errorf("Expected 4 parts in a row, got %d on row %d", len(parts), idx+1))
    }

    MCheyne[idx][0] = Passage(parts[0])
    MCheyne[idx][1] = Passage(parts[1])
    MCheyne[idx][2] = Passage(parts[2])
    MCheyne[idx][3] = Passage(parts[3])

    loopEnd:
    idx++
  }
}
