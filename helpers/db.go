package helpers

import "github.com/jmoiron/sqlx"

// DropAndCreateTable удаляет таблицу, если она уже существует и создает заново
func DropAndCreateTable(schema string, tableName string, db *sqlx.DB) (bool, error) {
	var err error
	var rows *sqlx.Rows
	// Проверяем нет ли такой таблицы в базе
	rows, err = db.Queryx("SELECT to_regclass('" + tableName + "');")
	if err != nil {
		//fmt.Println("Error on check table '"+tableName+"':", err)
		return false, err
	}
	defer rows.Close()

	// И если есть удаляем
	rowsCount := 0
	for rows.Next() {
		rowsCount++
	}

	if rowsCount > 0 {
		_, err = db.Exec("DROP TABLE IF EXISTS " + tableName + ";")
		if err != nil {
			//fmt.Println("Error on drop table '"+tableName+"':", err)
			return false, err
		}
	}

	// Создаем таблицу
	_, err = db.Exec(schema)
	if err != nil {
		//fmt.Println("Error on create table '"+tableName+"':", err)
		return false, err
	}

	return true, nil
}
