package manager

import "strings"

func escapeFsRestrictedChars(name string) string {
	return strings.NewReplacer(
		"\\", "╲",
		"/", "╱",
		":", "꞉",
		"*", "＊",
		"?", "？",
		"\"", "＂",
		"<", "˂",
		">", "˃",
		"|", "∣",
	).Replace(name)
}
