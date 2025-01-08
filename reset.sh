#!/bin/bash

# Navigate to the backend directory
cd backend || { echo "Failed to navigate to backend directory"; exit 1; }

# Open sqlite3 with the database and run the commands
sqlite3 wasa.db <<EOF
PRAGMA foreign_keys = OFF;
DELETE FROM Messages;
DELETE FROM GroupsMembers;
DELETE FROM Conversations;
DELETE FROM Dialogs;
DELETE FROM \`Groups\`;
DELETE FROM Users;
DELETE FROM sqlite_sequence;
PRAGMA foreign_keys = ON;
EOF

# Exit message to confirm execution
echo "Database cleanup commands executed successfully."

