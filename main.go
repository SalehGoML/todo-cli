package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

var (
	userStorage     []User
	categoryStorage []Category
	taskStorage     []Task

	authenticatedUser *User
	serializationMode  string
)

const (
	useerStoragePath              = "user.txt"
	ManDaravardiSerializationMode = "mandaravardi"
	JsonSerializationMode         = "json"
)
var userFileToStore = fileStore {
	filePath: useerStoragePath,
}

func main() {
	serializationMode := flag.String("m", ManDaravardiSerializationMode, "Serialization mode")
	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	// log user storage from file
	//loadUserStorageFromFile(*serializationMode)


	//var userReadFileStore userReadStore
	//
	//var userReadStore = fileStore{
	//	filePath:  "./store/data.txt",
	//}
	//
	//userReadFileStore = userReadStore()

	loadUserFromStorage(userFileStore, *serializeMode)

	fmt.Println("Hello to TODO app")

	switch *serializeMode {
	case ManDaravardiSerializationMode:
		serializationMode = ManDaravardiSerializationMode
	default:
		serializationMode = JsonSerializationMode
	}
	// if there is a user record with corresponding data allow the user to continue

	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil
		}


		switch command {
		case "create-task":
			createTask()
		case "create-category":
			createCategory()
		case "register-user":
			registerUser(userFileStore)
		case "list-task":
			listTask()
		case "login":
			login()
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("command is not valid", command)
	}
}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, dueDate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category")
	scanner.Scan()
	category = scanner.Text()

	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("categor-id is not valid integer, #{err}\n")

		return
	}

	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenticatedUser.ID {
			isFound = true

			break
		}
	}

	if !isFound {
		fmt.Printf("category-id is not found\n")

		return
	}
	fmt.Println("please enter the task due date")
	scanner.Scan()
	dueDate = scanner.Text()

	// validation
	// category validate

	task := Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    dueDate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}

	taskStorage = append(taskStorage, task)
}

func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()
	fmt.Println("category", title, color)

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title: title,
		Color: color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, c)
}

type userWriteStore interface {
	Save(u User)
}

type userReadStore interface {
	Load(serializationMode string) []User
}

func registerUser(store userStore) {
	scanner := bufio.NewScanner(os.Stdin)
	var id, name, email, password string

	fmt.Println("please enter the name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the email")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = email

	fmt.Println("user", id, email, password)

	user := User{
		ID: len(userStorage) + 1,
		Name: name,
		Email: email,
		Password: password,
	}

	userStorage = append(userStorage, user)

	// writeUserToFile(user)
	store.Save(user)
}

func login() {
	fmt.Println("login process")
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("please enter email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if user.Email == email && user.Password == hashThePassword(password) {
			authenticatedUser = &user

			break
		}
	}

	if authenticatedUseer == nil {
		fmt.Println("the email or password is not correct")
	}


}

func listTask() {
	for _, user := range  taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}

func loadUserFromStorage(store userReadStore) {
	users := store.Load(serializationMode)

	userStorage = append(userStorage, users...)

}

func (f fileStore) writeUserToFile(user User) {
	var file *os.File

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cant create or open file", err)

		return
	}
	defer file.Close()

	var data []byte
	// serialize the user struct/object
	if serializationMode == ManDaravardiSerializationMode {
		data = []byte(fmt.Sprintf("id: #{user.ID}, name: #{user.Name}, email: #{user.Email}, password: #{user.Password})\n"),)
	} else if serializationMode == JsonSerializationMode {
		//json

		var jErr error
		data, jErr = json.Marshal(user)
		if jErr != nil {
			fmt.Println("can't marshal user struct to json", jErr)

			return
		}
		data = append(data, []byte("\n")...)
	} else {
		fmt.Println("invaild serialization mode")

		return
	}

	numberOfWrittenBytes, wErr := fil.Write(data)
	if wErr != nil {
		fmt.Printf("can't write to the file #{wErr}\n")

		return
	}

	fmt.Println("numberOfWrittenBytes", numberOfWrittenBytes)
}

func deserilizeFromManDaravardi(userStr string) (User, error) {
	if userStr == "" : User{}, errors.New("user string is empty")

	var user = User{}

	userFields := strings.Split(userStr, ",")
	for _, userField := range userFields {
		values := strings.Split(field, ":")
		if len(values) != 2 {
			fmt.Println("field is not valid, skipping...", len(values))

			continue
		}
	}
	return user,nil

}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
type fileStore struct {
	filePath string
}

func (f fileStore) Save(u User) {
	f.writeUserToFile(u)
}

func (f fileStore) Load(serializationMode string) []User {
	var uStore []User

	file, err := os.Open(f.filePath)
	if err != nil {
		fmt.Println("can't open the file", err)
	}
	var data = make([]byte, 1024)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from the file", oErr)

		return nil
	}

	var dataStr = string(data)

	userSlice := strings.Split(dataStr, "\n")
	fmt.Println("len userSlice", len(userSlice), serializationMode)
	for _, u := range userSlice {
		var userStruct = User{}

		switch serializationMode {
		case ManDaravardiSerializationMode:
			var dErr error
			userStruct, dErr = deserilizeFromManDaravardi(u)
			if dErr != nil {
				fmt.Println("can't deserialize user record to user struct", dErr)

				return nil
			}
		case JsonSerializationMode:
			if u[0] != '{' && u[len(u)-1] != '}' {
				continue
			}

			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("can't deserialize user record to user struct with json model", uErr)

				return nil
			}
		default:
			fmt.Println("invalid serialization mode")

			return nil
		}

		uStore = append(uStore, userStruct)
	}

	return uStore
}