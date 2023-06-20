package calculator

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"webScraping/internal/db"
)

type Data struct {
	ID            int
	EnteredValue  string
	Result        string
	RecordingDate string
}

// calculator
func Calculator(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"./ui/html/page/calculator.page.html",
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
	rows, err := db.Db.Query("SELECT id, entered_value, result, recording_date FROM calculator ORDER BY id DESC LIMIT 10")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	// Создаем срез для хранения всех объектов Data
	var dataRecords []Data

	// Итерируем по строкам и заполняем структуру данными
	for rows.Next() {
		var dataRecord Data

		err = rows.Scan(&dataRecord.ID, &dataRecord.EnteredValue, &dataRecord.Result, &dataRecord.RecordingDate)
		if err != nil {
			log.Fatal(err)
		}

		// Добавляем объект Data в срез
		dataRecords = append(dataRecords, dataRecord)
	}

	err = ts.Execute(w, dataRecords)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере2", 500)
		return
	}
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем, что метод запроса является POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Используйте r.FormValue для получения значений полей формы
	data := r.FormValue("data")
	result := parseTheReceivedValue(data) //Приводим строку в числовой тип

	// Обрабатываем полученные данные (например, сохраняем в базу данных)
	//Вызов функции подключение к БД
	_, err := db.ConnectionDb()
	if err != nil {
		return
	}

	//Если веденные данные больше 10 символов обрезать
	if len(data) > 10 {
		data = data[:10]
	}

	//Добавляем данные в БД
	_, err = db.Db.Exec("INSERT INTO calculator (entered_value, result) VALUES (?, ?)", data, result)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/calculator", http.StatusSeeOther)
}

// Распарить полученное значение
func parseTheReceivedValue(dataParam string) float64 {
	var arithmeticallyParameter string
	for _, v := range dataParam {
		switch string(v) {
		case "+":
			arithmeticallyParameter = "+"
		case "-":
			arithmeticallyParameter = "-"
		case "*":
			arithmeticallyParameter = "*"
		case "/":
			arithmeticallyParameter = "/"
		case "%":
			arithmeticallyParameter = "%"
		}
	}

	strSplit := strings.Split(dataParam, arithmeticallyParameter)
	var result float64

	firstNumber, _ := strconv.ParseFloat(strSplit[0], 64)
	secondNumber, _ := strconv.ParseFloat(strSplit[1], 64)

	switch arithmeticallyParameter {
	case "+":
		result = firstNumber + secondNumber
	case "-":
		result = firstNumber - secondNumber
	case "*":
		result = firstNumber * secondNumber
	case "/":
		result = firstNumber / secondNumber
	case "%":
		result = float64(int(firstNumber) % int(secondNumber))
	}
	return result
}

func DeleteEntry(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Fprintf(w, "id: %d", id)

	//Вызов функции подключение к БД
	_, err = db.ConnectionDb()
	if err != nil {
		return
	}

	_, err = db.Db.Exec("DELETE FROM calculator WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/calculator", http.StatusSeeOther)
}
