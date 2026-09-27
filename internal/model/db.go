package model

type Database struct {
	results map[string]string
}

func CreateDatabase() Database {
	var database = &Database{}
	database.results = make(map[string]string)
	return *database
}

func (db *Database) Put(key string, value string) {
	if key == "" || value == "" {
		return
	}

	db.results[key] = value
}

func (db *Database) Get(key string) (bool, string) {
	if key == "" {
		return false, ""
	}

	res, ok := db.results[key]

	if ok {
		return true, res
	}

	return false, ""
}
