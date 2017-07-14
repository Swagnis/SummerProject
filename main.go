package main

import (
	"flag"
	"fmt"
)

// Computer class for Computer type
type Computer struct {
	login string
	pass  string
	ip    string
	port  int
}

func main() {
	startArg := flag.Bool("start", false, "a boolean")    // Запуск
	stopArg := flag.Bool("stop", false, "a boolean")      // Остановка
	ipArg := flag.String("ip", "0.0.0.0", "a string")     // Задает Ip
	conArg := flag.Bool("new", false, "a boolean")        // Создание нового подключения
	loginArg := flag.String("login", "admin", "a string") // Логин для подключения
	passArg := flag.String("pass", "", "a string")        // Пароль для подключения
	hostArg := flag.String("host", "Default", "a string") // Хостнейм
	portArg := flag.Int("port", 22, "an int")             // Порт SSH соединения

	flag.Parse()

	if *startArg == true {
		startProgram()
	}

	if *stopArg == true {
		stopProgram()
	}

}

// старт программы
func startProgram() {

	fmt.Println("Программа запущена")
}

// остановка программы.
func stopProgram() {
	fmt.Println("Программа остановлена")
}
