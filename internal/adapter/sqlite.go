package adapter

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Gift struct {
	ID        int `json:"id"`
	Name      string
	Price     float64 `json:"price"`
	Reserved  bool    `json:"reserved"`
	Category  string
	Buyers    int
	MaxBuyers int
	Image     string
	Link      string
	QRCode    string `json:"qrCode"`
}

type SqliteClient struct {
	db *sql.DB
}

func NewSqliteClient(dbName string) (*SqliteClient, func() error) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatal(err)
	}

	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	return &SqliteClient{db}, db.Close
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS gift (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL,
		reserved BOOLEAN NOT NULL,
		category TEXT NOT NULL,
		buyers INTEGER NOT NULL,
		maxBuyers INTEGER NOT NULL,
		image TEXT,
		link TEXT,
		qrCode TEXT
	);`
	_, err := db.Exec(query)
	return err
}

func (s SqliteClient) InsertGift(g Gift) error {
	_, err := s.db.Exec(`
		INSERT INTO gift (id, name, price, reserved, category, buyers, maxBuyers, image, link, qrCode)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, g.ID, g.Name, g.Price, g.Reserved, g.Category, g.Buyers, g.MaxBuyers, g.Image, g.Link, g.QRCode)
	return err
}

func (s SqliteClient) ListGifts() ([]Gift, error) {
	rows, err := s.db.Query(`SELECT id, name, price, reserved, category, buyers, maxBuyers, image, link, qrCode FROM gift`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gifts []Gift
	for rows.Next() {
		var g Gift
		err := rows.Scan(&g.ID, &g.Name, &g.Price, &g.Reserved, &g.Category, &g.Buyers, &g.MaxBuyers, &g.Image, &g.Link, &g.QRCode)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, g)
	}
	return gifts, nil
}

func (s SqliteClient) GetGiftByID(id int) (*Gift, error) {
	row := s.db.QueryRow(`SELECT id, name, price, reserved, category, buyers, maxBuyers, image, link, qrCode FROM gift WHERE id = ?`, id)
	var g Gift
	if err := row.Scan(&g.ID, &g.Name, &g.Price, &g.Reserved, &g.Category, &g.Buyers, &g.MaxBuyers, &g.Image, &g.Link, &g.QRCode); err != nil {
		return nil, err
	}
	return &g, nil
}

func (s SqliteClient) UpdateGiftQRCode(id int, qrCode string) error {
	result, err := s.db.Exec(`UPDATE gift SET qrcode = ? WHERE id = ?`, qrCode, id)
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(rows)
	return err
}

//
// func deleteGift(db *sql.DB, id int) error {
// 	_, err := db.Exec(`DELETE FROM gift WHERE id = ?`, id)
// 	return err
// }
