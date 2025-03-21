package database

import "os"

var dbFile string = os.Getenv("TODO_DBFILE")
var limit int = 50
var schemeSql string = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment VARCHAR(128) NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
    );
    CREATE INDEX scheduler_date ON scheduler (date)`
