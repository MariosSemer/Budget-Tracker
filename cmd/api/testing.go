package main

import (
	"budgettracker/internal/budget"
	"fmt"
)

func run(s budget.Store) {
	s.Add(budget.Transaction{Description: "Decathlon", Category: "salary", Amount: 100})
	s.Add(budget.Transaction{Description: "Coffee", Category: "entertainment", Amount: 2})
	s.Add(budget.Transaction{Description: "Sklavenitis", Category: "groceries", Amount: 25})
	s.Add(budget.Transaction{Description: "Reuma", Category: "bills", Amount: 25})
	fmt.Println(s.List())
	fmt.Println("Balance:\t\t", budget.Balance(s))
	s.Delete(2)
	fmt.Println(s.List())
	fmt.Println("Balance:\t\t", budget.Balance(s))
}

func main() {

	run(budget.NewMemoryStore())
}
