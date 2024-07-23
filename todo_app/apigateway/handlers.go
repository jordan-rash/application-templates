package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	services "github.com/nats-io/nats.go/micro"
)

type todoStatus string

var (
	todoNotStarted todoStatus = "Not Started"
	todoInProgress todoStatus = "In Progress"
	todoCompleted  todoStatus = "Completed"
)

type todo struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      todoStatus `json:"status,omitempty"`
}

func createHandler(nc *nats.Conn) func(req services.Request) {
	return func(req services.Request) {
		td := new(todo)
		err := json.Unmarshal(req.Data(), td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		if td.Name == "" {
			err = req.Respond([]byte("Name is required"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			}
			return
		}

		td.ID = uuid.NewString()
		td.Status = todoNotStarted

		td_raw, err := json.Marshal(td)
		if err != nil {
			err = req.Respond([]byte("Failed to marshal todo"))
			fmt.Fprintf(os.Stderr, "Failed to marshal todo: %v\n", err)
			return
		}

		// Create todo
		todoRaw, err := nc.Request(fmt.Sprintf("todo.%s.create", strings.ReplaceAll(VERSION, ".", "_")), td_raw, time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error requesting read: %v\n", err)
			_ = req.Respond([]byte("Error creating todo"))
			return
		}

		err = json.Unmarshal(todoRaw.Data, td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		resp, err := json.Marshal(td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		fmt.Println(string(todoRaw.Data))
		_ = req.Respond(resp)
	}
}

func readHandler(nc *nats.Conn) func(req services.Request) {
	return func(req services.Request) {
		td := new(todo)
		err := json.Unmarshal(req.Data(), td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		// Ask for todo by ID
		todoResp, err := nc.Request(fmt.Sprintf("todo.%s.read", strings.ReplaceAll(VERSION, ".", "_")), req.Data(), time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error requesting read: %v\n", err)
			_ = req.Respond([]byte("Error reading todo"))
			return
		}

		type readResp struct {
			Todo   json.RawMessage `json:"todo,omitempty"`
			Error  string          `json:"error,omitempty"`
			Status string          `json:"status"`
			Keys   []string        `json:"keys,omitempty"`
		}
		resp := new(readResp)
		err = json.Unmarshal(todoResp.Data, &resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		r, _ := json.Marshal(resp)
		_ = req.Respond(r)
	}
}

func updateHandler(nc *nats.Conn) func(req services.Request) {
	return func(req services.Request) {
		td := new(todo)
		err := json.Unmarshal(req.Data(), td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		if td.ID == "" {
			err = req.Respond([]byte("ID is required"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			}
			return
		}

		var newStatus todoStatus
		if td.Status == "" {
			err = req.Respond([]byte("Status is required"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			}
			return
		} else {
			newStatus = td.Status
		}

		// Get todo by ID
		todoRaw, err := nc.Request(fmt.Sprintf("todo.%s.read", strings.ReplaceAll(VERSION, ".", "_")), req.Data(), time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error requesting reat: %v\n", err)
			_ = req.Respond([]byte("Error reading todo"))
			return
		}
		err = json.Unmarshal(todoRaw.Data, td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		switch newStatus {
		case "Not Started":
			td.Status = todoNotStarted
			fmt.Fprintf(os.Stdout, "TODO [%s] status changed to %s\n", td.ID, todoNotStarted)
		case "In Progress":
			td.Status = todoInProgress
			fmt.Fprintf(os.Stdout, "TODO [%s] status changed to %s\n", td.ID, todoInProgress)
		case "Completed":
			td.Status = todoCompleted
			fmt.Fprintf(os.Stdout, "TODO [%s] status changed to %s\n", td.ID, todoCompleted)
		default:
			err = req.Respond([]byte("Invalid new status"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			}
			return
		}

		newTodo, err := json.Marshal(td)
		if err != nil {
			err = req.Respond([]byte("Failed to marshal todo"))
			fmt.Fprintf(os.Stderr, "Failed to marshal todo: %v\n", err)
			return
		}

		// Update todo by ID
		_, err = nc.Request(fmt.Sprintf("todo.%s.update", strings.ReplaceAll(VERSION, ".", "_")), newTodo, time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error requesting reat: %v\n", err)
			_ = req.Respond([]byte("Error reading todo"))
			return
		}

		req.Respond([]byte("TODO Updated"))
	}
}

func deleteHandler(nc *nats.Conn) func(req services.Request) {
	return func(req services.Request) {
		td := new(todo)
		err := json.Unmarshal(req.Data(), td)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		if td.ID == "" {
			err = req.Respond([]byte("ID is required"))
			fmt.Fprintf(os.Stderr, "Error responding to request: %v\n", err)
			return
		}

		// Delete todo by ID
		_, err = nc.Request(fmt.Sprintf("todo.%s.delete", strings.ReplaceAll(VERSION, ".", "_")), req.Data(), time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error requesting delete: %v\n", err)
			_ = req.Respond([]byte("Error deleting todo"))
			return
		}

		_ = req.Respond([]byte("Successfully deleted todo"))
	}
}
