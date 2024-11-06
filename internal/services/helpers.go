package services

//helpers.go defines helper functions used by ./course.go and ./person.go.

import "database/sql"

func dbQueryGetPeopleByName(firstName string, lastName string, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person" 
					WHERE LOWER(first_name) = LOWER($1)
					AND LOWER(last_name) = LOWER($2)`,
		firstName,
		lastName)
}
func dbQueryGetPeopleByNameAndAge(firstName string, lastName string, age int, db *sql.DB) (*sql.Rows, error) {
	return db.Query(`SELECT * FROM "person" 
					WHERE LOWER(first_name) = LOWER($1)
					AND LOWER(last_name) = LOWER($2)
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

// returns an []int of values that are in old, but not in new. New may contain values not in old. it is assumed items in old are unique
func getDifference(old []int, new []int) []int {
	//1. we increment all values in old as keys to a map[int][int] with a value of 2.
	//2. we add all values in new to the map with a value of 1
	//3. all keys in the list equal to 2 are in old but not new. return a slice of them.
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
