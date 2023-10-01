-- Go 工具会忽略任何名为 testdata 的目录，因此在编译应用程序时会忽略这些脚本（它也会忽略任何名称以 _ 或 .字符开头的目录或文件）。
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);
CREATE TABLE users (
   id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL,
   hashed_password CHAR(60) NOT NULL,
   created DATETIME NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
);
