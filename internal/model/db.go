package model

type database struct {
	results map[string]string
}

var db = &database{}

func Put(key string, value string) {
	if key == "" || value == "" {
		return
	}

	if db.results == nil {
		db.results = make(map[string]string)
	}

	db.results[key] = value
}

func Get(key string) (bool, string) {
	if key == "" || db.results == nil {
		return false, ""
	}

	res, ok := db.results[key]

	if ok {
		return true, res
	}

	return false, ""
}
