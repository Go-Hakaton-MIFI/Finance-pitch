package article

import "time"

type RawArticle struct {
	ID               int64     `db:"id"`
	Header           string    `db:"header"`
	SubHeader        string    `db:"sub_header"`
	Description      string    `db:"description"`
	ShortDescription string    `db:"short_description"`
	SpecialOffer     string    `db:"special_offer"`
	PartnerLink      string    `db:"partner_link"`
	Image            *string   `db:"image"`
	CategoryName     *string   `db:"category_name"`
	CategoryId       *int64    `db:"category_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
