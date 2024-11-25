package utils

import "regexp"

func IsValidHexColor(color string) bool {
    re := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
    return re.MatchString(color)
}
