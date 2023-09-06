USE hoge_development;

CREATE TABLE users (
  id BIGINT UNSIGNED AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(11) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  PRIMARY KEY(id)
) COMMENT="ログインユーザー";

ALTER TABLE users RENAME INDEX phone TO index_users_on_phone;
ALTER TABLE users RENAME INDEX email TO index_users_on_email;