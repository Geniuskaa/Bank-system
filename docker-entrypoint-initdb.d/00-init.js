db.createUser({
    user: "app",
    pwd: "pass",
    roles: [
        {
            role: "readWrite",
            db: "db",
        }
    ]
});

db.createCollection('orders');

db.createCollection('films');

// Example of collection creation
db.createCollection('user_payments', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['login', 'payments_transfers'],
            properties: {
                login: {
                    bsonType: 'string',
                    user: true
                },
                payments_transfers: {
                    bsonType: 'array',
                    required: ['element'],
                    properties: {
                        element: {
                            bsonType: 'object',
                            required: ['link_on_icon', 'description', 'link_on_web_site']
                        }
                    }
                }

            }
        }
    }
});

