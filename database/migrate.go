package database

// Migrate creates the necessary tables in the database
func Migrate() error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS estates (
		id UUID PRIMARY KEY,
		width INT NOT NULL,
		length INT NOT NULL
	);`)

	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS trees (
		id UUID PRIMARY KEY,
		estate_id UUID REFERENCES estates(id),
		x INT NOT NULL,
		y INT NOT NULL,
		height INT NOT NULL
	);`)

	return err
}
