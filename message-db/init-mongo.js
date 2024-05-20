// init-mongo.js

db = db.getSiblingDB('admin');

db.createUser({
  user: 'adminUser',
  pwd: 'adminPassword',
  roles: [{ role: 'root', db: 'admin' }]
});