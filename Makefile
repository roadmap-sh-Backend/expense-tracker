cmd-add:
	@if [ -z "$(description)" ] || [ -z "$(amount)" ] || [ -z "$(category)" ]; then \
		echo "Error: description, amount, and category are required. Usage: make cmd-add category=<string> description=<string> amount=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker add --category "$(category)" --description "$(description)" --amount $(amount)

cmd-update:
	@if [ -z "$(id)" ] || [ -z "$(description)" ] || [ -z "$(amount)" ] || [ -z "$(category)" ]; then \
		echo "Error: id, description, amount, and category are required. Usage: make cmd-update id=<number> category=<string> description=<string> amount=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker update --id $(id) --category "$(category)" --description "$(description)" --amount $(amount)

cmd-ls-all:
	go run . expense-tracker list

cmd-ls:
	@if [ -z "$(category)" ]; then \
		echo "Error: category is required. Usage: make cmd-ls category=<string>"; \
		exit 1; \
	fi
	go run . expense-tracker list --category "$(category)"

cmd-summ-all:
	go run . expense-tracker summary

cmd-summ:
	@if [ -z "$(month)" ]; then \
		echo "Error: month is required. Usage: make cmd-summ month=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker summary --month $(month)

cmd-del:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make cmd-del id=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker delete --id $(id)

cmd-set-budget:
	@if [ -z "$(month)" ] || [ -z "$(amount)" ]; then \
		echo "Error: month and amount are required. Usage: make cmd-set-budget month=<number> amount=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker set-budget --month $(month) --amount $(amount)

cmd-export:
	go run . expense-tracker export
