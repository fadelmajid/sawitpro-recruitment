package repositories

import (
	"database/sql"
	"sawitpro-recruitment/models"
	"github.com/google/uuid"
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
	_, err := r.db.Exec("INSERT INTO trees (id, estate_id, x, y, height) VALUES ($1, $2, $3, $4, $5)", tree.ID, tree.EstateID, tree.X, tree.Y, tree.Height)
	return err
}

// GetTreeByCoordinates retrieves a tree by its coordinates in a specific estate.
func (r *treeRepository) GetTreeByCoordinates(estateID uuid.UUID, x, y int) (*models.Tree, error) {
	tree := &models.Tree{}
	err := r.db.QueryRow("SELECT id, x, y, height FROM trees WHERE estate_id = $1 AND x = $2 AND y = $3", estateID, x, y).Scan(&tree.ID, &tree.X, &tree.Y, &tree.Height)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return tree, nil
}

// GetTreesByEstateID retrieves all trees associated with a specific estate.
func (r *treeRepository) GetTreesByEstateID(estateID uuid.UUID) (map[string]int, error) {
	rows, err := r.db.Query("SELECT x, y, height FROM trees WHERE estate_id = $1", estateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	treeHeights := make(map[string]int)
	for rows.Next() {
		var x, y, height int
		rows.Scan(&x, &y, &height)
		key := fmt.Sprintf("%d,%d", x, y)
		treeHeights[key] = height
	}
	return treeHeights, nil
}
