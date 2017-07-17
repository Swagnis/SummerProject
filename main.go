package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Computer class for Computer type
type Computer struct {
	login string
	pass  string
	ip    string
	port  int
}

func pComputer(p Computer) {
	fmt.Println("login: ", p.login, " pass: ", p.pass, " ip: ", p.ip)
}

func main() {
	helpArg := flag.Bool("help", false, "a boolean")     // Справочник
	newArg := flag.Bool("new", false, "a boolean")       // Подключение к компьютеру
	ipArg := flag.String("ip", "0.0.0.0", "a string")    // Задает Ip
	loginArg := flag.String("login", "root", "a string") // Логин для подключения
	passArg := flag.String("pass", "", "a string")       // Пароль для подключения
	portArg := flag.Int("port", 22, "an int")            // Порт SSH соединения

	flag.Parse()

	if *helpArg == true {
		fmt.Printf("\n----------------------------------------------------------------------------------------------")
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("\nДля подключения к компьютеру ввести данные по примеру: -new -login 'Логин пользователя' -pass 'Пароль пользователя' -ip 'ip компьютера' -port 'Порт соединения. По умолчание - 22' ")
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("\n----------------------------------------------------------------------------------------------")
	}

	if *newArg != false {
		var p Computer
		newConnection(p, *ipArg, *loginArg, *portArg, *passArg)
		copyFile(*loginArg, *passArg, *ipArg, *portArg)
	}

}

// Создание ssh-клинета и sftp соединения
func copyFile(loginArg string, passArg string, ipArg string, portArg int) {
	sftpConnect(sshCli(loginArg, passArg, ipArg, portArg))
}

// sshClient : Функция, создающая ssh клиент и sftp соединение
func sshCli(loginArg string, passArg string, ipArg string, portArg int) *ssh.Client {
	config := &ssh.ClientConfig{
		User: loginArg,
		Auth: []ssh.AuthMethod{
			ssh.Password(passArg),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", ipArg, portArg)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println("Successfully connected to ", ipArg, ":", portArg)

	return client
}

func sftpConnect(client *ssh.Client) {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer sftp.Close()

	srcPath := "/cf/conf/config.xml"
	dstPath := "C:/Go/Projects/New folder/Copied/"
	config := "config.xml"

	srcFile2, err := sftp.Open(srcPath + config)
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Printf("\n----------------------------------------------------------------------------------------------")
	}
	defer srcFile2.Close()

	dstFile2, err := os.Create(dstPath + config)
	if err != nil {
		fmt.Printf("Error: %s", err)
		fmt.Printf("\n----------------------------------------------------------------------------------------------")
	}
	defer dstFile2.Close()

	srcFile2.WriteTo(dstFile2)
}

func newConnection(p Computer, ip string, login string, port int, pass string) {
	p.ip = ip
	p.login = login
	p.port = port
	p.pass = pass
	fmt.Println("\n C компьютера по адресу: \t", "\t", ip, "\t", login, "\t", pass, "\t", port, "\t\t взят файл конфигурации")
	fmt.Println("\n----------------------------------------------------------------------------------------------")
}
