package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Yandex-Practicum/final/pkg/bootstrap"
	"github.com/Yandex-Practicum/final/pkg/dto"
	dateService "github.com/Yandex-Practicum/final/pkg/services/date"
)

func AddTask(task *dto.Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`
	res, err := bootstrap.Db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*dto.Task, error) {
	rows, err := bootstrap.Db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date DESC LIMIT :limit",
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return makeList(rows)
}

func SearchTasks(data string, limit int) ([]*dto.Task, error) {
	var rows *sql.Rows

	searchDate, err := time.Parse("02.01.2006", data)

	if err != nil {
		rows, err = bootstrap.Db.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE concat('%', :search, '%')
			OR comment LIKE concat('%', :search, '%')
			ORDER BY date LIMIT :limit`,
			sql.Named("search", data),
			sql.Named("limit", limit),
		)
	} else {
		rows, err = bootstrap.Db.Query(`
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = :date LIMIT :limit`,
			sql.Named("date", searchDate.Format(dateService.DateFormat)),
			sql.Named("limit", limit),
		)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return makeList(rows)
}

func makeList(rows *sql.Rows) ([]*dto.Task, error) {
	tasks := make([]*dto.Task, 0)

	for rows.Next() {
		task := dto.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id int) (*dto.Task, error) {
	task := dto.Task{}
	err := bootstrap.Db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(task *dto.Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := bootstrap.Db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id int) error {
	_, err := bootstrap.Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}
