package constants

const (
	QUERY_INSERT_USER          = `INSERT INTO users(first_name, last_name, age, email, user_name, password, phone_number, photo, city, gps, date_created, status, gender, app_language) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
	QUERY_GET_USER_ID          = "SELECT id, first_name, last_name, age, email, user_name, phone_number, photo, city, gps, date_created, status, gender, app_language FROM users WHERE id=?;"
	QUERY_GET_LIKED_PROPERTY   = "select property_id FROM liked_properties WHERE user_id=?"
	QUERY_GET_ALL_USERS        = "SELECT id, first_name, last_name, age, email, user_name, phone_number, photo, city, gps, date_created, status, gender, app_language FROM users;"
	QUERY_UPDATE_USER          = "UPDATE users SET first_name=?, last_name=?, age=?, email=?, user_name=?, phone_number=?, city=?, gps=?, status=?, gender=?, app_language=? WHERE id=?;"
	QUERY_UPDATE_USER_PHOTO    = "UPDATE users SET photo=? WHERE id=?;"
	QUERY_DELETE_USER          = "UPDATE users set status=? where id=?"
	QUERY_INSERT_LIKE_PROPERTY = "INSERT INTO liked_properties (property_id, user_id) values(?,?);"
)
