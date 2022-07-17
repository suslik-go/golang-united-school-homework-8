package main

import (
	"encoding/json"
	"errors"
	"flag" //command line parsing
	"fmt"
	"io"
	"log"
	"os"
)

type Arguments map[string]string

type User struct {
	id    int
	email string
	age   int
}

//usage: `./main.go -operation «add» -item ‘{«id»: "1", «email»: «email@test.com», «age»: 23}’ -fileName «users.json»`
func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

// Arguments fields: id, item, operation and fileName
func Perform(args Arguments, writer io.Writer) error {

	var user User
	json.Unmarshal([]byte(args["item"]), &user)
	//fmt.Printf("Id: %d, Email: %s, Age: %d", user.Id, user.Email, user.Age)

	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}

	switch {
	case args["operation"] == "add":
		add(user, args["fileName"])
	case args["operation"] == "list":
		writer.Write(list(args["fileName"]))
	case args["operation"] == "findById":
		findById(user, args["fileName"])
	case args["operation"] == "remove":
		remove(user, args["fileName"])
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	return nil
}

func parseArgs() Arguments {
	var id = flag.String("id", "", "command id")
	var userInfo = flag.String("item", "", "user information with this form '{«id»: '1', «email»: «email@test.com», «age»: 23}'")
	var operation = flag.String("operation", "", "operation applied to provided user info")
	var fileName = flag.String("fileName", "", "name of the file where info is being stored")

	flag.Parse()

	return Arguments{
		"id":        *id,
		"item":      *userInfo,
		"operation": *operation,
		"fileName":  *fileName,
	}

}

func add(user User, fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// decode
	var users []User
	jsonDecodeErr := json.Unmarshal(data, &users)
	if jsonDecodeErr != nil {
		fmt.Println("error:", jsonDecodeErr)
	}

	users = append(users, user)

	//encode
	data, jsonEncodeErr := json.Marshal(users)
	if jsonEncodeErr != nil {
		fmt.Println("error:", jsonEncodeErr)
	}

	if _, err := file.Write(data); err != nil {
		file.Close()
		log.Fatal(err)
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

}

func list(fileName string) []byte {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func findById(user User, fileName string) User {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// decode
	var users []User
	jsonDecodeErr := json.Unmarshal(data, &users)
	if jsonDecodeErr != nil {
		fmt.Println("error:", jsonDecodeErr)
	}

	for _, val := range users {
		if val.id == user.id {
			return user
		}
	}
	return user
}

func remove(user User, fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	// decode
	var users []User
	jsonDecodeErr := json.Unmarshal(data, &users)
	if jsonDecodeErr != nil {
		fmt.Println("error:", jsonDecodeErr)
	}
	var Id int = int(user.id)
	users = append(users[:Id], users[Id+1:]...)
}
