#  sudo apt install mysql-server

PASSWDDB="pass"
RPASS="2wsxzaq1"

# # replace "-" with "_" for database username
# MAINDB="snippetbox"

# sudo mysql_secure_installation <<EOF

# y
# $RPASS
# $RPASS
# y
# y
# y
# y
# EOF

sudo mysql <<EOF
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '$RPASS';
FLUSH PRIVILEGES;
EOF


mysql -u root -p${RPASS} <<EOF
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE snippetbox;

CREATE TABLE snippets (
id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
title VARCHAR(100) NOT NULL,
content TEXT NOT NULL,
created DATETIME NOT NULL,
expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
'An old silent pond',
'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
'Over the wintry forest',
'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);
INSERT INTO snippets (title, content, created, expires) VALUES (
'First autumn morning',
'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
UTC_TIMESTAMP(),
DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY '${PASSWDDB}';
EOF

mysql -D snippetbox -u root -p${RPASS} <<EOF
CREATE TABLE sessions (
token CHAR(43) PRIMARY KEY,
data BLOB NOT NULL,
expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
EOF

mysql -D snippetbox -u root -p${RPASS} <<EOF
CREATE TABLE users (
id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
name VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL,
hashed_password CHAR(60) NOT NULL,
created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
EOF

mysql -D snippetbox -u web -p${PASSWDDB} <<EOF
SELECT id, title, expires FROM snippets;
SELECT token, data, expiry FROM sessions;
SELECT name, email, hashed_password, created FROM users;
EOF

mysql -u root -p${RPASS} <<EOF
CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'test_web'@'localhost';
GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_snippetbox.* TO 'test_web'@'localhost';
ALTER USER 'test_web'@'localhost' IDENTIFIED BY 'pass';
EOF