DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS media CASCADE;
DROP TABLE IF EXISTS media_bindings;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS purchase_list_items;
DROP TABLE IF EXISTS account_access;
DROP TABLE IF EXISTS financial_groups CASCADE;
DROP TABLE IF EXISTS user_financial_groups;

DROP TYPE IF EXISTS UserStatus;
DROP TYPE IF EXISTS UserRole;
DROP TYPE IF EXISTS TransactionType;
DROP TYPE IF EXISTS CategoryEntityType;
DROP TYPE IF EXISTS MediaBind CASCADE;
DROP TYPE IF EXISTS AccessLevel CASCADE;

CREATE TYPE UserStatus AS ENUM ('banned', 'verified', 'disabled', 'locked', 'pending');

CREATE TYPE UserRole AS ENUM ('admin', 'user');

CREATE TYPE TransactionType AS ENUM ('income', 'expense', 'transfer');

CREATE TYPE CategoryEntityType AS ENUM ('income', 'expense', 'account');

CREATE TYPE MediaBind AS ENUM ('item', 'transaction');

CREATE TYPE AccessLevel AS ENUM ('view', 'edit', 'all');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    ip VARCHAR(45),
    password VARCHAR(255) NOT NULL,
    last_login TIMESTAMP,
    last_password_change TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    failed_tries INT DEFAULT 0,
    status UserStatus DEFAULT 'pending',
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    profile_id INT NOT NULL
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    picture_id INT,
    phone_number VARCHAR(15),
    address TEXT,
    name VARCHAR(255),
    family_name VARCHAR(255),
    middle_name VARCHAR(255)
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    file_path TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    metadata JSON,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    color VARCHAR(7) NOT NULL CHECK (color ~ '^#[0-9a-fA-F]{6}$'),
    icon_id INT REFERENCES media(id),
entity_type CategoryEntityType NOT NULL,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id),
    balance REAL DEFAULT 0.0,
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    account_id INT NOT NULL REFERENCES accounts(id),
    category_id INT NOT NULL REFERENCES categories(id),
    amout REAL NOT NULL,
    description TEXT,
    transaction_type TransactionType NOT NULL,
    linked_transaction_id INT REFERENCES transactions(id),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    image_id INT REFERENCES media(id),
    category_id INT NOT NULL REFERENCES categories(id),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE purchase_list_items (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    item_id INT NOT NULL REFERENCES items(id),
    count INT NOT NULL DEFAULT 1,
    unit_price REAL NOT NULL,
    transaction_id INT NOT NULL REFERENCES transactions(id),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE media_bindings (
    id SERIAL PRIMARY KEY,
    media_id INT NOT NULL REFERENCES media(id),
    bind_type MediaBind NOT NULL,
    bind_id INT NOT NULL
);

CREATE TABLE financial_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_financial_groups (
    id SERIAL PRIMARY KEY,
    financial_group_id INT NOT NULL REFERENCES financial_groups(id),
    user_id UUID NOT NULL REFERENCES users(id)
);

CREATE TABLE account_access (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    account_id INT NOT NULL REFERENCES accounts(id),
    access AccessLevel NOT NULL DEFAULT 'view',
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE media DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_profile_id;
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS fk_picture_id;

ALTER TABLE media ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users ADD CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles(id);
ALTER TABLE profiles ADD CONSTRAINT fk_picture_id FOREIGN KEY (picture_id) REFERENCES media(id);

CREATE OR REPLACE FUNCTION check_item_category()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM categories
        WHERE id = NEW.category_id
            AND entity_type IN ('Income', 'Expense')
    ) THEN
        RAISE EXCEPTION 'Invalid category_id: must reference a category with entity_type Income or Expense';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_account_category()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM categories
        WHERE id = NEW.category_id
            AND entity_type = 'Account'
    ) THEN
        RAISE EXCEPTION 'Invalid category_id: must reference a category with entity_type Account';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_item_category
BEFORE INSERT OR UPDATE ON items
FOR EACH ROW
EXECUTE FUNCTION check_item_category();

CREATE TRIGGER validate_item_category
BEFORE INSERT OR UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION check_account_category();

CREATE OR REPLACE FUNCTION update_date_on_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT (NEW IS DISTINCT FROM OLD) THEN
        RETURN NULL;
    END IF;

    NEW.update_date := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_date_trigger BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON profiles
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON media
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON purchase_list_items
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON account_access
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();
