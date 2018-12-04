CREATE TABLE `user` (
    `username` VARCHAR(64) PRIMARY KEY,
    `password` VARCHAR NULL
);

CREATE TABLE 'createdcreds' (
    'uid' INT(64) PRIMARY KEY AUTOINCREMENT,
    'date' DATE NOT NULL,
    'username' VARCHAR(64),
    FOREIGN KEY ('username') references 'user'('username')
);