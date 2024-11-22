CREATE TABLE book_sales (
                            id BIGSERIAL PRIMARY KEY,
                            book_id BIGINT NOT NULL REFERENCES books(id),
                            purchased_at TIMESTAMP WITH TIME ZONE NOT NULL
);
