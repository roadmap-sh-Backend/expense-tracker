package main

import "time"

type Expenses struct {
	Expenses []Expense `json:"expenses"`
}

type Expense struct {
	ID          int       `json:"id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpsertExpense struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}
