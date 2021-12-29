package storage

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Post encapsulates one item from rss feed
type Post struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}

// DB encapsulates database connection
type DB struct {
	// pool is thread-safe so we don't need mutex
	pool *pgxpool.Pool
}

// New returns new instanse of database connection
func New() (*DB, error) {
	conn := os.Getenv("rssdb")
	if conn == "" {
		return nil, errors.New("no PG connection string")
	}
	pool, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// Posts fetches last posts from database
// number of posts is limieted by lim
func (db *DB) Posts(lim int) ([]Post, error) {
	if lim == 0 {
		lim = 10
	}
	rows, err := db.pool.Query(context.Background(), `
		SELECT
			id,
			title,
			content,
			pub_time,
			link
		FROM
			devbase.rss_aggr.posts
		ORDER BY
			pub_time DESC
		LIMIT $1;
		`,
		lim)
	if err != nil {
		return nil, err
	}
	var posts []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// Save post inserts post to database
// if it's not already there
func (db *DB) SavePost(p Post) (int, error) {
	var id int
	// If post is already in database - return id
	err := db.pool.QueryRow(context.Background(), `
		INSERT INTO devbase.rss_aggr.posts (title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT post_link DO UPDATE SET link = EXCLUDED.link
		RETURNING id;
		`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
