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

	rows, err := dataaccess.DB.Query("SELECT * FROM itemsdb.items")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.ImageName, &item.Price, &item.SellerID); err != nil {
			if err == sql.ErrNoRows {
				fmt.Print(err)
				return nil, sql.ErrNoRows
			}
		}

		dbItems = append(dbItems, item)
	}

	return dbItems, nil
}

func (dataaccess *DataAccess) DeleteItem(itemID int64) error {
	_, err := dataaccess.DB.Exec("DELETE FROM itemsdb.items WHERE itemID = ?", itemID)
	return err
}

func (dataaccess *DataAccess) CreateItem(itemID string, itemName string, itemPrice string) error {

	_, err := dataaccess.DB.Exec("INSERT INTO itemsdb.items"+
		"(itemID, itemName, itemImageName, itemPrice, itemSellerID)"+
		"VALUES (?, ?, ?, ?, ?)", itemID, itemName, itemName+".svg", itemPrice, 1)

	if err != nil {
		return err
	}

	return nil
}

func (dataaccess *DataAccess) UpdatePrice(itemID int64, newPrice string) error {

	_, err := dataaccess.DB.Exec("UPDATE itemsdb.items SET itemPrice = (?) where itemID = (?)", newPrice, itemID)

	if err != nil {
		return err
	}

	return nil
}
