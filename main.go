package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"

	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/SalehGoML/constant"
	"github.com/SalehGoML/contract"
	"github.com/SalehGoML/entity"
	"github.com/SalehGoML/repository/filestore"
	"github.com/SalehGoML/repository/memorystore"
	"github.com/SalehGoML/service/Task"
)

var (
	userStorage       []entity.User
	categoryStorage   []entity.Category
	authenticatedUser *entity.User
	serializationMode string
)

const (
	useerStoragePath = "user.txt"
)

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()

	taskService := task.NewService(taskMemoryRepo)

	serializeMode := flag.String("m", constant.ManDaravardiSerializationMode, "Serialization mode")
	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

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
		runCommand(&userFileStore, *command, &taskService)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}
}

func runCommand(store *filestore.FileStore, command string, taskService *task.Service) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil {
			return
		}
	}

	switch command {
	case "create-task":

		createTask(taskService)
	case "create-category":

		createCategory()
	case "register-user":

		registerUser(&filestore)
	case "list-task":
		listTask(taskService)
	case "login":

		login()
	case "exit":

		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)
	}
}

func createTask(taskService *task.Service) {
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

	fmt.Println("please enter the task due date")
	scanner.Scan()
	dueDate = scanner.Text()

	response, err := taskService.Create(task.CreateRequest{
		Title:               title,
		DueDate:             dueDate,
		CategoryID:          categoryID,
		AuthenticatedUserID: authenticatedUser.ID,
	})

	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("create task:", response.Task)

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

	c := entity.Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
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
	email = scanner.Text()

	fmt.Println("please enter the email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the password")
	scanner.Scan()
	password = scanner.Text()

	id = email

	fmt.Println("user", id, email, password)

	user := entity.User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)

	// writeUserToFile(user)
	//store.Save(user)
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

	if authenticatedUser == nil {
		fmt.Println("the email or password is not correct")
	}

}

func listTask(taskService *task.Service) {
	userTasks, err := taskService.List(task.ListRequest{UserID: authenticatedUser.ID})

	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("user tasks", userTasks.Tasks)
}

func hashThePassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])
}
