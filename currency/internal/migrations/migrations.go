package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed *.sql
var migrationsFS embed.FS

// Migrator структура для применения миграций.
type Migrator struct {
	srcDriver source.Driver // Драйвер источника миграций.
}

// NewMigrator создает новый экземпляр Migrator с встроенными SQL-файлами миграций.
func NewMigrator() *Migrator {
	// Создаем новый драйвер источника миграций с встроенными SQL-файлами.
	d, err := iofs.New(migrationsFS, ".")
	if err != nil {
		panic(err)
	}
	return &Migrator{
		srcDriver: d,
	}
}

// ApplyMigrations применяет миграции к базе данных используя pgxpool.
func (m *Migrator) ApplyMigrations(pool *pgxpool.Pool) error {
	// Получаем stdlib подключение из пула pgx для совместимости с migrate
	conn := stdlib.OpenDBFromPool(pool)

	// Создаем экземпляр драйвера базы данных для PostgreSQL с pgx.
	driver, err := pgx.WithInstance(conn, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %w", err)
	}

	// Создаем новый экземпляр мигратора с использованием драйвера источника и драйвера базы данных PostgreSQL.
	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %w", err)
	}

	// Закрываем мигратор и подключение в конце работы функции.
	defer func() {
		migrator.Close()
		conn.Close()
	}()

	// Применяем миграции.
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return nil
}
