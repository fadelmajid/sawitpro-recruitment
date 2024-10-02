package repositories

import (
    "database/sql"
    "sawitpro-recruitment/models"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "fmt"
)

// TreeRepository defines the methods for tree-related database operations.
type TreeRepository interface {
    AddTreeToEstate(tree *models.Tree) error
    GetTreeByCoordinates(estateID uuid.UUID, x, y int) (*models.Tree, error)
    GetTreesByEstateID(estateID uuid.UUID) (map[string]int, error)
}

// treeRepository is the concrete implementation of the TreeRepository interface.
type treeRepository struct {
    db *sql.DB
}

// NewTreeRepository returns a new instance of treeRepository.
func NewTreeRepository(db *sql.DB) TreeRepository {
    return &treeRepository{
        db: db,
    }
}

// AddTreeToEstate inserts a new tree into the specified estate.
func (r *treeRepository) AddTreeToEstate(tree *models.Tree) error {
    logrus.Infof("Adding tree with ID: %v to estate ID: %v", tree.ID, tree.EstateID)
    _, err := r.db.Exec("INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5)", tree.ID, tree.EstateID, tree.X, tree.Y, tree.Height)
    if err != nil {
        logrus.Errorf("Failed to add tree with ID %v to estate ID %v: %v", tree.ID, tree.EstateID, err)
    }
    return err
}

// GetTreeByCoordinates retrieves a tree by its coordinates in a specific estate.
func (r *treeRepository) GetTreeByCoordinates(estateID uuid.UUID, x, y int) (*models.Tree, error) {
    logrus.Infof("Retrieving tree at coordinates (%d, %d) for estate ID: %v", x, y, estateID)
    tree := &models.Tree{}
    err := r.db.QueryRow("SELECT id, x, y, height FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3", estateID, x, y).Scan(&tree.ID, &tree.X, &tree.Y, &tree.Height)
    if err != nil {
        if err == sql.ErrNoRows {
            logrus.Warnf("No tree found at coordinates (%d, %d) for estate ID: %v", x, y, estateID)
            return nil, nil
        }
        logrus.Errorf("Failed to retrieve tree at coordinates (%d, %d) for estate ID %v: %v", x, y, estateID, err)
        return nil, err
    }
    logrus.Infof("Tree retrieved successfully at coordinates (%d, %d) for estate ID: %v", x, y, estateID)
    return tree, nil
}

// GetTreesByEstateID retrieves all trees associated with a specific estate.
func (r *treeRepository) GetTreesByEstateID(estateID uuid.UUID) (map[string]int, error) {
    logrus.Infof("Retrieving all trees for estate ID: %v", estateID)
    rows, err := r.db.Query("SELECT x, y, height FROM trees WHERE estate_id = $1", estateID)
    if err != nil {
        logrus.Errorf("Failed to retrieve trees for estate ID %v: %v", estateID, err)
        return nil, err
    }
    defer rows.Close()

    treeHeights := make(map[string]int)
    for rows.Next() {
        var x, y, height int
        if err := rows.Scan(&x, &y, &height); err != nil {
            logrus.Errorf("Failed to scan tree row for estate ID %v: %v", estateID, err)
            return nil, err
        }
        key := fmt.Sprintf("%d,%d", x, y)
        treeHeights[key] = height
    }
    if err := rows.Err(); err != nil {
        logrus.Errorf("Error occurred during rows iteration for estate ID %v: %v", estateID, err)
        return nil, err
    }
    logrus.Infof("All trees retrieved successfully for estate ID: %v", estateID)
    return treeHeights, nil
}