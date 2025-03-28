package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) error {
	expireFormated := time.Now().AddDate(0, 0, expires)
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, NOW(), $3)`

	_, err := m.DB.Exec(stmt, title, content, expireFormated)
	if err != nil {
		return err
	}

	return nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	var snippet Snippet

	stmt := `SELECT * FROM snippets WHERE expires > NOW() AND id = $1`

	err := m.DB.QueryRow(stmt, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	var snippets []Snippet

	stmt := `SELECT * FROM snippets WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snippet Snippet
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
