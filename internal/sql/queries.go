package sql

const (
	createExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`

	createTable = `CREATE TABLE IF NOT EXISTS keyring (
		Id UUID NOT NULL,
		Userid UUID NOT NULL,
		Url TEXT,
		Username TEXT,
		Sitename TEXT,
		Sitepassword TEXT,
		Folder TEXT,
		Notes TEXT,
		Favorite BOOLEAN,
		PRIMARY KEY(ID)
	)`

	selectKeyEntry = `SELECT * FROM keyring WHERE Userid == $1 ORDER BY Id LIMIT $2 OFFSET $3`
	insertKeyEntry = `INSERT INTO keyring(
		Id, Userid, Url, Username, Sitename, Sitepassword, Folder, Notes, Favorite) 
		VALUES(uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8)`

	updateKeyEntry = `UPDATE keyring
	SET Url = $1, Username = $2, Sitename = $3, Sitepassword = $4, Folder = $5, Notes = $6, Favorite = $7
	WHERE Id == $8 AND Userid == $9`
)
