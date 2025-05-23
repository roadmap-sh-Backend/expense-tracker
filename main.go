package main

import (
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"os"
	"time"
)

const FILE_NAME = "expenses.json"

func main() {
	err := initializeStorage(FILE_NAME)
	if err != nil {
		log.Fatalf("Failed initialize storage: %v", err)
	}

	args := os.Args
	if len(args) < 3 {
		log.Fatalf("LACKS COMMAND")
	}

	cmd := args[1]
	if cmd != "expense-tracker" {
		log.Fatalf("INVALID COMMAND")
	}

	feat := args[2]

	switch feat {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		category := addCmd.String("category", "", "A category of expenses")
		description := addCmd.String("description", "", "A description of expenses")
		amount := addCmd.Int64("amount", 0, "The amount of expenses")

		err = addCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse add command flags: %v", err)
		}

		expenses, err := CreateExpense(FILE_NAME, UpsertExpense{
			Category:    *category,
			Description: *description,
			Amount:      *amount,
		})
		if err != nil {
			log.Fatalf("Failed adding expense: %v", err)
		}

		slog.Info("Result",
			"status", "Successfully added expense data",
			"expenses", expenses)
	case "delete":
		delCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		id := delCmd.Int("id", 0, "An id of expense")

		err = delCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse delete command flags: %v", err)
		}

		expenses, err := DeleteExpense(FILE_NAME, id)
		if err != nil {
			log.Fatalf("Failed to delete a expense: %v", err)
		}

		slog.Info("Result",
			"status", "Successfully deleted expense data",
			"expenses", expenses)
	case "summary":
		summCmd := flag.NewFlagSet("summary", flag.ExitOnError)
		month := summCmd.Int("month", 0, "A month of expense")

		err = summCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse summary command flags: %v", err)
		}

		expenses, err := GetExpenses(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed to retrieve all expenses: %v", err)
		}

		total := int64(0)
		if *month == 0 {
			for _, expense := range expenses.Expenses {
				if *month == 0 {
					total += expense.Amount
				}
			}

			slog.Info("Expense summary",
				"Total", total,
			)
		} else {
			currYear := time.Now().Year()
			for _, expense := range expenses.Expenses {
				if *month == int(expense.CreatedAt.Month()) && currYear == expense.CreatedAt.Year() {
					total += expense.Amount
				}
			}

			slog.Info("Expense summary",
				"Month", time.Month(*month),
				"Total", total,
			)
		}
	case "list":
		lsCmd := flag.NewFlagSet("list", flag.ExitOnError)
		category := lsCmd.String("category", "", "A category of expense")

		err = lsCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse list command flags: %v", err)
		}

		expenses, err := GetExpenses(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed to retrieve all expenses: %v", err)
		}

		if category != nil {
			slog.Info("Expense list",
				"Category", *category)
			for _, expense := range expenses.Expenses {
				if expense.Category == *category {
					slog.Info("Expense",
						"ID", expense.ID,
						"Description", expense.Description,
						"Amount", expense.Amount,
						"CreatedAt", expense.CreatedAt,
						"UpdatedAt", expense.UpdatedAt,
					)
				}
			}
		} else {
			slog.Info("Expense list")
			for _, expense := range expenses.Expenses {
				slog.Info("Expense",
					"ID", expense.ID,
					"Description", expense.Description,
					"Amount", expense.Amount,
					"CreatedAt", expense.CreatedAt,
					"UpdatedAt", expense.UpdatedAt,
				)
			}
		}
	default:
		log.Fatalf("UNKNOWN COMMAND")
	}
}

func initializeStorage(name string) error {
	_, err := os.Stat(name)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		return os.WriteFile(name, []byte(`{"expenses": []}`), 0644)
	}

	return err
}

func WriteExpense(fileName string, expenses *Expenses) error {
	jsonData, err := json.MarshalIndent(expenses, "", "\t")
	if err != nil {
		log.Println("Error marshalling", err)
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Println("Error writing into a file", err)
		return err
	}

	return nil
}

func GetExpenses(fileName string) (*Expenses, error) {
	byteData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var expenses Expenses
	err = json.Unmarshal(byteData, &expenses)
	if err != nil {
		return nil, err
	}

	return &expenses, nil
}

func CreateExpense(fileName string, payload UpsertExpense) (*Expenses, error) {
	expenses, err := GetExpenses(fileName)
	if err != nil {
		return nil, err
	}

	currID := 0
	if len(expenses.Expenses) > 0 {
		currID = expenses.Expenses[len(expenses.Expenses)-1].ID
	}

	expense := Expense{
		ID:          currID + 1,
		Category:    payload.Category,
		Description: payload.Description,
		Amount:      payload.Amount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expenses.Expenses = append(expenses.Expenses, expense)
	err = WriteExpense(fileName, expenses)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func DeleteExpense(fileName string, id *int) (*Expenses, error) {
	expenses, err := GetExpenses(fileName)
	if err != nil {
		return nil, err
	}

	newExpenses := Expenses{
		Expenses: []Expense{},
	}
	for _, expense := range expenses.Expenses {
		if expense.ID != *id {
			newExpenses.Expenses = append(newExpenses.Expenses, expense)
		}
	}

	err = WriteExpense(fileName, &newExpenses)
	if err != nil {
		return nil, err
	}

	return &newExpenses, nil
}
