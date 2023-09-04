package goqp

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strconv"
)

type QueryParser[T interface{}] struct {
	Error  error
	Query  *url.Values
	Data   *T
	params map[string]string
}
type QueryParserCheckFn[T interface{}] func(d *T) error
type QueryParserParamFn[T interface{}, U interface{}] func(v U, d *T)
type QueryParserParamsFn[T interface{}, U interface{}] func(v []U, d *T)
type QueryParserParamErrorFn[T interface{}, U interface{}] func(v U, d *T) error
type QueryParserCustomErrorFn[T interface{}] func(d *T) error
type QueryParserExtrasFn[T interface{}] func(extras map[string]string, d *T)

//goland:noinspection GoUnusedExportedFunction
func NewQueryParser[T interface{}](q *url.Values, d *T) *QueryParser[T] {
	qp := QueryParser[T]{
		Query:  q,
		Data:   d,
		Error:  nil,
		params: map[string]string{},
	}

	return &qp
}

func (qp *QueryParser[T]) RegisterParam(key string, typeName string) *QueryParser[T] {
	qp.params[key] = typeName

	return qp
}
func (qp *QueryParser[T]) HasError() bool {
	return qp.Error != nil
}
func (qp *QueryParser[T]) Fn(key string, fn QueryParserParamFn[T, string]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.RegisterParam(key, "fn")
	v := qp.Query.Get(key)
	fn(v, qp.Data)
	return qp
}
func (qp *QueryParser[T]) CustomErrorFn(fn QueryParserCustomErrorFn[T]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.Error = fn(qp.Data)
	return qp
}
func (qp *QueryParser[T]) ErrorFn(key string, fn QueryParserParamErrorFn[T, string]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.RegisterParam(key, "fn")
	v := qp.Query.Get(key)
	qp.Error = fn(v, qp.Data)
	return qp
}
func (qp *QueryParser[T]) String(key string, defaultValue string, fn QueryParserParamFn[T, string]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.RegisterParam(key, "string")
	v := qp.Query.Get(key)
	if len(v) == 0 {
		v = defaultValue
	}
	fn(v, qp.Data)
	return qp
}
func (qp *QueryParser[T]) Base64String(key string, defaultValue string, fn QueryParserParamFn[T, string]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.RegisterParam(key, "string")
	v := qp.Query.Get(key)
	if len(v) == 0 {
		v = defaultValue
	}
	if len(v) != 0 {
		decoded, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			v = ""
		} else {
			v = string(decoded)
		}
	}
	fn(v, qp.Data)
	return qp
}

func (qp *QueryParser[T]) Int(key string, defaultValue int, fn QueryParserParamFn[T, int]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	qp.RegisterParam(key, "int")
	vv := qp.Query.Get(key)
	v := defaultValue
	if len(vv) > 0 {
		vvv, err := strconv.Atoi(vv)
		if err == nil {
			v = vvv
		}
	}
	fn(v, qp.Data)
	return qp
}
func (qp *QueryParser[T]) Ints(keys []string, defaultValues []int, fn QueryParserParamsFn[T, int]) *QueryParser[T] {
	if qp.HasError() {
		return qp
	}
	if len(keys) != len(defaultValues) {
		qp.Error = errors.New("non-matching size for int params")
		return qp
	}
	values := make([]int, len(keys))

	for i, key := range keys {
		qp.RegisterParam(key, "int")
		vv := qp.Query.Get(key)
		v := defaultValues[i]
		if len(vv) > 0 {
			vvv, err := strconv.Atoi(vv)
			if err == nil {
				v = vvv
			}
		}
		values[i] = v
	}

	fn(values, qp.Data)
	return qp
}
func (qp *QueryParser[T]) Parse(fn QueryParserCheckFn[T]) error {
	if qp.HasError() {
		return qp.Error
	}
	qp.Error = fn(qp.Data)

	return qp.Error
}
func (qp *QueryParser[T]) Extras(fn QueryParserExtrasFn[T]) *QueryParser[T] {
	cleanedParams := make(map[string]string)
	for k, v := range *qp.Query {
		if _, ok := qp.params[k]; !ok {
			if len(k) <= 20 && len(v[0]) <= 255 {
				cleanedParams[k] = v[0]
			}
		}
	}
	fn(cleanedParams, qp.Data)
	return qp
}
