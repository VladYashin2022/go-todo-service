package cli

import (
	"bufio"
	"cli_todo/service"
	"errors"
	"fmt"
	"os"
)

func Run() error {

	//UI
	for {
		fmt.Println("Выберите действие:\n" +
			"1. Создать новую задачу\n" +
			"2. Прочитать все задачи.\n" +
			"3. Прочитать одну задачу.\n" +
			"4. Изменить задачу.\n" +
			"5. Удалить задачу из списка.\n" +
			"6. Выход из программы.")

		var choose int
		_, err := fmt.Scan(&choose)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch choose {
		case 1:
			//Создание новой задачи
			name, err := writeNameAsk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			date, err := writeDateAsk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			_, err = service.CreateTask(name, date)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case 2:
			//Вывод всех задач
			fmt.Println(service.ReadAllTasks(service.AllTasks))

		case 3:
			//Вывод одной задачи по id
			id, err := chooseIDAsk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			data, errRead := service.ReadTask(id, service.AllTasks)
			if errRead != nil {
				fmt.Println(errRead)
				continue
			}
			fmt.Println(data)

		case 4:
			//Изменение параметров задачи
			id, err := chooseIDAsk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			chooseUpdate, newName, newDate, err := chooseAsk()
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = service.UpdateTask(chooseUpdate, id, newName, newDate, &service.AllTasks)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(service.ReadAllTasks(service.AllTasks))
		case 5:
			//Удаление задачи из списка
			for {
				id, err := chooseIDAsk()
				if err != nil {
					fmt.Println(err)
					break
				}
				errDelete := service.DeleteTask(id, &service.AllTasks)
				if errDelete != nil {
					fmt.Println(errDelete)
					break
				}
				break
			}
		case 6:
			os.Exit(0)
		default:
			fmt.Println("Не существует такого варианта ответа.")
			continue
		}

	}
}

// функция для вопроса пользователю какое имя для задчи выбрать
func writeNameAsk() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите название задачи: ")

	taskName, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read name task: %w", err)
	}
	return taskName, nil
}

// функция для вопроса пользователю какую дату для задчи выбрать
func writeDateAsk() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите дату задачи в формате (23-12-2025 15:30): ")
	readData, errRead := reader.ReadString('\n')
	if errRead != nil {
		return "", fmt.Errorf("read date task string: %w", errRead)
	}
	return readData, nil
}

func chooseIDAsk() (int, error) {
	fmt.Print("Выберите ID задачи: ")
	var choose int
	_, err := fmt.Scan(&choose)
	if err != nil {
		fmt.Println(err, "Ошибка в chooseIDAsk()")
		return 0, err
	}
	return choose, nil
}

func chooseUpdateAsk() (int, error) {
	fmt.Println("Что изменить?\n1. Название задачи.\n2. Дату задачи.\n3. Оба параметра")
	var choose int
	_, err := fmt.Scan(&choose)
	if err != nil {
		fmt.Println(err, "Ошибка в chooseUpdateAsk()")
		return 0, err
	}
	if choose < 1 || choose > 3 {
		return 0, errors.New("Такого варианта ответа не существует.")
	}
	return choose, nil
}

// функция для ввода в консоль решения и в зависимости от него спрагиваются новое имя и(или) дата для задачи
func chooseAsk() (int, string, string, error) {
	choose, err := chooseUpdateAsk()
	if err != nil {
		fmt.Println(err)
		return 0, "", "", err
	}

	var newName, newDate string

	switch choose {
	case 1:
		newName, err = writeNameAsk()

		if err != nil {
			return 0, "", "", fmt.Errorf("choose ask: %w", err)
		}

		return choose, newName, "", nil

	case 2:
		newDate, err = writeDateAsk()

		if err != nil {
			return 0, "", "", fmt.Errorf("choose ask: %w", err)
		}

		return choose, "", newDate, nil

	case 3:
		newName, err = writeNameAsk()

		if err != nil {
			return 0, "", "", fmt.Errorf("choose ask: %w", err)
		}

		newDate, err = writeDateAsk()

		if err != nil {
			return 0, "", "", fmt.Errorf("choose ask: %w", err)
		}

		return choose, newName, newDate, nil
	default:
		return 0, "", "", fmt.Errorf("choose ask: %w", err)
	}
}
