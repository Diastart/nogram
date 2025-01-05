CREATE TABLE Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    token TEXT NOT NULL,
    photo TEXT NOT NULL
);

CREATE TABLE Dialogs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    companion1Id INTEGER NOT NULL,
    companion2Id INTEGER NOT NULL,
    FOREIGN KEY (companion1Id) REFERENCES Users(id),
    FOREIGN KEY (companion2Id) REFERENCES Users(id)
);

CREATE TABLE `Groups` (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    groupname TEXT NOT NULL,
    photo TEXT NOT NULL
);

CREATE TABLE GroupsMembers (
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE Conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dialog_id INTEGER,
    group_id INTEGER,
    FOREIGN KEY (dialog_id) REFERENCES Dialogs(id),
    FOREIGN KEY (group_id) REFERENCES Groups(id)
);

CREATE TABLE Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    conversation_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(id),
    FOREIGN KEY (conversation_id) REFERENCES Conversations(id)
);