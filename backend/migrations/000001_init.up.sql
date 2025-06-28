CREATE TABLE users (
    id SERIAL PRIMARY KEY, 
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio VARCHAR(100) DEFAULT '',
    avatar_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tweets (
    id SERIAL PRIMARY KEY, 
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
    content VARCHAR(280) NOT NULL,
    parent_tweet_id INT REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    like_count INT DEFAULT 0,
    reply_count INT DEFAULT 0,
    retweet_count INT DEFAULT 0
);

CREATE TABLE tweet_media (
    tweet_id INT REFERENCES tweets(id) ON DELETE CASCADE,
    media_url VARCHAR(255),
    media_type VARCHAR(10) NOT NULL CHECK (media_type IN ('image', 'video', 'gif'))
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
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id INT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, tweet_id)
);