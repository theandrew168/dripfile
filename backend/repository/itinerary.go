package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/domain"
)

// ensure ItineraryRepository interface is satisfied
var _ ItineraryRepository = (*PostgresItineraryRepository)(nil)

type ItineraryRepository interface {
	Create(itinerary *domain.Itinerary) error
	List() ([]*domain.Itinerary, error)
	Read(id uuid.UUID) (*domain.Itinerary, error)
	Delete(itinerary *domain.Itinerary) error
}

type Itinerary struct {
	ID uuid.UUID `db:"id"`

	Pattern        string    `db:"pattern"`
	FromLocationID uuid.UUID `db:"from_location_id"`
	ToLocationID   uuid.UUID `db:"to_location_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostgresItineraryRepository struct {
	conn database.Conn
}

func NewPostgresItineraryRepository(conn database.Conn) *PostgresItineraryRepository {
	repo := PostgresItineraryRepository{
		conn: conn,
	}
	return &repo
}

func (repo *PostgresItineraryRepository) marshal(itinerary *domain.Itinerary) (Itinerary, error) {
	row := Itinerary{
		ID: itinerary.ID(),

		Pattern:        itinerary.Pattern(),
		FromLocationID: itinerary.FromLocationID(),
		ToLocationID:   itinerary.ToLocationID(),

		CreatedAt: itinerary.CreatedAt(),
		UpdatedAt: itinerary.UpdatedAt(),
	}
	return row, nil
}

func (repo *PostgresItineraryRepository) unmarshal(row Itinerary) (*domain.Itinerary, error) {
	itinerary := domain.LoadItinerary(
		row.ID,
		row.FromLocationID,
		row.ToLocationID,
		row.Pattern,
		row.CreatedAt,
		row.UpdatedAt,
	)
	return itinerary, nil
}

func (repo *PostgresItineraryRepository) Create(itinerary *domain.Itinerary) error {
	stmt := `
		INSERT INTO itinerary
			(id, from_location_id, to_location_id, pattern, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)`

	row, err := repo.marshal(itinerary)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.FromLocationID,
		row.ToLocationID,
		row.Pattern,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = repo.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	return nil
}

func (repo *PostgresItineraryRepository) List() ([]*domain.Itinerary, error) {
	stmt := `
		SELECT
			id,
			from_location_id,
			to_location_id,
			pattern,
			created_at,
			updated_at
		FROM itinerary
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	itineraryRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Itinerary])
	if err != nil {
		return nil, checkListError(err)
	}

	var itineraries []*domain.Itinerary
	for _, row := range itineraryRows {
		itinerary, err := repo.unmarshal(row)
		if err != nil {
			return nil, err
		}

		itineraries = append(itineraries, itinerary)
	}

	return itineraries, nil
}

func (repo *PostgresItineraryRepository) Read(id uuid.UUID) (*domain.Itinerary, error) {
	stmt := `
		SELECT
			id,
			from_location_id,
			to_location_id,
			pattern,
			created_at,
			updated_at
		FROM itinerary
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Itinerary])
	if err != nil {
		return nil, checkReadError(err)
	}

	return repo.unmarshal(row)
}

func (repo *PostgresItineraryRepository) Delete(itinerary *domain.Itinerary) error {
	stmt := `
		DELETE FROM itinerary
		WHERE id = $1
		RETURNING id`

	err := itinerary.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, itinerary.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return checkDeleteError(err)
	}

	return nil
}
