CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(64) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio VARCHAR(100) DEFAULT '',
    gen VARCHAR(10) DEFAULT 'male' NOT NULL,
    avatar_url VARCHAR(255),
    is_superuser BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tweets (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    content VARCHAR(280) NOT NULL,
    parent_tweet_id INT DEFAULT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tweet_media (
    tweet_id INT PRIMARY KEY REFERENCES tweets(id) ON DELETE CASCADE,
    media_url VARCHAR(255),
    media_type VARCHAR(15) NOT NULL
);

CREATE TABLE follows (
    follower_id INT REFERENCES users(id) ON DELETE CASCADE,
    following_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE TABLE likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE retweets (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE secret_questions (
    user_id INT PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question VARCHAR(100) NOT NULL,
    answer VARCHAR(100) NOT NULL
);