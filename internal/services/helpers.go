package services

import "database/sql"

func dbQueryGetPeopleByName(firstName string, lastName string, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person" 
					WHERE first_name = $1
					AND last_name = $2`,
		firstName,
		lastName)
}
func dbQueryGetPeopleByNameAndAge(firstName string, lastName string, age int, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person" 
					WHERE first_name = $1
					AND last_name = $2
					AND age = $3`,
		firstName,
		lastName,
		age)
}
func dbQueryGetPeopleByAge(age int, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person" 
					WHERE age = $1`,
		age)
}
func dbQueryGetPeople(db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person"`)
}

// returns an []int of values that are in old, but not in new. New may contain values not in old, which should be ignored. items in old are unique.
func getDifference(old []int, new []int) []int {
	//we add all old ones to the map with a value of 2.
	//we then add all new ones to the map with a value of 1
	//all keys in the list equal to 2 are in old but not new
	list := make(map[int]int)
	for _, val := range old {
		list[val] += 2
	}
	for _, val := range new {
		list[val] += 1
	}
	result := make([]int, 0)
	for key, val := range list {
		if val == 2 {
			result = append(result, key)
		}
	}
	return result
}
