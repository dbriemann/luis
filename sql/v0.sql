CREATE TABLE IF NOT EXISTS users (
	id 			INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 	INTEGER,
	updated_at 	INTEGER,
	email 		TEXT,
	secret 		TEXT,
	name 			TEXT,
	is_admin 	INTEGER
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
ON users(email);

CREATE TABLE IF NOT EXISTS files (
	id 				INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 		INTEGER,
	updated_at 		INTEGER,
	name 				TEXT,
	thumb 			TEXT,
	type 				INTEGER,
	title 			TEXT,
	description 	TEXT,
	owner_id 		INTEGER NOT NULL,
	FOREIGN KEY(owner_id)
	REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id 			INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 	INTEGER,
	updated_at 	INTEGER,
	text 			TEXT,
	file_id 		INTEGER NOT NULL,
	FOREIGN KEY(file_id) 
	REFERENCES files(id)
);

CREATE TABLE IF NOT EXISTS stars (
	id 	 		INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 	INTEGER,
	updated_at 	INTEGER,
	file_id 		INTEGER NOT NULL,
	FOREIGN KEY(file_id)
	REFERENCES files(id)
);

CREATE TABLE IF NOT EXISTS permissions (
	id 				INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 		INTEGER,
	updated_at 		INTEGER,
	type 				INTEGER, 
	collection_id 	INTEGER NOT NULL,
	user_id 			INTEGER NOT NULL,
	FOREIGN KEY(collection_id)
	REFERENCES collections (id),
	FOREIGN KEY(user_id) 
	REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS collections (
	id 				INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at 		INTEGER,
	updated_at 		INTEGER,
	name 				TEXT,
	description 	TEXT,
	owner_id 		INTEGER NOT NULL,
	cover_id 		INTEGER NOT NULL,
	FOREIGN KEY(owner_id)
	REFERENCES users(id),
	FOREIGN KEY(cover_id) 
	REFERENCES files(id)
);

CREATE TABLE IF NOT EXISTS collections_files (
	id 				INTEGER PRIMARY KEY AUTOINCREMENT,
	collection_id 	INTEGER NOT NULL,
	file_id 			INTEGER NOT NULL,
	FOREIGN KEY(collection_id)
	REFERENCES collections(id),
	FOREIGN KEY(file_id)
	REFERENCES files(id)
);
