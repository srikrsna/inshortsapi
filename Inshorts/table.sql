CREATE TABLE articles (
        id serial NOT NULL PRIMARY KEY,
        title text NOT NULL,
        subtitle text NOT NULL,
        content text date,
        creationtimestamp timestamp
);