package testutil

/*mock_course.go is used for testing purposes in
../services/course_test.go,
../services/person_test.go,
../services/setup_test.go,
*/
import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
)

func MustStructsToRows[T any](slice []T) *sqlmock.Rows {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Sprintf("expected a slice but got %T", slice))
	}

	if v.Len() == 0 {
		panic("slice is empty")
	}

	elemType := reflect.TypeOf(slice).Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected a slice of structs but got a slice of %v", elemType.Kind()))
	}

	numFields := elemType.NumField()
	columns := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		colName := elemType.Field(i).Name
		colNameSnake := toSnake(colName)
		columns[i] = colNameSnake
	}

	rows := sqlmock.NewRows(columns)

	for i := 0; i < v.Len(); i++ {
		var values []driver.Value
		elem := v.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			values = append(values, elem.Field(j).Interface())
		}
		rows.AddRow(values...)
	}

	return rows
}

func toSnake(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if (i != 0 || i == l-1) &&
			((i > 0 && rune(camel[i-1]) >= 'a') || (i < l-1 && rune(camel[i+1]) >= 'a')) {
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}
