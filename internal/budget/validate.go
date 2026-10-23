package budget

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

func AddValidation(req *Transaction, today time.Time) error {

	req.Category = Category(Normalise(string(req.Category)))

	if req.Category == "" {
		return errors.New("Category is required.")
	}

	if _, ok := req.Category.Kind(); !ok {
		return errors.New("Uknown category.")
	}

	if req.Amount <= 0 {
		return errors.New("Amount must be positive.")
	}

	parsedDate, err := ParseDate(req.Date)
	if err != nil {
		return err
	}

	if parsedDate.After(today) {
		return errors.New("Date cannot be in the future.")
	}

	req.Description = strings.TrimSpace(req.Description)
	if utf8.RuneCountInString(req.Description) > 100 {
		return errors.New("Must be equal or less than 100 characters.")
	}

	return nil
}
