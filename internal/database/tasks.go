package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

//СТРУКТУТУРА
/* func ConnectDB() *sql.DB {
	connStr := "user=postgres password=postgres123 dbname=todo_app sslmode=disable  "
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Успешное подключение PostgreSQL!")
	return db
}*/
func ConnectDB() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv ("DB_USER", "postgres")
	password:= getEnv("DB_PASSWORD", "postgres123")
	 dbname := getEnv("DB_NAME", "todo_app")    
	  
	 connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	 )
	db,err := sql.Open("postgres", connStr)
	if err != nil{
		log.Fatal(err)
	}
	err = db.Ping ()
	if err != nil{
		log.Fatal()
	}
	fmt.Println("✅ Подключено к PostgreSQL!")
	   createTable(db) 
return db
}

func createTable(db *sql.DB) {
    query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        done BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

    _, err := db.Exec(query)
    if err != nil {
        log.Printf("⚠️ Ошибка создания таблицы: %v", err)
    } else {
        fmt.Println("✅ Таблица 'tasks' создана/проверена")
    }
}


func getEnv (key, defaultValue string) string{
	if value, exists := os.LookupEnv(key); exists{
		return  value
	}
	return defaultValue
}

// просмотр задач
func GetAllTasks(db *sql.DB) ([]Task, error) {
	var tasks []Task // СЛАЙС, СРЕЗ

	query := "SELECT id,name, done, created_at FROM tasks ORDER BY id"
	rows, err := db.Query(query) // чтение строк
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Name, &task.Done, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// добвавление задачи -
func AddTask(db *sql.DB, name string) (int, error) {
	var id int
	query := "INSERT INTO tasks (name) VALUES ($1) RETURNING id"
	err := db.QueryRow(query, name).Scan(&id) // чтение строки
	if err != nil {
		return 0, err
	}
	return id, nil
}

// изменение статуса
func CompleteTask(db *sql.DB, id int) error {
	query := "UPDATE  tasks SET done = true WHERE id =$1 "
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}
	return nil
}

// удаление задачи
func DeleteTask(db *sql.DB, id int) error {
	query := " DELETE FROM tasks  WHERE id =$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
