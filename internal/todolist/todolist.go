package todolist

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
	"webScraping/internal/db"
)

type DataToDoList struct {
	ID        int
	Name      string
	Text      string
	DateAdded string
}

// To do list
func ToDoList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/page/toDoList.page.html",
		"./ui/html/layout/base.layout.html",
		"./ui/html/partial/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере", 500)
		return
	}

	//Вызов функции подключение к БД
	_, err = db.ConnectionDb()
	if err != nil {
		log.Println(err)
	}

	//Вывод данных из БД на страницу
	rows, err := db.Db.Query("SELECT id, name, text, date_added FROM todolist ORDER BY id DESC")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// Создаем срез для хранения всех объектов Data
	var dataRecordsToDoList []DataToDoList

	// Итерируем по строкам и заполняем структуру данными
	for rows.Next() {
		var dataRecordToDoList DataToDoList

		err = rows.Scan(&dataRecordToDoList.ID, &dataRecordToDoList.Name, &dataRecordToDoList.Text, &dataRecordToDoList.DateAdded)
		if err != nil {
			log.Fatal(err)
		}

		// Добавляем объект Data в срез
		dataRecordsToDoList = append(dataRecordsToDoList, dataRecordToDoList)
	}

	// Выполнить запрос на выборку последнего ID
	var lastID int
	err = db.Db.QueryRow("SELECT max(id) FROM todolist").Scan(&lastID)
	if err != nil {
		log.Println(err)
		return
	}

	err = ts.Execute(w, map[string]interface{}{
		"Records": dataRecordsToDoList,
		"ID":      lastID + 1,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере2", 500)
		return
	}

}

func FormHandlerToDoList(w http.ResponseWriter, r *http.Request) {
	//Проверяем, что метод запроса является POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Используйте r.FormValue для получения значений полей формы
	name := r.FormValue("mainFormsToDoListName")
	text := r.FormValue("mainFormsToDoListText")

	//Вызов функции подключение к БД
	_, err := db.ConnectionDb()
	if err != nil {
		log.Println(err)
	}

	//Данные с формы добавляем в БД
	_, err = db.Db.Exec("INSERT INTO todolist (name, text) VALUES (?,?)", name, text)
	if err != nil {
		return
	}

	//Перенаправление
	http.Redirect(w, r, "/todolist", http.StatusSeeOther)

}

// Удалить запись из Списка дел
func DeleteEntryToDoList(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}

	//Вызов функции подключение к БД
	_, err = db.ConnectionDb()
	if err != nil {
		return
	}

	_, err = db.Db.Exec("DELETE FROM todolist WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/todolist", http.StatusSeeOther)
}

// Редактировать запись из Списка дел
func EditPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	message := r.URL.Query().Get("message")

	//Вызов функции подключение к БД
	_, err = db.ConnectionDb()
	if err != nil {
		log.Println(err)
	}
	_, err = db.Db.Exec("UPDATE todolist SET text = ? WHERE id = ?", message, id)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/todolist", http.StatusSeeOther)
}
