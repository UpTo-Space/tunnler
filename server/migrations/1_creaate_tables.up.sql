CREATE TABLE connections (
    id SERIAL PRIMARY KEY,
    client_ip INET,
    local_port INTEGER,
    remote_port INTEGER,
    subdomain TEXT
);