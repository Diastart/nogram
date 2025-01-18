#!/bin/bash

# Clear screen
clear

# Print header (using simpler characters to avoid encoding issues)
echo -e "\033[32m======================================"
echo -e "          NOGRAM Database Reset          "
echo -e "======================================\033[0m"

# Navigate to the database directory
cd service/database || {
    echo -e "\033[31mError: Failed to locate database directory\033[0m"
    exit 1
}

# Check if database exists
if [ ! -f wasa.db ]; then
    echo -e "\033[31mError: Database file not found\033[0m"
    exit 1
fi

echo -e "\033[33mResetting database...\033[0m"

# Execute SQLite commands with proper EOF handling
sqlite3 wasa.db << 'EndOfFile'
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;
DELETE FROM Messages;
DELETE FROM GroupsMembers;
DELETE FROM Conversations;
DELETE FROM Dialogs;
DELETE FROM Groups;
DELETE FROM Users;
DELETE FROM sqlite_sequence;
COMMIT;
PRAGMA foreign_keys = ON;
EndOfFile

if [ $? -eq 0 ]; then
    echo -e "\033[32m======================================"
    echo -e "      Database Reset Successfully       "
    echo -e "======================================\033[0m"
else
    echo -e "\033[31m======================================"
    echo -e "         Failed to Reset Database       "
    echo -e "======================================\033[0m"
    exit 1
fi