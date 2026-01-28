package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"todo-web-app/internal/database"
)

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника (можно указать конкретный домен)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Обрабатываем предварительный запрос OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
func main() {

	// ОБСЛУЖИВАЕМ СТАТИЧЕСКИЕ ФАЙЛЫ 
	http.Handle("/", http.FileServer(http.Dir(".")))

	// ПОДКЛЮЧЕНИЕ К БД
	fmt.Println(" Запускаю сервер...")
	db := database.ConnectDB()
	defer db.Close()
	fmt.Println("✅ Подключено к БД")

	http.HandleFunc("/tasks", enableCORS(func(w http.ResponseWriter, r *http.Request){ // оброботчик функций(путь- адрес на который приходят запросы,анонимная функция(переменая и ее тип- интерфейс для записи ответа,переменная для указания структуры запроса))
		w.Header().Set("Content-Type", "application/json") // заголовок и тип возвращаемых данных
		switch r.Method {
		// просмотр задач
		case http.MethodGet:
			tasks, err := database.GetAllTasks(db)
			if err != nil {
				w.Write([]byte("ошибка, не удалось получить задачи "))
				return
			}
			fmt.Printf("получили %d задач из базы данных\n", len(tasks))

			response := map[string]interface{}{
				"status": "success",
				"count":  len(tasks),
				"tasks":  tasks, // ← МАССИВ СТРУКТУР!
			}

			// 3. Преобразуем в JSON
			jsonData, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(jsonData))
			w.Write([]byte(jsonData))

			//добавление задач
		case http.MethodPost:

			type CrateTaskRequest struct {
				Name string `json:"name"`
			}
			var reqData CrateTaskRequest
			err := json.NewDecoder(r.Body).Decode(&reqData)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(`{"error":"НЕверный формат"}`))
				return

			}
			if strings.TrimSpace(reqData.Name) == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "Имя задачи не может быть пустым"}`))
				return
			}
			taskID, err := database.AddTask(db, reqData.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"ошибка"}`))
				return
			}
			response := map[string]interface{}{
				"status":  "success",
				"message": "задача успешно добавлена",
				"task_id": taskID,
				"task": map[string]interface{}{
					"id":   taskID,
					"name": reqData.Name,
					"done": false,
				},
			}
			jsonData, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(jsonData))
			w.Write([]byte(jsonData))

		}

	}))

	http.HandleFunc("/tasks/", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path                    // Path - строка с путем после домена
		pathParts := strings.Split(path, "/") //strings.Split(path, "/") - функция, которая делит строку на части
		if len(pathParts) < 3 {
			http.Error(w, `{"error": "Не указан Id"}`, http.StatusBadRequest) //http.StatusBadRequest- неверный запрос
			return
		}
		taskIDStr := pathParts[2]
		taskId, err := strconv.Atoi(taskIDStr) // преобразование строки в число
		if err != nil || taskId <= 0 {
			http.Error(w, `{"error": "Неверный ID задачи"}`, http.StatusBadRequest)
			return
		}

		switch r.Method {
		// Отметить задачу как выполненную
		case http.MethodPut:
			err := database.CompleteTask(db, taskId)
			if err != nil {
				if strings.Contains(err.Error(), "не найдена") {
					http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
				} else {
					http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
				}
				return
			}

			// Успешный ответ
			response := map[string]interface{}{
				"status":  "success",
				"message": fmt.Sprintf("Задача #%d отмечена как выполненная", taskId),
				"id":      taskId,
				"done":    true,
			}
			json.NewEncoder(w).Encode(response)
		case http.MethodDelete:
			err := database.DeleteTask(db, taskId)
			if err != nil {
				if strings.Contains(err.Error(), "не найдена") {
					http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
				} else {
					http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
				}
				return
			}
			// Успешный ответ
			response := map[string]interface{}{
				"status":  "success",
				"message": fmt.Sprintf("Задача #%d удалена", taskId),
				"id":      taskId,
			}
			json.NewEncoder(w).Encode(response) // кодировщик(куда писать результат).метод кодировщика(что кодирует)

		}
	}))
	port := ":8081"
	log.Fatal(http.ListenAndServe(port, nil))
}
