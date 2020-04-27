package db

import (
	"fmt"
	"time"
	"encoding/binary"

	"github.com/boltdb/bolt"
)


const TasksBucket = "tasks"


type Task struct {
	Name      string
	Completed int64
}

type DB struct{
	bolt *bolt.DB
}

func New(dbPath string) (*DB, error){
	db := new(DB)
	var err error
	db.bolt, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TasksBucket))
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}


func (db *DB) Close() error {
	return db.bolt.Close()
}

func (db *DB) AddTask(taskName string) error {
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TasksBucket))
		err := b.Put([]byte(taskName), itob(0))
		return err
	})
	return err
}

func (db *DB) DoTask(id int) error {
	tasks, err := db.GetTasksNotCompleted()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Printf("No task for number %d.\n", id)
		return nil
	}
	completedTime := int(time.Now().Unix())
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TasksBucket))
		err := b.Put([]byte(tasks[id-1].Name), itob(completedTime))
		return err
	})
	fmt.Printf("You have completed the \"%s\" task.\n", tasks[id-1].Name)
	return err
}

func (db *DB) RemoveTask(id int) error {
	tasks, err := db.GetTasksNotCompleted()
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Printf("No task for number %d.\n", id)
		return nil
	}
	err = db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TasksBucket))
		err := b.Delete([]byte(tasks[id-1].Name))
		return err
	})
	fmt.Printf("You have removed the \"%s\" task.\n", tasks[id-1].Name)
	return err
}

func (db *DB) getTasks() ([]*Task, error) {
	tasks := []*Task{}
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TasksBucket))
		err := b.ForEach(func(k, v []byte) error {
			task := &Task{}
			task.Name = string(k)
			if len(v) > 0 {
				task.Completed = int64(btoi(v))
			}
			tasks = append(tasks, task)
			return nil
		})
		return err
	})
	return tasks, err
}


func (db *DB) GetTasksCompleted() ([]*Task, error) {
	now := time.Now()
	tasks, err := db.getTasks()
	if err != nil {
		return nil, err
	}
	result := []*Task{}
	for _, task := range tasks {
		if task.Completed != 0 {
			t := time.Unix(task.Completed, 0)
			if (t.Day() == now.Day() && t.Month() == now.Month() && t.Year() == now.Year()) {
				result = append(result, task)
			}
		}
	}
	return result, err
}

func (db *DB) GetTasksNotCompleted() ([]*Task, error) {
	tasks, err := db.getTasks()
	if err != nil {
		return nil, err
	}
	result := []*Task{}
	for _, task := range tasks {
		if task.Completed == 0 {
			result = append(result, task)
		}
	}
	return result, err
}


func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
