@echo off
echo === ТЕСТИРУЕМ API ===
echo.

echo 1. СОЗДАЕМ ЗАДАЧУ:
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"name\": \"Тестовая задача\"}"
echo.
echo.

echo 2. ПОЛУЧАЕМ ВСЕ ЗАДАЧИ:
curl http://localhost:8080/tasks
echo.
echo.

echo 3. ОТМЕЧАЕМ КАК ВЫПОЛНЕННУЮ (ID 1):
curl -X PUT http://localhost:8080/tasks/1
echo.
echo.

echo 4. УДАЛЯЕМ ЗАДАЧУ:
curl -X DELETE http://localhost:8080/tasks/1
echo.
echo.

echo 5. ПРОВЕРЯЕМ ЧТО УДАЛИЛОСЬ:
curl http://localhost:8080/tasks
echo.
echo.

echo === ТЕСТ ОШИБОК ===
echo.
echo 6. ПЫТАЕМСЯ УДАЛИТЬ НЕСУЩЕСТВУЮЩУЮ:
curl -X DELETE http://localhost:8080/tasks/999
echo.
echo.

echo 7. СОЗДАТЬ С ПУСТЫМ ИМЕНЕМ:
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"name\": \"\"}"
echo.
echo.

echo  ТЕСТ ЗАВЕРШЕН!
pause