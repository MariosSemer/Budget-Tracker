package budget

import (
	"fmt"
	"log/slog"
)

func Balance(s Store) int64 {
	var total int64
	for _, t := range s.List() {
		kind, ok := t.Category.Kind()
		fmt.Printf("id=%d category=%q kind=%q ok=%v\n", t.ID, t.Category, kind, ok)
		if !ok {
			slog.Warn("Skipping transaction with iknown category", "id", t.ID, "category", t.Category, "kind", kind, "ok", ok)
			continue
		}
		switch kind {
		case "income":
			total += t.Amount
		case "expense":
			total -= t.Amount
		}
	}
	return total
}
