package shopapi

import (
	"database/sql"
	"fmt"
)

type DataAccess struct {
	DB *sql.DB
}

func (dataaccess *DataAccess) GetAllItems() ([]Item, error) {

	dbItems := make([]Item, 0, 10)
	var isSoldByte []byte

	rows, err := dataaccess.DB.Query("SELECT * FROM itemsdb.items")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) DeleteItem(itemID string) error {
	_, err := dataaccess.DB.Exec("DELETE FROM itemsdb.items WHERE ItemID = ?", itemID)
	return err
}

func (dataaccess *DataAccess) CreateItem(item Item) error {

	_, err := dataaccess.DB.Exec("INSERT INTO itemsdb.items"+
		"(ItemID, ItemName, ItemGender, ItemDescription, ItemImageName, ItemGalleryImage, ItemUploadDate, ItemPrice, ItemIsSold, ItemSize, ItemCategory, ItemCondition, ItemSellerID, ItemSellerName)"+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", item.ID, item.Name, item.Gender, item.Description, item.ImageName, item.GalleryImage, item.DateUploaded, item.Price, item.IsSold, item.Size, item.Category, item.Condition, item.SellerID, item.SellerName)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdateItem(item Item) error {

	_, err := dataaccess.DB.Exec("UPDATE itemsdb.items SET "+
		"ItemName = (?), ItemGender= (?), ItemDescription= (?), ItemPrice= (?), ItemIsSold= (?), ItemSize= (?), ItemCategory= (?), ItemCondition = (?)"+
		"where ItemID = (?)", item.Name, item.Gender, item.Description, item.Price, item.IsSold, item.Size, item.Category, item.Condition, item.ID)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) GetItem(itemID string) (Item, error) {
	row := dataaccess.DB.QueryRow("SELECT * FROM itemsdb.items WHERE itemID = ?", itemID)

	var item Item
	var isSoldByte []byte

	if err := row.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
		if err == sql.ErrNoRows {
			fmt.Print(err)
			return Item{}, sql.ErrNoRows
		}
	}

	item.IsSold = isSoldByte[0] == 1

	return item, nil
}

func (dataaccess *DataAccess) GetItemsByQueryTerm(query string) ([]Item, error) {
	rows, err := dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemName LIKE ?", query)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	dbItems := make([]Item, 0, 100)

	for rows.Next() {
		var item Item
		var isSoldByte []byte

		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) getSortedItemsByPriceInc(query, category string) ([]Item, error) {
	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemCategory LIKE ? ORDER BY ItemPrice", category)
	} else if query != "" {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemName LIKE ? ORDER BY ItemPrice", query)
	} else {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items ORDER BY ItemPrice")
	}

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	dbItems := make([]Item, 0, 100)

	for rows.Next() {
		var item Item
		var isSoldByte []byte

		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) getSortedItemsByPriceDec(query, category string) ([]Item, error) {
	var rows *sql.Rows
	var err error
	if category != "" {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemCategory LIKE ? ORDER BY ItemPrice DESC", category)
	} else if query != "" {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemName LIKE ? ORDER BY ItemPrice DESC", query)
	} else {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items ORDER BY ItemPrice DESC")
	}

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	dbItems := make([]Item, 0, 100)

	for rows.Next() {
		var item Item
		var isSoldByte []byte

		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) getItemsByCategory(category string) ([]Item, error) {

	rows, err := dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemCategory LIKE ?", category)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	dbItems := make([]Item, 0, 100)

	for rows.Next() {
		var item Item
		var isSoldByte []byte

		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) GetItemsBySeller(userID, limit int) ([]Item, error) {
	var rows *sql.Rows
	var err error
	if limit < 0 {
		// Get all items by seller
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemSellerID = ?", userID)
	} else {
		rows, err = dataaccess.DB.Query("SELECT * FROM itemsdb.items WHERE ItemSellerID = ? LIMIT ?", userID, limit)

	}

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	dbItems := make([]Item, 0, 100)

	for rows.Next() {
		var item Item
		var isSoldByte []byte

		if err := rows.Scan(&item.ID, &item.Name, &item.Gender, &item.Description, &item.ImageName, &item.GalleryImage, &item.DateUploaded, &item.Price, &isSoldByte, &item.Size, &item.Category, &item.Condition, &item.SellerID, &item.SellerName); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		item.IsSold = isSoldByte[0] == 1

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) ListAsSold(itemID string) error {

	_, err := dataaccess.DB.Exec("UPDATE itemsdb.items SET ItemIsSold = 1 where itemID = (?)", itemID)

	if err != nil {
		return err
	}

	return nil
}
