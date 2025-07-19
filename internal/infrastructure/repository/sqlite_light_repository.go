package repository

import (
	"MrAlbertX/server/internal/core/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // O driver do SQLite, importado com _
	"io/ioutil"
	"log"
)

type SQLiteLightRepository struct {
	db *sql.DB
}

func NewSQLiteLightRepository(dbPath string) *SQLiteLightRepository {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := &SQLiteLightRepository{db: db}
	repo.runMigrations()

	return repo
}

func (r *SQLiteLightRepository) runMigrations() {
	content, err := ioutil.ReadFile("./migrations/001_create_lights_table.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	if _, err := r.db.Exec(string(content)); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}
}

func (r *SQLiteLightRepository) Save(light models.Light) error {
	stmt, err := r.db.Prepare("INSERT INTO lights(id, is_on) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(light.ID, light.IsOn)
	return err
}

func (r *SQLiteLightRepository) Update(light models.Light) error {
	stmt, err := r.db.Prepare("UPDATE lights SET is_on = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(light.IsOn, light.ID)
	return err
}

func (r *SQLiteLightRepository) Delete(id string) error {
	stmt, err := r.db.Prepare("DELETE FROM lights WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (r *SQLiteLightRepository) FindByID(id string) (*models.Light, error) {
	stmt, err := r.db.Prepare("SELECT id, is_on FROM lights WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var light models.Light
	err = stmt.QueryRow(id).Scan(&light.ID, &light.IsOn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Não é um erro, apenas não encontrou
		}
		return nil, err
	}
	return &light, nil
}

func (r *SQLiteLightRepository) GetAll() ([]models.Light, error) {
	rows, err := r.db.Query("SELECT id, is_on FROM lights")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lights []models.Light
	for rows.Next() {
		var light models.Light
		if err := rows.Scan(&light.ID, &light.IsOn); err != nil {
			return nil, err
		}
		lights = append(lights, light)
	}
	return lights, nil
}