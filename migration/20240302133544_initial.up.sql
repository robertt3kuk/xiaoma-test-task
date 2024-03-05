-- Customer migration
CREATE TABLE customer (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL UNIQUE, -- Added UNIQUE constraint
    balance FLOAT8 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE 
);

-- Item migration
CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL UNIQUE, -- Added UNIQUE constraint
    cost FLOAT8 NOT NULL,
    price FLOAT8 NOT NULL,
    sort INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE 
);

-- Transaction migration
CREATE TABLE transaction (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    item_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    price FLOAT8 NOT NULL,
    amount FLOAT8 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (customer_id) REFERENCES customer(id),
    FOREIGN KEY (item_id) REFERENCES item(id)
);
