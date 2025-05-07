package domain

import "time"

type Article struct {
	ID          int64   `db:"id"`
	Header      string  `db:"header"`
	SubHeader   string  `db:"sub_header"`
	Description string  `db:"description"`
	Image       *string `db:"image"`
	Categories  []Category
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
