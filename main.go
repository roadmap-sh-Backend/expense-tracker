package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

		personExpenses, err := InsertExpense(FILE_NAME, UpsertExpense{
			Category:    *category,
			Description: *description,
			Amount:      *amount,
		})
		if err != nil {
			log.Fatalf("Failed adding expense: %v", err)
		}

		msg, err := CheckBudgetUsage(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed while checking current budget usage: %v", err)
		} else if msg != "" {
			log.Println(msg)
		}

		slog.Info("Result",
			"status", "Successfully added expense data",
			"expenses", personExpenses)
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		id := updateCmd.Int("id", 0, "An id of expense")
		category := updateCmd.String("category", "", "A category of expenses")
		description := updateCmd.String("description", "", "A description of expenses")
		amount := updateCmd.Int64("amount", 0, "The amount of expenses")

		err = updateCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse update command flags: %v", err)
		}

		expense, err := GetExpenseByID(FILE_NAME, *id)
		if err != nil {
			log.Fatalf("Failed retrieving expense by id: %v", err)
		}
		expenses, err := UpdateExpense(FILE_NAME, expense, UpsertExpense{
			Category:    *category,
			Description: *description,
			Amount:      *amount,
		})
		if err != nil {
			log.Fatalf("Failed updated expense: %v", err)
		}

		msg, err := CheckBudgetUsage(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed while checking current budget usage: %v", err)
		} else if msg != "" {
			log.Println(msg)
		}

		slog.Info("Result",
			"status", "Successfully updated expense data",
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

		personExpenses, err := GetPersonExpense(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed to retrieve all expenses: %v", err)
		}

		total := int64(0)
		if *month < 0 {
			log.Fatalf("Month number cannot be negative")
		} else if *month == 0 {
			for _, expense := range personExpenses.Expenses.Expenses {
				if *month == 0 {
					total += expense.Amount
				}
			}

			slog.Info("Expense summary",
				"Total", total,
			)
		} else {
			currYear := time.Now().Year()
			for _, expense := range personExpenses.Expenses.Expenses {
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

		personExpenses, err := GetPersonExpense(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed to retrieve all expenses: %v", err)
		}

		if *category != "" {
			slog.Info("Expense list",
				"Category", *category)
			for _, expense := range personExpenses.Expenses.Expenses {
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
			for _, expense := range personExpenses.Expenses.Expenses {
				slog.Info("Expense",
					"ID", expense.ID,
					"Description", expense.Description,
					"Amount", expense.Amount,
					"CreatedAt", expense.CreatedAt,
					"UpdatedAt", expense.UpdatedAt,
				)
			}
		}
	case "set-budget":
		setBudgetCmd := flag.NewFlagSet("set-budget", flag.ExitOnError)
		amount := setBudgetCmd.Float64("amount", 0.0, "A budget limit amount")
		month := setBudgetCmd.Int("month", int(time.Now().Month()), "The month for specifying budget limit")
		year := setBudgetCmd.Int("year", int(time.Now().Year()), "The year for specifying budget limit")

		err = setBudgetCmd.Parse(args[3:])
		if err != nil {
			log.Fatalf("Failed to parse set-budget command flags: %v", err)
		}

		personExpense, err := GetPersonExpense(FILE_NAME)
		if err != nil {
			log.Fatalf("Failed to retrieve all expenses: %v", err)
		}

		personExpense.Budget = *amount
		err = WriteExpense(FILE_NAME, personExpense)
		if err != nil {
			log.Fatalf("Failed updating budget limit on month: [%s]. Error: %v", time.Month(*month), err)
		}

		slog.Info("Successfully set budget limit",
			"Month", time.Month(*month),
			"Year", *year,
			"Amount", *amount,
		)
	default:
		log.Fatalf("UNKNOWN COMMAND")
	}
}

func initializeStorage(name string) error {
	_, err := os.Stat(name)
	if err == nil {
		return nil
	}

	strData := fmt.Sprint(`{"budget": 0, "expenses": []}`)

	if os.IsNotExist(err) {
		return os.WriteFile(name, []byte(strData), 0644)
	}

	return err
}

func WriteExpense(fileName string, data *PersonExpense) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
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

func GetPersonExpense(fileName string) (*PersonExpense, error) {
	byteData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var expenses PersonExpense
	err = json.Unmarshal(byteData, &expenses)
	if err != nil {
		return nil, err
	}

	return &expenses, nil
}

func GetExpenseByID(fileName string, id int) (*Expense, error) {
	byteData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var expenses Expenses
	err = json.Unmarshal(byteData, &expenses)
	if err != nil {
		return nil, err
	}

	for _, expense := range expenses.Expenses {
		if expense.ID == id {
			return &expense, nil
		}
	}

	return nil, fmt.Errorf("Expense with id: %d is not found", id)
}

func InsertExpense(fileName string, payload UpsertExpense) (*PersonExpense, error) {
	if payload.Amount < 0 {
		return nil, fmt.Errorf("Amount cannot be negative")
	}

	personExpense, err := GetPersonExpense(fileName)
	if err != nil {
		return nil, err
	}

	currID := 0
	if len(personExpense.Expenses.Expenses) > 0 {
		currID = personExpense.Expenses.Expenses[len(personExpense.Expenses.Expenses)-1].ID
	}

	expense := Expense{
		ID:          currID + 1,
		Category:    payload.Category,
		Description: payload.Description,
		Amount:      payload.Amount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	personExpense.Expenses.Expenses = append(personExpense.Expenses.Expenses, expense)
	err = WriteExpense(fileName, personExpense)
	if err != nil {
		return nil, err
	}

	return personExpense, nil
}

func UpdateExpense(fileName string, expense *Expense, payload UpsertExpense) (*PersonExpense, error) {
	if payload.Amount < 0 {
		return nil, fmt.Errorf("Amount cannot be negative")
	}

	personExpense, err := GetPersonExpense(fileName)
	if err != nil {
		return nil, err
	}

	expense = &Expense{
		ID:          expense.ID,
		Category:    payload.Category,
		Description: payload.Description,
		Amount:      payload.Amount,
		CreatedAt:   expense.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	personExpense, err = DeleteExpense(fileName, &expense.ID)
	if err != nil {
		return nil, err
	}

	personExpense.Expenses.Expenses = append(personExpense.Expenses.Expenses, *expense)
	err = WriteExpense(fileName, personExpense)
	if err != nil {
		return nil, err
	}

	return personExpense, nil
}

func DeleteExpense(fileName string, id *int) (*PersonExpense, error) {
	personExpense, err := GetPersonExpense(fileName)
	if err != nil {
		return nil, err
	}

	newExpenses := Expenses{
		Expenses: []Expense{},
	}
	found := false
	for _, expense := range personExpense.Expenses.Expenses {
		if expense.ID != *id {
			newExpenses.Expenses = append(newExpenses.Expenses, expense)
		} else {
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("Expense with id: %d is not found", *id)
	}

	personExpense.Expenses = newExpenses
	err = WriteExpense(fileName, personExpense)
	if err != nil {
		return nil, err
	}

	return personExpense, nil
}

func CheckBudgetUsage(fileName string) (string, error) {
	personExpense, err := GetPersonExpense(fileName)
	if err != nil {
		return "", err
	}

	amountSum := 0.0
	for _, expense := range personExpense.Expenses.Expenses {
		amountSum += float64(expense.Amount)
	}

	if amountSum >= personExpense.Budget {
		return fmt.Sprintf(
			"Your current summarize expenses for this month: [%s] is exceeding your budget!!!",
			time.Month(time.Now().Month())), nil
	}

	return "", nil
}
