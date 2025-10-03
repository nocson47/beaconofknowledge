// mongo-init.js
// Run with mongo <db> mongo-init.js or mount into docker-entrypoint-initdb.d for mongo container

const dbName = process.env.MONGO_DBNAME || 'goweb';

db = db.getSiblingDB(dbName);

print('Creating indexes for reports and audit/debug collections...');

db.reports.createIndex({ kind: 1, target_id: 1 });
db.reports.createIndex({ status: 1, created_at: -1 });
db.reports.createIndex({ created_at: -1 });

db.audit.createIndex({ ts: -1 });
db.audit.createIndex({ actor_id: 1 });

// TTL for debug collection: 30 days
db.debug.createIndex({ ts: 1 }, { expireAfterSeconds: 60 * 60 * 24 * 30 });

print('Mongo init complete');
