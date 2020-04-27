package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MrBTTF/gophercises/task/db"
	"github.com/spf13/cobra"
)

type key int

const CtxDB = key(0)

var (
	Root = &cobra.Command{
		Use:   "task",
		Short: "task is a CLI for managing your TODOs.",
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task to your TODO list",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskName := strings.Join(args, " ")
			db := cmd.Context().Value(CtxDB).(*db.DB)
			err := db.AddTask(taskName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Added \"%s\" to your task list.\n", taskName)
		},
	}

	doCmd = &cobra.Command{
		Use:   "do",
		Short: "Mark a task on your TODO list as complete",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db := cmd.Context().Value(CtxDB).(*db.DB)
			for _, idStr := range args {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Fatal(err)
				}
				err = db.DoTask(id)
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	removeCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove task from list",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			db := cmd.Context().Value(CtxDB).(*db.DB)
			for _, idStr := range args {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Fatal(err)
				}
				err = db.RemoveTask(id)
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all of your incomplete tasks",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := cmd.Context().Value(CtxDB).(*db.DB)
			fmt.Println("You have the following tasks:")
			tasks, err := db.GetTasksNotCompleted()
			for i, task := range tasks {
				if task.Completed == 0 {
					fmt.Printf("%d. %s\n", i+1, task.Name)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	completedCmd = &cobra.Command{
		Use:   "completed",
		Short: "List completed tasks today",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			db := cmd.Context().Value(CtxDB).(*db.DB)
			fmt.Println("You have finished the following tasks today:")
			tasks, err := db.GetTasksCompleted()
			for i, task := range tasks {
				fmt.Printf("%d. %s\n", i+1, task.Name)
			}
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	Root.AddCommand(addCmd)
	Root.AddCommand(doCmd)
	Root.AddCommand(removeCmd)
	Root.AddCommand(listCmd)
	Root.AddCommand(completedCmd)
}
