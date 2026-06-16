package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dani-susanto/go-common/validator"
	"github.com/lib/pq"
)

type QueryOperator string

const (
	QueryOpEqual              QueryOperator = "="
	QueryOpNotEqual           QueryOperator = "!="
	QueryOpGreaterThan        QueryOperator = ">"
	QueryOpGreaterThanOrEqual QueryOperator = ">="
	QueryOpLessThan           QueryOperator = "<"
	QueryOpLessThanOrEqual    QueryOperator = "<="
	QueryOpLike               QueryOperator = "LIKE"
	QueryOpILike              QueryOperator = "ILIKE"
	QueryOpIn                 QueryOperator = "IN"
	QueryOpNotIn              QueryOperator = "NOT IN"
	QueryOpIsNull             QueryOperator = "IS NULL"
	QueryOpIsNotNull          QueryOperator = "IS NOT NULL"
	QueryOpBetween            QueryOperator = "BETWEEN"
)

type QueryOrderDirection string

const (
	QueryOrderAsc  QueryOrderDirection = "ASC"
	QueryOrderDesc QueryOrderDirection = "DESC"
)

type QueryWhere struct {
	Field    string
	Operator QueryOperator
	Value    any
}

type QueryUpdate struct {
	Field string
	Value any
}

type QueryOrder struct {
	Field     string
	Direction QueryOrderDirection
}

func BuildQueryWhere(clauses []QueryWhere, args *[]any) string {
	var parts []string

	for _, c := range clauses {
		if c.Operator != QueryOpIsNull && c.Operator != QueryOpIsNotNull {
			if validator.IsEmpty(c.Value) {
				continue
			}
		}

		prefix := ""
		if len(parts) == 0 {
			prefix = "WHERE "
		}

		switch c.Operator {
		case QueryOpIsNull, QueryOpIsNotNull:
			parts = append(parts, fmt.Sprintf("%s%s %s", prefix, c.Field, c.Operator))

		case QueryOpIn, QueryOpNotIn:
			*args = append(*args, pq.Array(c.Value))
			if c.Operator == QueryOpIn {
				parts = append(parts, fmt.Sprintf("%s%s = ANY($%d)", prefix, c.Field, len(*args)))
			} else {
				parts = append(parts, fmt.Sprintf("%s%s != ALL($%d)", prefix, c.Field, len(*args)))
			}

		case QueryOpBetween:
			v := reflect.ValueOf(c.Value)
			if (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) || v.Len() != 2 {
				continue
			}
			*args = append(*args, v.Index(0).Interface())
			from := len(*args)
			*args = append(*args, v.Index(1).Interface())
			to := len(*args)
			parts = append(parts, fmt.Sprintf("%s%s BETWEEN $%d AND $%d", prefix, c.Field, from, to))

		default:
			*args = append(*args, c.Value)
			parts = append(parts, fmt.Sprintf("%s%s %s $%d", prefix, c.Field, c.Operator, len(*args)))
		}
	}

	return strings.Join(parts, " AND ")
}

func BuildQueryUpdate(clauses []QueryUpdate, args *[]any) string {
	var parts []string

	for _, c := range clauses {
		if validator.IsEmpty(c.Value) {
			continue
		}
		*args = append(*args, c.Value)
		parts = append(parts, fmt.Sprintf("%s = $%d", c.Field, len(*args)))
	}

	if len(parts) == 0 {
		return ""
	}

	return "SET " + strings.Join(parts, ", ")
}

func BuildQueryOrder(clauses []QueryOrder) string {
	var parts []string

	for _, c := range clauses {
		if c.Field == "" {
			continue
		}
		dir := c.Direction
		if dir != QueryOrderAsc && dir != QueryOrderDesc {
			dir = QueryOrderAsc
		}
		parts = append(parts, fmt.Sprintf("%s %s", c.Field, dir))
	}

	if len(parts) == 0 {
		return ""
	}

	return "ORDER BY " + strings.Join(parts, ", ")
}
