package repositories

import (
    "database/sql"
    "errors"
    "testing"
    "sawitpro-recruitment/models"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestEstateRepository_CreateEstate(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estate := &models.Estate{
        ID:     uuid.New(),
        Width:  100,
        Length: 200,
    }

    mock.ExpectExec(`INSERT INTO estates \(id, width, length\) VALUES \(\$1, \$2, \$3\)`).
        WithArgs(estate.ID, estate.Width, estate.Length).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.CreateEstate(estate)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_CreateEstate_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estate := &models.Estate{
        ID:     uuid.New(),
        Width:  100,
        Length: 200,
    }

    mock.ExpectExec(`INSERT INTO estates \(id, width, length\) VALUES \(\$1, \$2, \$3\)`).
        WithArgs(estate.ID, estate.Width, estate.Length).
        WillReturnError(errors.New("insert error"))

    err = repo.CreateEstate(estate)
    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_GetEstateByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estateID := uuid.New()
    expectedEstate := &models.Estate{
        ID:     estateID,
        Width:  100,
        Length: 200,
    }

    rows := sqlmock.NewRows([]string{"id", "width", "length"}).
        AddRow(expectedEstate.ID, expectedEstate.Width, expectedEstate.Length)

    mock.ExpectQuery(`SELECT id, width, length FROM estates WHERE id = \$1`).
        WithArgs(estateID).
        WillReturnRows(rows)

    estate, err := repo.GetEstateByID(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedEstate, estate)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_GetEstateByID_NoRows(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estateID := uuid.New()

    mock.ExpectQuery(`SELECT id, width, length FROM estates WHERE id = \$1`).
        WithArgs(estateID).
        WillReturnError(sql.ErrNoRows)

    estate, err := repo.GetEstateByID(estateID)
    assert.NoError(t, err)
    assert.Nil(t, estate)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_GetEstateByID_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estateID := uuid.New()

    mock.ExpectQuery(`SELECT id, width, length FROM estates WHERE id = \$1`).
        WithArgs(estateID).
        WillReturnError(errors.New("query error"))

    estate, err := repo.GetEstateByID(estateID)
    assert.Error(t, err)
    assert.Nil(t, estate)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_GetEstateStats(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estateID := uuid.New()
    expectedCount, expectedMax, expectedMin, expectedMedian := 10, 50, 5, 25

    rows := sqlmock.NewRows([]string{"count", "max", "min", "median"}).
        AddRow(expectedCount, expectedMax, expectedMin, expectedMedian)

    query := `
        SELECT 
            COUNT\(\*\), 
            MAX\(height\), 
            MIN\(height\), 
            PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY height\) 
        FROM trees 
        WHERE estate_id = \$1
    `

    mock.ExpectQuery(query).
        WithArgs(estateID).
        WillReturnRows(rows)

    count, max, min, median, err := repo.GetEstateStats(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedCount, count)
    assert.Equal(t, expectedMax, max)
    assert.Equal(t, expectedMin, min)
    assert.Equal(t, expectedMedian, median)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEstateRepository_GetEstateStats_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewEstateRepository(db)

    estateID := uuid.New()

    query := `
        SELECT 
            COUNT\(\*\), 
            MAX\(height\), 
            MIN\(height\), 
            PERCENTILE_CONT\(0.5\) WITHIN GROUP \(ORDER BY height\) 
        FROM trees 
        WHERE estate_id = \$1
    `

    mock.ExpectQuery(query).
        WithArgs(estateID).
        WillReturnError(errors.New("query error"))

    count, max, min, median, err := repo.GetEstateStats(estateID)
    assert.Error(t, err)
    assert.Equal(t, 0, count)
    assert.Equal(t, 0, max)
    assert.Equal(t, 0, min)
    assert.Equal(t, 0, median)
    assert.NoError(t, mock.ExpectationsWereMet())
}