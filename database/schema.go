package database

var (
	Users_Table = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		age INTEGER,
		gender TEXT,
		firstname TEXT,
		lastname TEXT,
		email TEXT UNIQUE,
		password TEXT
	);`

	Posts_Table = `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		category TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		image_path TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	Likes_Table = `CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);`

	Unlikes_Table = `CREATE TABLE IF NOT EXISTS unlikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		post_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (post_id) REFERENCES posts(id)
	);`

	Comments_table = `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`

	UnlikesComment_Table = `CREATE TABLE IF NOT EXISTS unlikescomment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		comment_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (comment_id) REFERENCES comment(id)
	);`

	CommentLikes_table = `CREATE TABLE IF NOT EXISTS commentlikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	Categorie_table = `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);`
)
