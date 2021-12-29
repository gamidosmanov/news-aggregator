BEGIN;

    DROP TABLE IF EXISTS rss_aggr.posts;

    DROP SCHEMA IF EXISTS rss_aggr;

    CREATE SCHEMA IF NOT EXISTS rss_aggr
        AUTHORIZATION postgres;

    GRANT ALL ON SCHEMA rss_aggr TO postgres;

    DROP TABLE IF EXISTS
        rss_aggr.posts;

    CREATE TABLE IF NOT EXISTS rss_aggr.posts
    (
        id SERIAL,
        title TEXT,
        content TEXT,
        pub_time BIGINT,
        link TEXT NOT NULL,
        CONSTRAINT post_id PRIMARY KEY (id),
        CONSTRAINT post_link UNIQUE(link)
    );   

END;