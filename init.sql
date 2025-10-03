CREATE DATABASE IF NOT EXISTS goweb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE goweb;

-- Users
CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(32) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE,
  pass_hash VARCHAR(255) NOT NULL,
  role ENUM('user','admin') NOT NULL DEFAULT 'user',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login_at DATETIME NULL
) ENGINE=InnoDB;

-- Threads
CREATE TABLE IF NOT EXISTS threads (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  title VARCHAR(200) NOT NULL,
  body MEDIUMTEXT NOT NULL,
  is_locked TINYINT(1) NOT NULL DEFAULT 0,
  is_deleted TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Replies
CREATE TABLE IF NOT EXISTS replies (
  id INT AUTO_INCREMENT PRIMARY KEY,
  thread_id INT NOT NULL,
  user_id INT NOT NULL,
  parent_id INT NULL,
  body MEDIUMTEXT NOT NULL,
  is_deleted TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NULL,
  FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_id) REFERENCES replies(id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- Votes (รองรับทั้ง thread และ reply)
CREATE TABLE IF NOT EXISTS votes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  thread_id INT NULL,
  reply_id INT NULL,
  value TINYINT NOT NULL, -- 1 = upvote, -1 = downvote
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
  FOREIGN KEY (reply_id) REFERENCES replies(id) ON DELETE CASCADE,
  CONSTRAINT unique_vote UNIQUE (user_id, thread_id, reply_id)
) ENGINE=InnoDB;

-- Indexes
CREATE INDEX idx_threads_title ON threads(title);
CREATE FULLTEXT INDEX ft_threads_body ON threads(body);

-- Tags (many-to-many with threads)
CREATE TABLE IF NOT EXISTS tags (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS thread_tags (
  thread_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (thread_id, tag_id),
  FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- Useful indexes
CREATE INDEX idx_thread_tags_tag_id ON thread_tags(tag_id);
CREATE INDEX idx_replies_thread_id ON replies(thread_id);

-- Sample seed data (development only)
-- Insert users for development
INSERT INTO users (username, email, pass_hash, role, created_at)
SELECT 'alice','alice@example.com','dev-hash','user',NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'alice');

INSERT INTO users (username, email, pass_hash, role, created_at)
SELECT 'bob','bob@example.com','dev-hash','user',NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'bob');

-- Insert a sample thread
INSERT INTO threads (user_id, title, body, is_locked, is_deleted, created_at)
SELECT u.id, 'Welcome to the mini webboard', 'This is a seeded thread to help development.', 0, 0, NOW()
FROM users u WHERE u.username = 'alice' AND NOT EXISTS (SELECT 1 FROM threads WHERE title = 'Welcome to the mini webboard');

-- Insert a reply to the sample thread
INSERT INTO replies (thread_id, user_id, parent_id, body, is_deleted, created_at)
SELECT t.id, u.id, NULL, 'This is a reply from Bob.', 0, NOW()
FROM threads t CROSS JOIN users u
WHERE t.title = 'Welcome to the mini webboard' AND u.username = 'bob'
  AND NOT EXISTS (SELECT 1 FROM replies r WHERE r.body = 'This is a reply from Bob.');

-- Insert a sample tag and link it
INSERT INTO tags (name) SELECT 'introduction' WHERE NOT EXISTS (SELECT 1 FROM tags WHERE name = 'introduction');
INSERT INTO thread_tags (thread_id, tag_id)
SELECT t.id, tg.id FROM threads t JOIN tags tg ON tg.name = 'introduction' WHERE t.title = 'Welcome to the mini webboard'
  AND NOT EXISTS (SELECT 1 FROM thread_tags tt WHERE tt.thread_id = t.id AND tt.tag_id = tg.id);

-- A sample vote (Alice upvotes the seeded thread)
INSERT INTO votes (user_id, thread_id, reply_id, value, created_at)
SELECT u.id, t.id, NULL, 1, NOW() FROM users u JOIN threads t ON t.title = 'Welcome to the mini webboard'
WHERE u.username = 'alice' AND NOT EXISTS (SELECT 1 FROM votes v WHERE v.user_id = u.id AND v.thread_id = t.id);

-- End of init.sql

-- Reports: user-submitted reports for threads or users
  CREATE TABLE IF NOT EXISTS reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reporter_id INT NULL,
    kind ENUM('thread','user') NOT NULL,
    target_id INT NOT NULL,
    reason TEXT,
    status ENUM('open','resolved','dismissed') NOT NULL DEFAULT 'open',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_by INT NULL,
    resolved_at DATETIME NULL,
    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL
  ) ENGINE=InnoDB;



  CREATE TABLE IF NOT EXISTS password_resets (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(128) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  used BOOLEAN DEFAULT false
);

CREATE INDEX idx_password_resets_user ON password_resets(user_id);
CREATE INDEX idx_password_resets_token_hash ON password_resets(token_hash);