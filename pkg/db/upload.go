package db

import "time"

// Upload represents an upload record stored in the database
type Upload struct {
	Username  string    `json:"username"`
	Filename  string    `json:"filename"`
	S3Url     string    `json:"s3_url"`
	Timestamp time.Time `json:"timestamp"`
}

// InsertUpload records a given upload in the database
func (c *Client) InsertUpload(u Upload) error {
	stmt :=
		`INSERT INTO uploads(username, filename, s3url, timestamp)
		 VALUES($1, $2, $3, $4)`

	u.Timestamp = time.Now()
	if _, err := c.db.Exec(stmt, u.Username, u.Filename, u.S3Url, u.Timestamp); err != nil {
		return err
	}

	return nil
}

// ListUploads retrieves the list of uploads for a given user
func (c *Client) ListUploads(username string) ([]Upload, error) {
	stmt :=
		`SELECT username, filename, s3url, timestamp
		 FROM uploads
		 WHERE username=$1`

	rows, err := c.db.Query(stmt, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Upload, 0)

	for rows.Next() {
		u := Upload{}
		if err := rows.Scan(&u.Username, &u.Filename, &u.S3Url, &u.Timestamp); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	return result, nil
}
