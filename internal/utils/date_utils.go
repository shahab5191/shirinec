package utils

import (
	"fmt"
	"time"
)

func DurationToPostgresqlInterval(d time.Duration) string {
    minutes := int64(d.Minutes())
    return fmt.Sprintf("%02d Minutes", minutes)
}
