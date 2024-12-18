DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS media CASCADE;
DROP TABLE IF EXISTS media_transaction;
DROP TABLE IF EXISTS items CASCADE;
DROP TABLE IF EXISTS purchase_list_items;
DROP TABLE IF EXISTS account_access;
DROP TABLE IF EXISTS financial_groups CASCADE;
DROP TABLE IF EXISTS user_financial_groups;

DROP TYPE IF EXISTS UserStatus;
DROP TYPE IF EXISTS UserRole;
DROP TYPE IF EXISTS TransactionType;
DROP TYPE IF EXISTS CategoryEntityType;
DROP TYPE IF EXISTS AccessLevel CASCADE;
DROP TYPE IF EXISTS MediaStatus CASCADE;
DROP TYPE IF EXISTS MediaAccess;
DROP TYPE IF EXISTS AccountType;

CREATE TYPE UserStatus AS ENUM ('banned', 'verified', 'disabled', 'locked', 'pending');

CREATE TYPE UserRole AS ENUM ('admin', 'user');

CREATE TYPE TransactionType AS ENUM ('income', 'expense', 'transfer');

CREATE TYPE CategoryEntityType AS ENUM ('income', 'expense', 'account');

CREATE TYPE AccessLevel AS ENUM ('view', 'edit', 'all');

CREATE TYPE MediaStatus AS ENUM ('temp', 'attached', 'removed');

CREATE TYPE MediaAccess AS ENUM ('owner', 'group', 'public');

CREATE TYPE AccountType AS ENUM ('self', 'external');

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
    middle_name VARCHAR(255),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    file_path TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    metadata TEXT,
    status MediaStatus NOT NULL DEFAULT 'temp',
    financial_group_id INT NOT NULL,
    access MediaAccess NOT NULL DEFAULT 'owner',
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
    type AccountType DEFAULT 'self',
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

CREATE TABLE media_transaction (
    id SERIAL PRIMARY KEY,
    media_id INT NOT NULL REFERENCES media(id),
    transaction_id INT NOT NULL REFERENCES transactions(id)
);

CREATE TABLE financial_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    image_id INT REFERENCES media(id),
    user_id UUID NOT NULL REFERENCES users(id),
    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_financial_groups (
    id SERIAL PRIMARY KEY,
    financial_group_id INT NOT NULL REFERENCES financial_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
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
ALTER TABLE media DROP CONSTRAINT IF EXISTS fk_financial_group_id;
ALTER TABLE financial_groups DROP CONSTRAINT IF EXISTS fk_fg_image_id;

ALTER TABLE media ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users ADD CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles(id);
ALTER TABLE profiles ADD CONSTRAINT fk_picture_id FOREIGN KEY (picture_id) REFERENCES media(id);
ALTER TABLE media ADD CONSTRAINT fk_financial_group_id FOREIGN KEY (financial_group_id) REFERENCES financial_groups(id);
ALTER TABLE financial_groups ADD CONSTRAINT fk_fg_image_id FOREIGN KEY (image_id) REFERENCES media(id);

CREATE OR REPLACE FUNCTION check_item_category()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM categories
        WHERE id = NEW.category_id
            AND entity_type = 'expense'
    ) THEN
        RAISE EXCEPTION 'Invalid category_id: must reference a category with entity_type Income'
            USING ERRCODE = 'S0001';
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
            AND entity_type = 'account'
    ) THEN
        RAISE EXCEPTION 'Invalid category_id: must reference a category with entity_type Account'
            USING ERRCODE = 'S0001';
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
    IF NEW IS DISTINCT FROM OLD THEN
        NEW.update_date := CURRENT_TIMESTAMP;
    END IF;
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
CREATE TRIGGER update_date_trigger BEFORE UPDATE ON financial_groups
    FOR EACH ROW EXECUTE PROCEDURE update_date_on_change();

CREATE OR REPLACE FUNCTION update_profile_picture_check()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS(
        SELECT 1
        FROM media
        WHERE id = NEW.picture_id
    ) THEN
         RAISE EXCEPTION 'Foreign key violation: % does not exists in media', NEW.picture_id
            USING ERRCODE = 'S0002';
    END IF;
       
    IF NOT EXISTS(
        SELECT 1
        FROM media
        JOIN users on users.profile_id = NEW.id
        WHERE id = NEW.picture_id
        AND user_id = users.id
    ) THEN
        RAISE EXCEPTION 'Unauthorized access: % does not belong to user %', NEW.picture_id, NEW.user_id
            USING ERRCODE = 'S0004';
    END IF;

    UPDATE media
    SET status = 'attached'
    WHERE id = NEW.picture_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_profile_picture_on_change BEFORE UPDATE OF picture_id ON profiles
    FOR EACH ROW EXECUTE PROCEDURE update_profile_picture_check();

CREATE OR REPLACE FUNCTION update_item_image_check()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS(
        SELECT 1
        FROM media
        WHERE id = NEW.image_id
        AND user_id = NEW.user_id
    ) THEN
        UPDATE media
        SET status = 'attached'
        WHERE id = NEW.image_id;
    ELSE
        RAISE EXCEPTION 'Foreign key violation: % does not exists in media', NEW.image_id
            USING ERRCODE = 'S0002';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_item_image_on_change BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE PROCEDURE update_item_image_check();
CREATE TRIGGER update_item_image_on_change BEFORE UPDATE ON financial_groups
    FOR EACH ROW EXECUTE PROCEDURE update_item_image_check();

CREATE OR REPLACE FUNCTION update_category_icon_check()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS(
        SELECT 1
        FROM media
        WHERE id = NEW.icon_id
        AND user_id = NEW.user_id
    ) THEN
        UPDATE media
        SET status = 'attached'
        WHERE id = NEW.icon_id;
    ELSE 
        RAISE EXCEPTION 'Foreign key violation: % does not exists in media', NEW.icon_id
            USING ERRCODE = 'S0002';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_category_icon_on_change BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE PROCEDURE update_category_icon_check();

CREATE OR REPLACE FUNCTION add_user_to_financial_group_check()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS(
        SELECT 1
        FROM user_financial_groups
        WHERE user_id = NEW.user_id
        AND financial_group_id = NEW.financial_group_id
    ) THEN
        RAISE EXCEPTION 'User is already in selected group'
            USING ERRCODE = 'S0003';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_creating_existing_group_member BEFORE INSERT ON user_financial_groups
    FOR EACH ROW EXECUTE PROCEDURE add_user_to_financial_group_check();

CREATE OR REPLACE FUNCTION add_financial_group_insert_owner()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO user_financial_groups (financial_group_id, user_id)
    VALUES (NEW.id, NEW.user_id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER add_owner_to_financial_group AFTER INSERT ON financial_groups
    FOR EACH ROW EXECUTE PROCEDURE add_financial_group_insert_owner();
