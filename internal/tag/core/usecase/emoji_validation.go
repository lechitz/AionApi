package usecase

import (
	"strings"
	"unicode"
)

func normalizeTagIcon(icon *string) string {
	if icon == nil {
		return DefaultTagIcon
	}

	trimmed := strings.TrimSpace(*icon)
	if trimmed == "" {
		return DefaultTagIcon
	}
	return trimmed
}

func isSingleEmoji(value string) bool {
	if value == "" {
		return false
	}

	runes := []rune(value)
	if isFlagEmoji(runes) || isKeycapEmoji(runes) {
		return true
	}

	i := 0
	if !isEmojiBase(runes[i]) {
		return false
	}
	i++
	i = consumeEmojiTail(runes, i)

	for i < len(runes) {
		if !isZeroWidthJoiner(runes[i]) {
			return false
		}
		i++
		if i >= len(runes) || !isEmojiBase(runes[i]) {
			return false
		}
		i++
		i = consumeEmojiTail(runes, i)
	}

	return i == len(runes)
}

func consumeEmojiTail(runes []rune, i int) int {
	for i < len(runes) {
		if isVariationSelector(runes[i]) || isEmojiModifier(runes[i]) || unicode.Is(unicode.Mn, runes[i]) {
			i++
			continue
		}
		break
	}
	return i
}

func isEmojiBase(r rune) bool {
	return unicode.Is(unicode.So, r)
}

func isEmojiModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

func isVariationSelector(r rune) bool {
	return r == 0xFE0F || r == 0xFE0E
}

func isZeroWidthJoiner(r rune) bool {
	return r == 0x200D
}

func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

func isFlagEmoji(runes []rune) bool {
	return len(runes) == 2 && isRegionalIndicator(runes[0]) && isRegionalIndicator(runes[1])
}

func isKeycapEmoji(runes []rune) bool {
	if len(runes) == 2 {
		return isKeycapBase(runes[0]) && runes[1] == 0x20E3
	}
	if len(runes) == 3 {
		return isKeycapBase(runes[0]) && isVariationSelector(runes[1]) && runes[2] == 0x20E3
	}
	return false
}

func isKeycapBase(r rune) bool {
	return (r >= '0' && r <= '9') || r == '#' || r == '*'
}
