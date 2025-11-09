db = db.getSiblingDB('admin');

// Create adminUser with root role - for local development
db.createUser({
  user: 'adminUser',
  pwd: 'adminPassword',
  roles: [{ role: 'root', db: 'admin' }]
});

