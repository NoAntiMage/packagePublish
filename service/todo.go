package service

type Todo interface{}

func NewTodoService() Todo {
	return &todo{}
}

type todo struct{}
