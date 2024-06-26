BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS USERS (
    id UUID PRIMARY KEY,
    email VARCHAR UNIQUE,
    name VARCHAR,
    password VARCHAR,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS cats (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    race VARCHAR NOT NULL,
    sex VARCHAR NOT NULL,
    age_in_month INTEGER NOT NULL,
    description VARCHAR NOT NULL,
    image_urls VARCHAR[] NOT NULL,
    has_matched BOOLEAN DEFAULT FALSE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);
-- {
-- 	"matchCatId": "",
-- 	"userCatId": "",
-- 	"message": "" // not null, minLength: 5, maxLength: 120
-- }
CREATE TABLE IF NOT EXISTS match_cats (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    match_cat_id INT REFERENCES cats(id) ON DELETE CASCADE,
    user_cat_id INT REFERENCES cats(id) ON DELETE CASCADE,
    message VARCHAR NOT NULL,
    has_approved BOOLEAN DEFAULT NULL CHECK (has_approved IS NULL OR has_approved IN (TRUE, FALSE)),
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

COMMIT TRANSACTION;