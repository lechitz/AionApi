INSERT INTO aion_api.users (name, username, password, email, roles)
VALUES
    ('user','user','$2a$10$AQjsEvApPM59JVg4XIqEluc.JtKfnyIWPLK7QEYW6CPDN9cU/OH7O','user@aionapi.dev','user'),
    ('Bruno Amaral','bruno','$2a$10$AQjsEvApPM59JVg4XIqEluc.JtKfnyIWPLK7QEYW6CPDN9cU/OH7O','bruno@aionapi.dev','user'),
    ('Carla Fernandes','carla','$2a$10$AQjsEvApPM59JVg4XIqEluc.JtKfnyIWPLK7QEYW6CPDN9cU/OH7O','carla@aionapi.dev','user'),
    ('Daniel Dias','daniel','$2a$10$AQjsEvApPM59JVg4XIqEluc.JtKfnyIWPLK7QEYW6CPDN9cU/OH7O','daniel@aionapi.dev','user'),
    ('Ester Rodrigues','ester','$2a$10$AQjsEvApPM59JVg4XIqEluc.JtKfnyIWPLK7QEYW6CPDN9cU/OH7O','ester@aionapi.dev','user')
ON CONFLICT (username) DO NOTHING;
