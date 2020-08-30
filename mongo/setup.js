// Create database user
db.createUser({
  user: "user",
  pwd: "123456",
  roles: [
    {
      role: "readWrite",
      db: "app",
    },
  ],
});

// Create collection
db.createCollection("accounts");

const users = Array(5)
  .fill({})
  .map((_, index) => {
    userID = index + 1;
    return {
      id: userID.toString(),
      name: "user" + userID,
      credential: "auth" + userID,
      transactions: [
        {
          from: "system",
          to: userID.toString(),
          amount: 1000,
          date: new Date().toISOString(),
          message: "Initial funds",
        },
      ],
    };
  });

db.accounts.insertMany(users);
