// Create database user
db.createUser({
  user: "user",
  pwd: "123456",
  roles: [
    {
      role: "read",
      db: "app",
    },
  ],
});

// Create collection
db.createCollection("data");

// Create base data
