CREATE TABLE books (
                       id BIGSERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       isbn TEXT NOT NULL UNIQUE,
                       description TEXT,
                       price DECIMAL(10,2) NOT NULL,
                       author_id BIGINT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
                       published_date DATE NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);