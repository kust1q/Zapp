CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(64) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio VARCHAR(100) DEFAULT '',
    gen VARCHAR(10) DEFAULT 'male' NOT NULL,
    is_superuser BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tweets (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    content VARCHAR(280) NOT NULL,
    parent_tweet_id INT DEFAULT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tweet_media (
    id SERIAL PRIMARY KEY,
    tweet_id INT NOT NULL UNIQUE REFERENCES tweets(id) ON DELETE CASCADE,
    media_url TEXT NOT NULL,
    mime_type VARCHAR(15) NOT NULL,
    size_bytes BIGINT
);  

CREATE TABLE IF NOT EXISTS avatars (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    media_url TEXT NOT NULL,
    mime_type VARCHAR(15) NOT NULL,
    size_bytes BIGINT
);

CREATE TABLE IF NOT EXISTS follows (
    follower_id INT REFERENCES users(id) ON DELETE CASCADE,
    following_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE TABLE IF NOT EXISTS likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE IF NOT EXISTS retweets (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS secret_questions (
    user_id INT PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question VARCHAR(100) NOT NULL,
    answer VARCHAR(100) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_tweets_user_created_at ON tweets(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tweets_parent ON tweets(parent_tweet_id);

CREATE INDEX IF NOT EXISTS idx_follows_following ON follows(following_id);

CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id);

CREATE INDEX IF NOT EXISTS idx_likes_tweet ON likes(tweet_id);

CREATE INDEX IF NOT EXISTS idx_retweets_tweet ON retweets(tweet_id);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_media_tweet ON tweet_media(tweet_id);