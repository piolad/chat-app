db = db.getSiblingDB('admin');

// Create adminUser with root role
db.createUser({
  user: 'adminUser',
  pwd: 'adminPassword',
  roles: [{ role: 'root', db: 'admin' }]
});