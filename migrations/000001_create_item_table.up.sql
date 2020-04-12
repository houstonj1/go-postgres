CREATE TABLE IF NOT EXISTS public.item(
    id text NOT NULL PRIMARY KEY,
    name text NOT NULL UNIQUE,
    description text
);