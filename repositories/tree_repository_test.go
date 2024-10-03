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

func TestTreeRepository_AddTreeToEstate(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    tree := &models.Tree{
        ID:       uuid.New(),
        EstateID: uuid.New(),
        X:        10,
        Y:        20,
        Height:   30,
    }

    mock.ExpectExec("INSERT INTO trees").
        WithArgs(tree.ID, tree.EstateID, tree.X, tree.Y, tree.Height).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.AddTreeToEstate(tree)
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_AddTreeToEstate_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    tree := &models.Tree{
        ID:       uuid.New(),
        EstateID: uuid.New(),
        X:        10,
        Y:        20,
        Height:   30,
    }

    mock.ExpectExec("INSERT INTO trees").
        WithArgs(tree.ID, tree.EstateID, tree.X, tree.Y, tree.Height).
        WillReturnError(errors.New("insert error"))

    err = repo.AddTreeToEstate(tree)
    assert.Error(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_GetTreeByCoordinates(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    estateID := uuid.New()
    expectedTree := &models.Tree{
        ID:     uuid.New(),
        X:      10,
        Y:      20,
        Height: 30,
    }

    rows := sqlmock.NewRows([]string{"id", "x", "y", "height"}).
        AddRow(expectedTree.ID, expectedTree.X, expectedTree.Y, expectedTree.Height)

    mock.ExpectQuery("SELECT id, x, y, height FROM trees WHERE estate_id = \\$1 AND x = \\$2 AND y = \\$3").
        WithArgs(estateID, 10, 20).
        WillReturnRows(rows)

    tree, err := repo.GetTreeByCoordinates(estateID, 10, 20)
    assert.NoError(t, err)
    assert.Equal(t, expectedTree, tree)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_GetTreeByCoordinates_NoRows(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    estateID := uuid.New()

    mock.ExpectQuery("SELECT id, x, y, height FROM trees WHERE estate_id = \\$1 AND x = \\$2 AND y = \\$3").
        WithArgs(estateID, 10, 20).
        WillReturnError(sql.ErrNoRows)

    tree, err := repo.GetTreeByCoordinates(estateID, 10, 20)
    assert.NoError(t, err)
    assert.Nil(t, tree)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_GetTreeByCoordinates_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    estateID := uuid.New()

    mock.ExpectQuery("SELECT id, x, y, height FROM trees WHERE estate_id = \\$1 AND x = \\$2 AND y = \\$3").
        WithArgs(estateID, 10, 20).
        WillReturnError(errors.New("query error"))

    tree, err := repo.GetTreeByCoordinates(estateID, 10, 20)
    assert.Error(t, err)
    assert.Nil(t, tree)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_GetTreesByEstateID(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    estateID := uuid.New()
    rows := sqlmock.NewRows([]string{"x", "y", "height"}).
        AddRow(10, 20, 30).
        AddRow(15, 25, 35)

    mock.ExpectQuery("SELECT x, y, height FROM trees WHERE estate_id = \\$1").
        WithArgs(estateID).
        WillReturnRows(rows)

    expectedTreeHeights := map[string]int{
        "10,20": 30,
        "15,25": 35,
    }

    treeHeights, err := repo.GetTreesByEstateID(estateID)
    assert.NoError(t, err)
    assert.Equal(t, expectedTreeHeights, treeHeights)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTreeRepository_GetTreesByEstateID_Error(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTreeRepository(db)

    estateID := uuid.New()

    mock.ExpectQuery("SELECT x, y, height FROM trees WHERE estate_id = \\$1").
        WithArgs(estateID).
        WillReturnError(errors.New("query error"))

    treeHeights, err := repo.GetTreesByEstateID(estateID)
    assert.Error(t, err)
    assert.Nil(t, treeHeights)
    assert.NoError(t, mock.ExpectationsWereMet())
}