package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/smtnn-ks/test-task-inhousead/db"
)

func checkAuth(name, password string) error {
	row := db.Client.QueryRow("SELECT password FROM usrs WHERE nickname = $1", name) // nickname is unique
	var p string
	err := row.Scan(&p)
	if err != nil && p != password {
		return fiber.ErrUnauthorized
	}
	return nil
}

func register(name, password string) error {
	_, err := db.Client.Exec("INSERT INTO usrs (nickname, password) VALUES ($1, $2)", name, password)
	sqlError, ok := err.(*pq.Error)
	if ok && sqlError.Code == "23505" {
		return fiber.ErrBadRequest
	} else {
		log.Println("[ERROR] ", err)
		return fiber.ErrInternalServerError
	}
}

func getCategories() ([]entity_t, error) {
	rows, err := db.Client.Query("SELECT * FROM categories")
	if err != nil {
		log.Println("[ERROR] ", err)
		return nil, fiber.ErrInternalServerError
	}
	var categories []entity_t
	for rows.Next() {
		category := entity_t{}
		err := rows.Scan(&category.Id, &category.Title)
		if err != nil {
			log.Println("[ERROR] ", err)
			return nil, fiber.ErrInternalServerError
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func getItems(categoryId int) ([]entity_t, error) {
	rows, err := db.Client.Query(
		`SELECT items.id, items.title 
        FROM item_category 
        JOIN items 
        ON item_category.item_id = items.id 
        WHERE item_category.category_id = $1`,
		categoryId,
	)
	if err != nil {
		log.Println("[ERROR] ", err)
		return nil, fiber.ErrInternalServerError
	}

	var items []entity_t
	for isNext := rows.Next(); isNext; isNext = rows.Next() {
		item := entity_t{}
		err := rows.Scan(&item.Id, &item.Title)
		if err != nil {
			log.Println("[ERROR] ", err)
			return nil, fiber.ErrInternalServerError
		}
		items = append(items, item)
	}
	return items, nil
}

func createCategory(title string) error {
	_, err := db.Client.Exec("INSERT INTO categories (title) VALUES ($1)", title)
	if err != nil {
		sqlError, ok := err.(*pq.Error)
		if ok && sqlError.Code == "23505" {
			return fiber.ErrBadRequest
		}
		return err
	}
	return nil
}

func updateCategory(id int, title string) error {
	_, err := db.Client.Exec("UPDATE categories SET title = $1 WHERE id = $2", title, id)
	if err != nil {
		return err
	}
	return nil
}

func removeCategory(id int) error {
	_, err := db.Client.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func createItem(title string, categoryIds []int) error {
	_, err := db.Client.Exec(
		`WITH new_item AS (INSERT INTO items (title) VALUES ($1) RETURNING id) 
        INSERT INTO item_category (item_id, category_id)
        SELECT new_item.id, unnest($2::INT[]) FROM new_item`,
		title,
		pq.Array(categoryIds),
	)
	if err != nil {
		sqlError, ok := err.(*pq.Error)
		if ok && sqlError.Code == "23505" {
			return fiber.ErrBadRequest
		}
		return err
	}
	return nil
}

func updateItem(id int, title string, categoryIds []int) error {
	if title != "" {
		_, err := db.Client.Exec("UPDATE items SET title = $1 WHERE id = $2", title, id)
		if err != nil {
			return err
		}
	}
	if categoryIds != nil && len(categoryIds) != 0 {
		_, err := db.Client.Exec(
			"SELECT update_item_categories($1, $2::int[])",
			id,
			pq.Array(categoryIds),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeItem(id int) error {
	_, err := db.Client.Exec(
		"DELETE FROM items WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
