-- Insert authors
INSERT INTO authors (name, bio)
VALUES ('J.K. Rowling', 'British author best known for the Harry Potter series'),
       ('George R.R. Martin', 'American novelist and short story writer, author of A Song of Ice and Fire'),
       ('Jane Austen', 'English novelist known for her romance fiction'),
       ('Stephen King', 'American author of horror, supernatural fiction, and fantasy'),
       ('Agatha Christie', 'British crime novelist, short story writer, and playwright');

-- Insert books
INSERT INTO books (title, isbn, description, price, author_id, published_date)
VALUES ('Harry Potter and the Philosopher''s Stone', '9780747532699', 'The first book in the Harry Potter series',
        19.99, 1, '1997-06-26'),
       ('Harry Potter and the Chamber of Secrets', '9780747538486', 'The second book in the Harry Potter series', 19.99,
        1, '1998-07-02'),
       ('A Game of Thrones', '9780553103540', 'The first book in A Song of Ice and Fire series', 24.99, 2,
        '1996-08-01'),
       ('A Clash of Kings', '9780553108033', 'The second book in A Song of Ice and Fire series', 24.99, 2,
        '1998-11-16'),
       ('Pride and Prejudice', '9780141439518', 'A romantic novel of manners', 15.99, 3, '1813-01-28'),
       ('Emma', '9780141439587', 'A novel about youthful hubris and romantic misunderstandings', 15.99, 3,
        '1815-12-23'),
       ('The Shining', '9780307743657', 'A horror novel set in an isolated hotel', 21.99, 4, '1977-01-28'),
       ('The Stand', '9780307743681', 'A post-apocalyptic horror/fantasy novel', 22.99, 4, '1978-10-03'),
       ('Murder on the Orient Express', '9780007119318', 'A detective novel featuring Hercule Poirot', 18.99, 5,
        '1934-01-01'),
       ('Death on the Nile', '9780007119325', 'A mystery novel featuring Hercule Poirot', 18.99, 5, '1937-11-01');

-- Book sales (10k total sales randomised across all books)
WITH
    date_range AS (
        SELECT
            '2024-01-01'::timestamp AS start_date,
            NOW() AS end_date
    )
INSERT INTO book_sales (book_id, purchased_at)
SELECT
    books.id,
    date_range.start_date + (random() * (date_range.end_date - date_range.start_date))
FROM
    books,
    date_range,
    generate_series(1, 10000);
