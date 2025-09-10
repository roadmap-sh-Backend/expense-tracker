
This is my work from [Roadmap.sh](https://roadmap.sh) on [third task](https://roadmap.sh/projects/expense-tracker)  

# Technology
I'm just using Go as programming language.  

# How to run
There are two ways for running this CLI app. First using `go run . expense-tracker <command> <args>` or using my Makefile script `make cmd-<command>`, `ex: make cmd-add category=<string> description=<string> amount=<number>`, `make cmd-update id=<number> category=<string> description=<string> amount=<number>`, etc

# Architecture
I just separate the codes by it's meaning like `types.go` for the structs declaration,  and `main.go` for logic and main executioner of this programs

# Features
1. Users can add an expense with a description and amount. ✅
2. Users can update an expense. ✅
3. Users can delete an expense. ✅
4. Users can view all expenses. ✅
5. Users can view a summary of all expenses. ✅
6. Users can view a summary of expenses for a specific month (of current year). ✅
7. Add expense categories and allow users to filter expenses by category. ✅
8. Allow users to set a budget for each month and show a warning when the user exceeds the budget. ✅
9. Allow users to export expenses to a CSV file. ✅

