-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: GetAuthorStats :one
SELECT 
    a.id, 
    a.name, 
    COUNT(b.id) as total_books, 
    ROUND(AVG(b.price::numeric), 2) as average_price,
    TO_CHAR(MIN(b.published_date), 'YYYY-MM-DD') as earliest_publication,
    TO_CHAR(MAX(b.published_date), 'YYYY-MM-DD') as latest_publication,
    ROUND(SUM(bs.total_revenue), 2) as total_revenue,
    jsonb_object_agg(year::text, book_count)::text AS books_published_per_year
FROM 
    authors a 
LEFT JOIN 
    books b ON a.id = b.author_id 
LEFT JOIN 
    (
        SELECT 
            book_id, 
            COUNT(bs.id) AS total_sales, 
            SUM(b.price::numeric) * COUNT(bs.id) AS total_revenue
        FROM 
            book_sales bs
        JOIN 
            books b ON bs.book_id = b.id
        GROUP BY 
            book_id
    ) bs ON b.id = bs.book_id
LEFT JOIN 
    (
        SELECT 
            EXTRACT(YEAR FROM b.published_date) AS year,
            COUNT(b.id) AS book_count
        FROM 
            books b
        GROUP BY 
            year
    ) yearly_books ON EXTRACT(YEAR FROM b.published_date) = yearly_books.year
WHERE 
    a.id = $1 
GROUP BY 
    a.id, a.name;

