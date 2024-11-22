CREATE TABLE authors (
                         id BIGSERIAL PRIMARY KEY,
                         name TEXT NOT NULL,
                         bio TEXT,
                         created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);