package controller

import (
	"strings"
	"time"
)

var weekdayAliases = map[time.Weekday][]string{
	time.Sunday:    {"domingo", "dom"},
	time.Monday:    {"segunda", "segunda feira", "segunda-feira", "seg"},
	time.Tuesday:   {"terca", "terca feira", "terca-feira", "terça", "terça feira", "terça-feira", "ter"},
	time.Wednesday: {"quarta", "quarta feira", "quarta-feira", "qua"},
	time.Thursday:  {"quinta", "quinta feira", "quinta-feira", "qui"},
	time.Friday:    {"sexta", "sexta feira", "sexta-feira", "sex"},
	time.Saturday:  {"sabado", "sabado feira", "sábado", "sábado feira", "sab"},
}

func parseRecordsDayQuery(dateStr string, now time.Time) (time.Time, error) {
	loc := resolveQueryLocation()
	return parseRecordsDayQueryWithLocation(dateStr, now, loc)
}

func parseRecordsDayQueryWithLocation(dateStr string, now time.Time, loc *time.Location) (time.Time, error) {
	raw := strings.TrimSpace(dateStr)
	nowInLoc := now.In(loc)
	today := startOfDay(nowInLoc)

	if raw == "" {
		return today, nil
	}

	if parsed, ok := parseAbsoluteDate(raw, loc); ok {
		return startOfDay(parsed.In(loc)), nil
	}

	norm := normalizeDateExpression(raw)
	switch norm {
	case "hoje", "today":
		return today, nil
	case "ontem", "yesterday":
		return today.AddDate(0, 0, -1), nil
	case "amanha", "amanhã", "tomorrow":
		return today.AddDate(0, 0, 1), nil
	}

	if weekday, ok := extractWeekday(norm); ok {
		forcePrevious := strings.Contains(norm, "passad") || strings.Contains(norm, "ultim")
		return resolveWeekday(nowInLoc, weekday, forcePrevious), nil
	}

	return time.Time{}, ErrInvalidRecordQueryDate
}

func resolveQueryLocation() *time.Location {
	loc, err := time.LoadLocation(DefaultQueryTimezone)
	if err != nil {
		return time.Local
	}
	return loc
}

func parseAbsoluteDate(raw string, loc *time.Location) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
		"02.01.2006",
		"2.1.2006",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, raw, loc); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

func normalizeDateExpression(raw string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "â", "a", "ã", "a",
		"é", "e", "ê", "e",
		"í", "i",
		"ó", "o", "ô", "o", "õ", "o",
		"ú", "u",
		"ç", "c",
		"-", " ", "_", " ", ",", " ", ";", " ", ".", " ", ":", " ",
	)
	norm := strings.ToLower(strings.TrimSpace(raw))
	norm = replacer.Replace(norm)
	norm = strings.Join(strings.Fields(norm), " ")
	return norm
}

func extractWeekday(norm string) (time.Weekday, bool) {
	for weekday, aliases := range weekdayAliases {
		for _, alias := range aliases {
			if containsToken(norm, normalizeDateExpression(alias)) {
				return weekday, true
			}
		}
	}
	return 0, false
}

func containsToken(text, token string) bool {
	if token == "" {
		return false
	}
	padded := " " + text + " "
	return strings.Contains(padded, " "+token+" ")
}

func resolveWeekday(now time.Time, target time.Weekday, forcePrevious bool) time.Time {
	today := startOfDay(now)
	diff := (int(today.Weekday()) - int(target) + 7) % 7
	if diff == 0 && forcePrevious {
		diff = 7
	}
	return today.AddDate(0, 0, -diff)
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
