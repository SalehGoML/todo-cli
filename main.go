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

	"github.com/SalehGoML/constant"
	"github.com/SalehGoML/contract"
	"github.com/SalehGoML/entity"
	"github.com/SalehGoML/filestore"
)



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
	userStorage     []entity.User
	categoryStorage []Category
	taskStorage     []Task

	authenticatedUser *entity.User
	serializationMode  string
)

const (
	useerStoragePath              = "user.txt"

)

//var userFileStore = filestore.New (useerStoragePath, serializationMode)

func main() {
	serializeMode := flag.String("m", constant.ManDaravardiSerializationMode, "Serialization mode")
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

	//loadUserFromStorage(userFileStore, *serializeMode)

	fmt.Println("Hello to TODO app")

	switch *serializeMode {
	case constant.ManDaravardiSerializationMode:

		serializationMode = constant.ManDaravardiSerializationMode
	default:
		serializationMode = constant.JsonSerializationMode
	}

	var userFileStore = filestore.New(useerStoragePath, serializationMode)

	users := userFileStore.Load()
	userStorage = append(userStorage, users...)
	// if there is a user record with corresponding data allow the user to continue

	for {
		runCommand(userFileStore, *command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(store contract.UserReadStore, command string) {
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
			registerUser(store)
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



func registerUser(store contract.UserReadStore) {
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

	user := entity.User{
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

//func loadUserFromStorage(store userReadStore) {
//	users := store.Load(serializationMode)
//
//	userStorage = append(userStorage, users...)
//
//}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
