package goblazer

import (
	"fmt"
	"time"
)

// TimeCostStatistics is
func TimeCostStatistics(start time.Time, name string) {
	t := time.Since(start)
	fmt.Printf("[%-32s] costs %12d nanoseconds.\n", name, t)
}
