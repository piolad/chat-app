db = db.getSiblingDB('admin');

// Create adminUser with root role
db.createUser({
  user: 'adminUser',
  pwd: 'adminPassword',
  roles: [{ role: 'root', db: 'admin' }]
});

// db = db.getSiblingDB('message-db');

// // Define the schema for the messages collection
// const messageSchema = {
//   $jsonSchema: {
//       bsonType: "object",
//       required: ["message", "timestamp", "status"],
//       properties: {
//           message: {
//               bsonType: "string",
//               description: "must be a string and is required"
//           },
//           timestamp: {
//               bsonType: "date",
//               description: "must be a date and is required"
//           },
//           status: {
//               enum: ["unread", "read", "deleted"],
//               description: "can only be one of the enum values and is required"
//           }
//       }
//   }
// };

// // Create messages collection with schema validation
// db.createCollection("messages", {
//   validator: messageSchema
// });

// // Insert example data into the messages collection
// db.messages.insertMany([
//   {
//     message: "Hello, how are you?",
//     timestamp: new Date("2024-06-01T08:00:00Z"),
//     status: "unread"
//   },
//   {
//     message: "I'm doing well, thank you!",
//     timestamp: new Date("2024-06-01T08:15:00Z"),
//     status: "read"
//   },
//   {
//     message: "Would you like to grab lunch later?",
//     timestamp: new Date("2024-06-01T08:30:00Z"),
//     status: "unread"
//   },
//   {
//     message: "Sure, let's meet at 12:30 PM.",
//     timestamp: new Date("2024-06-01T08:45:00Z"),
//     status: "unread"
//   }
// ]);