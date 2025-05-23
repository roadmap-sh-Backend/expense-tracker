cmd-add:
	@if [ -z "$(description)" ] || [ -z "$(amount)"]; then \
		echo "Error: description and amount is required. Usage: make cmd-add description=<string> amount=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker add --description $(description) --amount $(amount)

cmd-ls:
	go run . expense-tracker list

cmd-summ-all:
	go run . expense-tracker summary

cmd-summ:
	go run . expense-tracker summary --month $(month)

cmd-del:
	@if [ -z "$(id)" ]; then \
		echo "Error: id is required. Usage: make cmd-del id=<number>"; \
		exit 1; \
	fi
	go run . expense-tracker delete --id $(id)
