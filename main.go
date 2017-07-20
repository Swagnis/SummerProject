package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Computer class for Computer type
type Computer struct {
	num   int
	name  string
	login string
	pass  string
	ip    string
	port  int
}

func main() {
	helpArg := flag.Bool("help", false, "a boolean")              // Отображает доступные команды
	pathArg := flag.String("path", "./", "a string")              // Указывает путь на data.json
	newArg := flag.Bool("new", false, "a boolean")                // Создает новое подключение к компьютеру
	nameArg := flag.String("name", "Unknown", "a string")         // Для задания псеводнима компьютера
	loginArg := flag.String("login", "admin", "a string")         // Для задания логина для подключения к роутеру
	passArg := flag.String("pass", "", "a string")                // Для задания пароля для подключения к роутеру
	ipArg := flag.String("ip", "0.0.0.0", "a string")             // Для задания ip компьютера
	portArg := flag.Int("port", 22, "an int")                     // Для задания порта SSH соединения
	doArg := flag.Bool("do", false, "a boolean")                  // Запись одного или всех конфигов
	allArg := flag.Bool("all", false, "a boolean")                // Писать вместе с do, если необходимо записать конфиги со всех компьютеров
	conectedArg := flag.Bool("conected", false, "a boolean")      // Выводит список подключенных компьютеров
	hashtorageArg := flag.Bool("hashstorage", false, "a boolean") // Выводит hashstorage по имени компьютера
	exportArg := flag.Bool("export", false, "a boolean")          // Экспортирует конфиг в папку на компьютере
	timeArg := flag.String("time", "20:24", "a string")           // Писать вместе с export для задания времени для экспорта
	dateArg := flag.String("date", "18.07.2017", "a str")         // Писать вместе с export для задания даты для экспорта
	flag.Args()
	flag.Parse()

	params := getData(*pathArg)
	names := flag.Args()

	if *helpArg == true {
		fmt.Printf("\n------------------------------------------------------------------------")
		fmt.Printf("\nДля создания подключения к компьютеру ввести данные по примеру: -new -name (Псевдоним  компьютера) -login (Логин пользователя) -pass (Пароль пользователя) -ip (ip компьютера) -port (Порт соединения. По умолчание - 22")
		fmt.Printf("\n")
		fmt.Printf("\nДля получения конфига с одного компьютера ввести данные по примеру: -do (Псевдоним компьютера)")
		fmt.Printf("\n")
		fmt.Printf("\nЧтобы получить конфиги со всех подключенных компьютеров ввести данные по примеру: -do -all ")
		fmt.Printf("\n")
		fmt.Printf("\nДля того, чтобы узнать какие омпьютеры подключены ввести даные по примеру: -conected")
		fmt.Printf("\n")
		fmt.Printf("\nДля того, чтобы вывести хранилище с хэшами ввести данные по примеру: -hashstorage -name (Псевдоним компьютера)")
		fmt.Printf("\n")
		fmt.Printf("\nДля того, чтобы экспортировать файл из таблицы конфигов ввести данные по примеру: -export -name (Псевдоним компьютера) -date (Дата занесения в таблицу) -time (Время занесения в таблицу)")
		fmt.Printf("\n")
		fmt.Printf("\n------------------------------------------------------------------------")
	}

	if *newArg != false {
		var r Computer
		newConnection(r, *nameArg, *ipArg, *loginArg, *portArg, *passArg, params)
		fmt.Println()
	}

	if *doArg != false {

		if *allArg != false {
			doAllConfig(params, *ipArg, *portArg, *loginArg, *passArg, names)
		} else {
			doConfig(params, *ipArg, *portArg, *loginArg, *passArg, names)
		}
	}

	if *conectedArg != false {
		current := time.Now()
		fmt.Println("Date: ", current.Format("02.01.2006"))
		fmt.Println("Time: ", current.Format("15:04"))
		printAllConnected(connectDB(params[0]))
	}

	if *hashtorageArg != false {
		hashStorage(*nameArg, params)
	}

	if *exportArg != false {
		export(*nameArg, *dateArg, *timeArg, params)
	}

}

// newConnection
func newConnection(r Computer, name string, ip string, login string, port int, pass string, params [2]string) {
	r.name = name
	r.ip = ip
	r.login = login
	r.port = port
	r.pass = pass
	fmt.Println("----------------------------------------------------------------------------------------------")
	fmt.Println("\nКомпьютер: \t", r.name, "\t", r.ip, "\t", r.login, "\t", r.pass, "\t", r.port, "\t\tбыл подключен")
	fmt.Println("\n----------------------------------------------------------------------------------------------")

	addNewComputer(r, params[0])
}

// sshComputer
func sshComputer(login string, pass string, ip string, port int) *ssh.Client {
	config := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
	}
	fmt.Println("Successfully connected to ", ip, ":", port)

	return client
}

// sftpComputer
func sftpComputer(client *ssh.Client, path string) {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		fmt.Printf("1: %s", err)
	}
	defer sftp.Close()

	srcPath := "/cf/conf/"
	config := "config.xml"

	srcFile, err := sftp.Open(srcPath + config)
	if err != nil {
		fmt.Printf("2: %s", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(path + config)
	if err != nil {
		fmt.Printf("#3: %s", err)
	}
	defer dstFile.Close()

	srcFile.WriteTo(dstFile)

}

func getData(path string) [2]string {
	data := "data.json"

	var str [2]string
	file, err := os.Open(path + data)
	if err != nil {
		fmt.Println("20")
	}
	defer file.Close()
	f := bufio.NewReader(file)
	for i := 0; i < len(str); i++ {
		str[i], _ = f.ReadString('\n')
	}

	return str
}

// ConnectDB
func connectDB(settings string) *sql.DB {

	// Подключение к БД
	db, err := sql.Open("postgres", settings)
	if err != nil {
		fmt.Println("4")
	}
	return db
}

// addNewComputer
func addNewComputer(r Computer, settings string) {
	db, err := sql.Open("postgres", settings)

	var lastInsertID int
	err = db.QueryRow("INSERT INTO computerstest(ip, login, pass, port, name) VALUES($1,$2,$3,$4,$5) returning id;", r.ip, r.login, r.pass, r.port, r.name).Scan(&lastInsertID)
	if err != nil {
		fmt.Println("5")
	}

	defer db.Close()
}

// addNewHash
func addNewHash(md5 string, sha1 string, name string, params [2]string) {
	current := time.Now()
	time := current.Format("15:04")
	date := current.Format("02.01.2006")

	db, err := sql.Open("postgres", params[0])
	if err != nil {
		fmt.Println("sql.DB.Open()")
	}

	defer db.Close()

	var lastInsertID int
	err = db.QueryRow("INSERT INTO hashstoragetest(date, time, md5config, sha1config, name) VALUES($1,$2,$3,$4,$5) returning id;", date, time, md5, sha1, name).Scan(&lastInsertID)
	if err != nil {
		fmt.Println("5,5")
	}

}

// addNewFile
func addNewFile(params [2]string, r Computer) {
	current := time.Now()
	time := current.Format("15:04")
	date := current.Format("02.01.2006")

	db, err := sql.Open("postgres", params[0])
	if err != nil {
		fmt.Println("6")
	}

	defer db.Close()

	file, err := os.Open(params[1] + "config.xml")
	if err != nil {
		fmt.Println("7")
	}
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	bytes := make([]byte, fileSize)

	var lastInsertID int
	err = db.QueryRow("INSERT INTO configstest(date, time, name, conf) VALUES($1,$2,$3,$4) returning id;", date, time, r.name, bytes).Scan(&lastInsertID)
	if err != nil {
		fmt.Println("8")
	}

}

// printAllConnected
func printAllConnected(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM computerstest")
	if err != nil {
		fmt.Println("9")
	}

	for rows.Next() {
		var id int
		name := ""
		ip := ""
		port := 0
		login := ""
		pass := ""
		err = rows.Scan(&id, &ip, &login, &pass, &port, &name)
		if err != nil {
			fmt.Println("10")
		}
		fmt.Println("id | ip | login | pass | port | name")
		fmt.Printf("%3v | %8v | %8v | %6v | %8v | %3v\n", id, ip, login, pass, port, name)

	}
	defer db.Close()
}

// ComputerData
func ComputerData(db *sql.DB, params [2]string) {
	var r Computer
	rows, err := db.Query("SELECT * FROM computerstest")
	if err != nil {
		fmt.Println("11")
	}

	for rows.Next() {
		var id int
		name := ""
		ip := ""
		port := 0
		login := ""
		pass := ""
		err = rows.Scan(&id, &ip, &login, &pass, &port, &name)
		if err != nil {
			fmt.Println("12")
		}

		fmt.Println("id | ip | login | pass | port | name")
		fmt.Printf("%3v | %8v | %6v | %8v | %6v | %3v\n", id, ip, login, pass, port, name)

		r.num = id
		r.name = name
		r.ip = ip
		r.port = port
		r.login = login
		r.pass = pass

		sftpComputer(sshComputer(r.login, r.pass, r.ip, r.port), params[1])
		sqlComputer(params, r)

	}
	defer db.Close()
}

func sqlComputer(params [2]string, r Computer) {

	addNewHash(hashMD5config(params[1]), hashSHA1config(params[1]), r.name, params)
	// fmt.Println("MD5: ", hashMD5config(params[1]), " SHA-1: ", hashSHA1config(params[1]))
	addNewFile(params, r)

}

func hashMD5config(path string) string {
	config := "config.xml"

	configFile, err := os.Open(path + config)
	if err != nil {
		fmt.Print("18", err)
	}
	defer configFile.Close()

	configHash := md5.New()
	if _, err := io.Copy(configHash, configFile); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", configHash.Sum(nil))
}

func hashSHA1config(path string) string {
	config := "config.xml"

	configFile, err := os.Open(path + config)
	if err != nil {
		fmt.Print("19", err)
	}
	defer configFile.Close()

	configHash := sha1.New()
	if _, err := io.Copy(configHash, configFile); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", configHash.Sum(nil))
}

// ComputerDataNm
func ComputerDataNm(db *sql.DB, params [2]string, names []string) {
	var r Computer
	rows, err := db.Query("SELECT * FROM computerstest")
	if err != nil {
		fmt.Println("13")
	}

	for rows.Next() {
		var id int
		name := ""
		ip := ""
		port := 0
		login := ""
		pass := ""
		err = rows.Scan(&id, &ip, &login, &pass, &port, &name)
		if err != nil {
			fmt.Println("14")
		}
		fmt.Println("id | ip | login | pass | port | name")
		fmt.Printf("%3v | %8v | %6v | %8v | %6v | %3v\n", id, ip, login, pass, port, name)

		r.num = id
		r.name = name
		r.ip = ip
		r.port = port
		r.login = login
		r.pass = pass

		for i := 0; i < len(names); i++ {
			if names[i] == name {
				sftpComputer(sshComputer(r.login, r.pass, r.ip, r.port), params[1])
				sqlComputer(params, r)
			}
		}
	}
	defer db.Close()
}

// hashStorage
func hashStorage(name string, params [2]string) {
	db, err := sql.Open("postgres", params[0])
	if err != nil {
		fmt.Println("15")
	}

	cmd := "SELECT * FROM hashstoragetest WHERE name LIKE '" + name + "'"
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println("16")
	}

	for rows.Next() {
		var id int
		name := ""
		date := ""
		time := ""
		md5config := ""
		sha1config := ""
		err = rows.Scan(&id, &date, &time, &md5config, &sha1config, &name)
		if err != nil {
			fmt.Println("17")
		}

		fmt.Println("id | date | time | md5config | sha1config | name |")
		fmt.Printf("%3v | %8v |  %4v | %8v | %8v | %8v\n", id, date, time, md5config, sha1config, name)
	}
	defer db.Close()
}

// doAllConfig
func doAllConfig(params [2]string, ip string, port int, login string, pass string, names []string) {
	ComputerData(connectDB(params[0]), params)
}

// doConfig
func doConfig(params [2]string, ip string, port int, login string, pass string, names []string) {
	ComputerDataNm(connectDB(params[0]), params, names)
}

// export
func export(_name string, _date string, _time string, params [2]string) {

	cmd := "SELECT * FROM configstest WHERE name LIKE '" + _name + "'"
	db := connectDB(params[0])
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Println("21")
	}

	for rows.Next() {

		var id int
		var time string
		var conf []byte
		var name string
		var date string
		fmt.Println(_name)
		err = rows.Scan(&id, &date, &time, &name, &conf)
		if err != nil {
			fmt.Println("22")
		}
		fmt.Println("id | date | time | name | conf")
		fmt.Printf("%3v | %6v | %3v | %8v | %1v\n", id, date, time, name, conf)
		if _date == date {
			if _time == time {
				convertToFile(conf)

			}
		}

	}
	defer db.Close()

}

// convertToFile
func convertToFile(bytes []byte) {
	permissions := os.FileMode(0644)
	bytes = []byte("to be written to a file\n")
	cfgerr := ioutil.WriteFile("C:/Go/Projects/Newfolder/Pasted/config.xml", bytes, permissions)
	if cfgerr != nil {
		fmt.Println("23")
	}

}
