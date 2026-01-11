package main

import (
	"fintrack/internal/models"
	"fmt"
)

func AddTrannsaction() error {

	var amount int
	var category, description string

	fmt.Println("Введите сумму: ")
	fmt.Scan(&amount)
	fmt.Println("Введите категорию: ")
	fmt.Scan(&category)
	fmt.Println("Введите описание: ")
	fmt.Scan(&description)

	if amount <= 0 {
		return fmt.Errorf("сумма не может быть меньше нуля\n")
	}

	if category == "" {
		return fmt.Errorf("категория не может быть пустой\n")
	} else if category != "income" && category != "expense" {
		return fmt.Errorf("невернро введена категория\n")
	}

	if description == "" {
		return fmt.Errorf("описание не может быть пустым\n")
	}

	fmt.Printf("сумма равна: %d\nописание: %s\nкатегория: %s\n", amount, description, category)

	return nil

}

func main() {
	var choise int
	// var service *services.TransactionService

	for {

		fmt.Println("----------------------Меню----------------------")
		fmt.Println("1. Добавить транзакцию\n2. Показать транзакции\n0. Выход")

		fmt.Scan(&choise)

		if choise == 1 {

			err := AddTrannsaction()

			if err != nil {
				fmt.Print(err)
			}
			choise = 3
		} else if choise == 2 {
			models.GetDefaultCategories()
			choise = 3
		} else if choise == 0 {
			fmt.Println("До свидания!")
			break
		} else if choise == 3 {
			fmt.Println("--------------------------------------------------------------------")
		} else {
			fmt.Println("Только 3 действия в меню.")
		}
	}
}
