package dbconn

import (
	"errors"
	"fmt"
)

// grant flag
const (
	GRANT_USER = 1 << iota
	GRANT_PRODUCT_ADD
	GRANT_PRODUCT_EDIT
	GRANT_ITEM_READ
	GRANT_ITEM_ADD
	GRANT_ITEM_EDIT

	// GRANT_SUPER = (1 << iota) - 1
)

type ItemInfomation struct {
	Id     string
	Name   string
	Brand  string
	Color  string
	Size   string
	Price  float64
	Status int
}
type Table struct {
	Header  []string
	Content [][]string
}

func (tb *Table) Init(header ...string) {
	tb.Header = header[:]
	tb.Content = nil
}

var saltStrRander, itemIDStrRander RandStringMaker

func init() {
	saltStrRander.Set("abcdefghijklmnopqrstuvwxyz", 7)
	itemIDStrRander.Set("0123456789ABCDEFGHJKLMNPQRSTUVWXYZ", 20)
}

func GetItemInfomation(id string) (info ItemInfomation, err error) {
	tb_name := "item INNER JOIN product ON it_pd_id = pd_id INNER JOIN pattern ON pd_pt_id = pt_id"
	tb, err := Select(tb_name, []string{"it_id ="}, []string{id}, []string{"it_id", "pt_name", "pt_brand", "pd_color", "pd_size", "pt_price", "it_status"})
	if err != nil {
		return
	}
	if len(tb.Content) == 0 {
		err = errors.New("货品不存在")
		return
	}
	info.Id = tb.Content[0][0]
	info.Name = tb.Content[0][1]
	info.Brand = tb.Content[0][2]
	info.Color = tb.Content[0][3]
	info.Size = tb.Content[0][4]
	fmt.Sscan(tb.Content[0][5], &info.Price)
	fmt.Sscan(tb.Content[0][6], &info.Status)
	return
}

// add remove set list
//user
func AddUser(name, pw string, grant int) (res int64, err error) {
	salt := saltStrRander.Get()
	passStr := getPassHex(salt, pw)
	return Insert("user", []string{"u_name", "u_salt", "u_pw", "u_grant"}, []string{name, salt, passStr, fmt.Sprintf("%d", grant)})
}
func SetUserPassword(id, pw string) (res int64, err error) {
	salt := saltStrRander.Get()
	passStr := getPassHex(salt, pw)
	return Update("user", []string{"u_salt", "u_pw"}, []string{salt, passStr}, []string{"u_id ="}, []string{id})
}

// func SetUserGrant(id string, grant int) (res int64, err error) {
// 	return Update("user", []string{"u_grant"}, []string{fmt.Sprintf("%d", grant)}, []string{"u_id ="}, []string{id})
// }
// func RemoveUser(id string) (res int64, err error) {
// 	return Delete("user", []string{"u_id ="}, []string{id})
// }
// func ListUser() (tb Table, err error) {
// 	return Select("user", nil, nil, []string{"u_id", "u_name", "u_grant"})
// }
func GetUserID(name string) (id string, err error) {
	id = ""
	tb, err := Select("user", []string{"u_name ="}, []string{name}, []string{"u_id"})
	if err != nil {
		return
	}
	if len(tb.Content) != 1 || len(tb.Content[0]) != 1 {
		err = errors.New("用户不存在")
		return
	}
	id = tb.Content[0][0]
	err = nil
	return
}

//Batch
// func ListBatch() (tb Table, err error) {
// 	return Select("batch", nil, nil, []string{"bt_id", "bt_u_id", "bt_time"})
// }
// func AddBatch(uid string) (res int64, err error) {
// 	return Insert("batch", []string{"bt_u_id"}, []string{fmt.Sprintf("%v", uid)})
// }
// func SetBatch(id string, change map[string]string) (res int64, err error) {
// 	var key, value []string
// 	keys := []string{"bt_id", "bt_u_id", "bt_time"}
// 	for _, k := range keys {
// 		if v, ok := change[k]; ok {
// 			key = append(key, k)
// 			value = append(value, v)
// 		}
// 	}
// 	return Update("batch", key, value, []string{"bt_id ="}, []string{id})
// }
// func RemoveBatch(id string) (res int64, err error) {
// 	Delete("item", []string{"it_bt_id ="}, []string{id})
// 	return Delete("batch", []string{"bt_id ="}, []string{id})
// }

//Pattern
// func ListPattern() (tb Table, err error) {
// 	return Select("pattern", nil, nil, []string{"pt_id", "pt_name", "pt_brand", "pt_price"})
// }
// func AddPattern(name, brand string, price float64) (res int64, err error) {
// 	return Insert("pattern", []string{"pt_name", "pt_brand", "pt_price"}, []string{name, brand, fmt.Sprintf("%v", price)})
// }
// func SetPattern(id string, change map[string]string) (res int64, err error) {
// 	var key, value []string
// 	keys := []string{"pt_id", "pt_name", "pt_brand", "pt_price"}
// 	for _, k := range keys {
// 		if v, ok := change[k]; ok {
// 			key = append(key, k)
// 			value = append(value, v)
// 		}
// 	}
// 	return Update("pattern", key, value, []string{"pt_id ="}, []string{id})
// }
// func RemovePattern(id string) (res int64, err error) {
// 	return Delete("pattern", []string{"pt_id ="}, []string{id})
// }

//Product
// func ListProduct() (tb Table, err error) {
// 	return Select("product", nil, nil, []string{"pd_id", "pd_pt_id", "pd_SKU", "pd_color", "pd_size"})
// }
// func AddProduct(pt_id int, SKU, color, size string) (res int64, err error) {
// 	return Insert("product", []string{"pd_pt_id", "pd_SKU", "pd_color", "pd_size"}, []string{fmt.Sprintf("%v", pt_id), SKU, color, size})
// }
// func SetProduct(id string, change map[string]string) (res int64, err error) {
// 	var key, value []string
// 	keys := []string{"pd_id", "pd_pt_id", "pd_SKU", "pd_color", "pd_size"}
// 	for _, k := range keys {
// 		if v, ok := change[k]; ok {
// 			key = append(key, k)
// 			value = append(value, v)
// 		}
// 	}
// 	return Update("product", key, value, []string{"pd_id ="}, []string{id})
// }
// func RemoveProduct(id string) (res int64, err error) {
// 	Delete("item", []string{"it_pd_id ="}, []string{id})
// 	return Delete("product", []string{"pd_id ="}, []string{id})
// }

//Item
// func ListItem() (tb Table, err error) {
// 	return Select("item", nil, nil, []string{"it_id", "it_pd_id", "it_bt_id", "it_status"})
// }
func AddItem(num, pd_id, bt_id, status int) (res int64, err error) {
	var val [4]string
	val[1] = fmt.Sprintf("%d", pd_id)
	val[2] = fmt.Sprintf("%d", bt_id)
	val[3] = fmt.Sprintf("%d", status)
	for i := int64(0); i < int64(num); {
		for {
			val[0] = itemIDStrRander.Get()
			tb, err := Select("item", []string{"it_id ="}, val[0:1], []string{"it_id"})
			if err != nil {
				return i, err
			}
			if len(tb.Content) == 0 {
				break
			}
		}
		res, err := Insert("item", []string{"it_id", "it_pd_id", "it_bt_id", "it_status"}, val[0:4])
		if err != nil || res == 0 {
			return i, err
		}
		i += res
	}
	return int64(num), nil
}

// func SetItem(id string, change map[string]string) (res int64, err error) {
// 	var key, value []string
// 	keys := []string{"item", "it_id", "it_pd_id", "it_bt_id", "it_status"}
// 	for _, k := range keys {
// 		if v, ok := change[k]; ok {
// 			key = append(key, k)
// 			value = append(value, v)
// 		}
// 	}
// 	return Update("item", key, value, []string{"it_id ="}, []string{id})
// }
// func RemoveItem(id string) (res int64, err error) {
// 	return Delete("item", []string{"it_id ="}, []string{id})
// }

// insert delete update selete
func Insert(tb_name string, keys []string, values []string) (res int64, err error) {
	if len(keys) != len(values) {
		return 0, errors.New("请求键值数量不一致")
	}
	if len(keys) == 0 {
		return 0, errors.New("无动作")
	}
	qstr := "INSERT INTO " + tb_name + " ("
	for i, s := range keys {
		if i != 0 {
			qstr += ", "
		}
		qstr += s
	}
	qstr += ") VALUES ("
	for i, _ := range values {
		if i != 0 {
			qstr += ", "
		}
		qstr += "?"
	}
	qstr += ")"
	values_interface := make([]interface{}, len(values))
	for i, v := range values {
		values_interface[i] = v
	}
	// fmt.Println(qstr, values)
	result, err := db.Exec(qstr, values_interface...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func Delete(tb_name string, search_keys []string, search_values []string) (res int64, err error) {
	if len(search_keys) != len(search_values) {
		return 0, errors.New("请求键值数量不一致")
	}

	qstr := "DELETE FROM " + tb_name

	for i, s := range search_keys {
		if i == 0 {
			qstr += " WHERE "
		} else {
			qstr += " and "
		}
		qstr += s + " ?"
	}
	values_interface := make([]interface{}, len(search_values))
	for i, v := range search_values {
		values_interface[i] = v
	}

	//fmt.Println(qstr, search_values)
	result, err := db.Exec(qstr, values_interface...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func Update(tb_name string, keys []string, values []string, search_keys []string, search_values []string) (res int64, err error) {
	if len(keys) != len(values) {
		return 0, errors.New("请求键值数量不一致")
	}
	if len(search_keys) != len(search_values) {
		return 0, errors.New("请求键值数量不一致")
	}
	if len(keys) == 0 {
		return 0, errors.New("无动作")
	}
	qstr := "UPDATE " + tb_name + " SET "
	for i, s := range keys {
		if i != 0 {
			qstr += ", "
		}
		qstr += s + " = ?"
	}
	for i, s := range search_keys {
		if i == 0 {
			qstr += " WHERE "
		} else {
			qstr += " and "
		}
		qstr += s + " ?"
	}
	values = append(values, search_values...)
	values_interface := make([]interface{}, len(values))
	for i, v := range values {
		values_interface[i] = v
	}
	//fmt.Println(qstr, values)
	result, err := db.Exec(qstr, values_interface...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func Select(tb_name string, search_keys []string, search_values []string, keys []string) (ret Table, err error) {
	ret.Init(keys...)
	if len(search_keys) != len(search_values) {
		return ret, errors.New("请求键值数量不一致")
	}
	col_num := len(keys)
	qstr := "SELECT "
	for i, s := range keys {
		if i != 0 {
			qstr += ", "
		}
		qstr += s
	}
	qstr += " FROM " + tb_name

	for i, s := range search_keys {
		if i == 0 {
			qstr += " WHERE "
		} else {
			qstr += " and "
		}
		qstr += s + " ?"
	}
	values_interface := make([]interface{}, len(search_values))
	for i, v := range search_values {
		values_interface[i] = v
	}
	//fmt.Println(qstr, search_values)
	rows, err := db.Query(qstr, values_interface...)
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		new_row := make([]string, col_num)
		pointers := make([]interface{}, col_num)
		for i := 0; i < col_num; i++ {
			pointers[i] = &new_row[i]
		}
		if err = rows.Scan(pointers...); err != nil {
			return ret, err // Handle scan error
		}
		ret.Content = append(ret.Content, new_row)
	}
	// check iteration error
	if rows.Err() != nil {
		fmt.Println(err)
	}
	return
}
