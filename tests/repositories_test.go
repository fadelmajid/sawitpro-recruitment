package tests

import (
    "database/sql"
	_ "github.com/lib/pq"
    "testing"

    "sawitpro-recruitment/models"
    "sawitpro-recruitment/repositories"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestEstateRepository_CreateEstate_Success(t *testing.T) {
    // Add import for the Postgres driver
    db, err := sql.Open("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        t.Fatalf("could not connect to the database: %v", err)
    }
    defer db.Close()

    estateRepo := repositories.NewEstateRepository(db)
    estate := &models.Estate{
        ID:     uuid.New(),
        Width:  10,
        Length: 20,
    }

    err = estateRepo.CreateEstate(estate)
    assert.NoError(t, err)

    createdEstate, err := estateRepo.GetEstateByID(estate.ID)
    assert.NoError(t, err)
    assert.NotNil(t, createdEstate)
    assert.Equal(t, estate.Width, createdEstate.Width)
    assert.Equal(t, estate.Length, createdEstate.Length)
}


func TestEstateRepository_GetEstateStats_Success(t *testing.T) {
    db, err := sql.Open("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        t.Fatalf("could not connect to the database: %v", err)
    }
    defer db.Close()

    estateRepo := repositories.NewEstateRepository(db)
    estateID := uuid.New() // Assume this ID exists and has trees
	expectedMedian := 15 // Change this based on the expected median in your database for the test

    // Capture all 5 returned values: count, max, min, median, and error
    count, max, min, median, err := estateRepo.GetEstateStats(estateID)
    
    assert.NoError(t, err)
    assert.Greater(t, count, 0)       // Assert that there are trees
    assert.LessOrEqual(t, min, max)   // Assert min <= max
    assert.Equal(t, expectedMedian, median) // Assert the median if expected
}

func TestTreeRepository_AddTreeToEstate_Success(t *testing.T) {
    // Setup in-memory database or a mock for testing
    db, err := sql.Open("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        t.Fatalf("could not connect to the database: %v", err)
    }
    defer db.Close()

    treeRepo := repositories.NewTreeRepository(db)
    tree := &models.Tree{
        ID:       uuid.New(),
        EstateID: uuid.New(), // Use a valid estate ID
        X:        1,
        Y:        1,
        Height:   15,
    }

    err = treeRepo.AddTreeToEstate(tree)
    assert.NoError(t, err)

    // Validate that the tree has been added
    createdTree, err := treeRepo.GetTreeByCoordinates(tree.EstateID, tree.X, tree.Y)
    assert.NoError(t, err)
    assert.NotNil(t, createdTree)
    assert.Equal(t, tree.Height, createdTree.Height)
}

func TestTreeRepository_GetTreesByEstateID_Success(t *testing.T) {
    // Setup in-memory database or a mock for testing
    db, err := sql.Open("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        t.Fatalf("could not connect to the database: %v", err)
    }
    defer db.Close()

    treeRepo := repositories.NewTreeRepository(db)
    estateID := uuid.New() // Assume this ID exists

    // Create some mock trees in the database first
    // Use treeRepo.AddTreeToEstate or raw SQL to insert test data

    trees, err := treeRepo.GetTreesByEstateID(estateID)
    assert.NoError(t, err)
    assert.NotNil(t, trees)
    assert.Greater(t, len(trees), 0) // Ensure there are trees
}
