package dto

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// ErrorJSON is wrapper struct for libindy's corresponding JSON type.
type ErrorJSON struct {
	Backtrace string `json:"backtrace,omitempty"`
	Message   string `json:"message"`
}

// NewErrorJSON creates a new ErrorJSON
func NewErrorJSON(js string) (msg *ErrorJSON) {
	if js == "" {
		js = "{}"
	}
	err := json.Unmarshal([]byte(js), &msg)
	if err != nil {
		fmt.Println("err marshalling from JSON: ", err.Error())
		return nil
	}
	return
}

// ToJSON is a helper to convert dto to JSON string.
func ToJSON(dto interface{}) string {
	return string(ToJSONBytes(dto))
}

// FromJSONStr is a helper to convert object from a JSON string.
func FromJSONStr(str string, dto interface{}) {
	FromJSON([]byte(str), dto)
}

// FromJSON is a helper to convert byte JSON to a data object.
func FromJSON(bytes []byte, dto interface{}) {
	err := json.Unmarshal(bytes, dto)
	if err != nil {
		glog.Errorf("%s: from JSON:\n%s\n", err.Error(), string(bytes))
		panic(err)
	}
}

// ToJSONBytes is a helper to convert dto to JSON byte data.
func ToJSONBytes(dto interface{}) []byte {
	output, err := json.Marshal(dto)
	if err != nil {
		fmt.Println("err marshalling to JSON:", err)
		return nil
	}
	return output
}

// JSONArray returns a JSON array in string from strings given.
func JSONArray(strs ...string) string {
	return DoJSONArray(strs)
}

// DoJSONArray returns a JSON array in string from strings given.
func DoJSONArray(strs []string) string {
	var b strings.Builder
	b.WriteRune('[')

	for i, s := range strs {
		if i > 0 {
			b.WriteRune(',')
		}
		b.WriteRune('"')
		b.WriteString(s)
		b.WriteRune('"')
	}
	b.WriteRune(']')
	return b.String()
}

// ToGOB returns bytes of the object in GOB format.
func ToGOB(dto interface{}) []byte {
	var buf bytes.Buffer
	// Create an encoder and send a value.
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dto)
	if err != nil {
		glog.Fatal("encode:", err)
	}
	return buf.Bytes()
}

// FromGOB reads object from bytes. Remember pass a pointer to preallocated
// object of the right type.
//
//	p := &PSM{}
//	dto.FromGOB(d, p)
//	return p
func FromGOB(data []byte, dto interface{}) {
	network := bytes.NewReader(data)
	dec := gob.NewDecoder(network)
	err := dec.Decode(dto)
	if err != nil {
		glog.Fatal("decode:", err)
	}
}

// Result is the return type result of the wrapper functions
type Result struct {
	Er   Err  `json:",omitempty"`
	Data Data `json:",omitempty"`
}

// Err is the error part of the Result type of the wrapper function return
// types.
type Err struct {
	Error string `json:",omitempty"`
	Code  int    `json:",omitempty"`
}

// Data is the actual data from wrapper function's return values when function
// call is successful.
type Data struct {
	Handle int    `json:"handle,omitempty"`
	Str1   string `json:"str_1,omitempty"`
	Str2   string `json:"str_2,omitempty"`
	Str3   string `json:"str_3,omitempty"`
	U64    uint64 `json:"u_64,omitempty"`
	Bytes  []byte `json:"bytes,omitempty"`
	Yes    bool   `json:"yes,omitempty"`
}

func (r Result) String() string {
	if r.Err() != nil {
		return r.Er.Error
	}
	// In JSON field names are visible
	return ToJSON(r.Data)
}

// Err returns Go error.
func (r Result) Err() error {
	if r.Er.Error != "" {
		return r
	}
	return nil
}

// SetErr sets the Go error for Result.
func (r *Result) SetErr(e error) {
	r.Er.Error = e.Error()
}

// Error is Go error method to make Result error compatible.
func (r Result) Error() string {
	// If there are no error, let String() do the job. For example printf-
	// family functions need that. They order they call interfaces are
	// 1. Error() 2. String()
	if r.Er.Error != "" {
		return r.Er.Error
	}
	return r.String()
}

// ErrCode returns an indy error code.
func (r Result) ErrCode() int {
	return r.Er.Code
}

// SetErrCode sets an indy error code.
func (r *Result) SetErrCode(c int) {
	r.Er.Code = c
}

// SetHandle sets the Handle data part of the Result.
func (r *Result) SetHandle(h int) {
	r.Data.Handle = h
}

// Handle return the Handle part of the result.
func (r Result) Handle() int {
	return r.Data.Handle
}

// SetYes sets the bool part of the Result.
func (r *Result) SetYes(v bool) {
	r.Data.Yes = v
}

// Yes return the bool part of the Result.
func (r Result) Yes() bool {
	return r.Data.Yes
}

// SetStr1 sets 1st string value of the Result.
func (r *Result) SetStr1(s string) {
	r.Data.Str1 = s
}

// Str1 returns 1st string.
func (r Result) Str1() string {
	return r.Data.Str1
}

// SetStr2 sets 1st string value of the Result.
func (r *Result) SetStr2(s string) {
	r.Data.Str2 = s
}

// Str2 returns 2nd string.
func (r Result) Str2() string {
	return r.Data.Str2
}

// SetBytes sets bytes part of the Result.
func (r *Result) SetBytes(bytes []byte) {
	r.Data.Bytes = bytes
}

// Bytes returns bytes.
func (r Result) Bytes() []byte {
	return r.Data.Bytes
}

// SetStr3 sets 3rd string value of the Result.
func (r *Result) SetStr3(s string) {
	r.Data.Str3 = s
}

// Str3 returns 3rd string.
func (r *Result) Str3() string {
	return r.Data.Str3
}

// SetU64 set corresponding part of the Result.
func (r *Result) SetU64(v uint64) {
	r.Data.U64 = v
}

// Uint64 return a corresponding part of the Result.
func (r *Result) Uint64() uint64 {
	return r.Data.U64
}

// SetErrorJSON sets error data from json string of the Result.
func (r *Result) SetErrorJSON(jsonString string) {
	ej := NewErrorJSON(jsonString)
	if ej != nil {
		r.Er.Error = ej.Message
	}
}
