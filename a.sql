CREATE TABLE usrs (
    id serial PRIMARY KEY,
    nickname varchar(30) UNIQUE NOT NULL,
    password varchar(30) NOT NULL
);

CREATE TABLE items (
    id serial PRIMARY KEY,
    title text UNIQUE NOT NULL
)

CREATE TABLE categories (
    id serial PRIMARY KEY,
    title text UNIQUE NOT NULL
)

CREATE TABLE item_category (
    item_id int REFERENCES items (id) ON DELETE CASCADE,
    category_id int REFERENCES category (id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, category_id)
)



CREATE FUNCTION update_item_categories(id int, category_ids int[]) RETURNS void LANGUAGE plpgsql AS
$$
BEGIN
	DELETE FROM item_category WHERE item_id = id;
	INSERT INTO item_category (item_id, category_id) VALUES (id, UNNEST(category_ids));
END;
$$;



CREATE FUNCTION update_content(t text, c1 text, c2 text) RETURNS void LANGUAGE plpgsql AS
$$
DECLARE 
    t_id int;
    c1_id int;
    c2_id int;
BEGIN
    select id into t_id from items where title = t;
    if not found then begin
        insert into items (title) values (t) returning id into t_id;
        select id into c1_id from categories where title = c1;
        if not found then insert into categories (title) values (c1) returning id into c1_id; end if;
        insert into item_category (item_id, category_id) values (t_id, c1_id);
        if c1 != c2 then begin 
            select id into c2_id from categories where title = c2;
            if not found then insert into categories (title) values (c2) returning id into c2_id; end if;
            insert into item_category (item_id, category_id) values (t_id, c2_id);
        end; end if;
    end; end if;
END;
$$;
