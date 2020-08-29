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
      credential: Math.random().toString(36).substring(6),
      transactions: [
        {
          from: "0", // 0 here as stand in for the system
          to: userID.toString(),
          amount: 1000,
        },
      ],
    };
  });

db.accounts.insertMany(users);
