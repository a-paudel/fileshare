#!/bin/bash
source .env
goose -dir migrations sqlite3 $DB_URL up 