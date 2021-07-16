package tools

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/zyfmix/go_tools/src/logs"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// Go 检查slice是否包含指定元素...
// [GoLang: Slice contains and dedup](https://samurailink3.com/blog/2016/05/11/golang-slice-contains-and-dedup/)

func ContainInt(intSlice []int, searchInt int) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func ContainInt64(intSlice []int64, searchInt int64) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func ContainUInt64(intSlice []uint64, searchInt uint64) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func ContainUInt8(intSlice []uint8, searchInt uint8) bool {
	for _, value := range intSlice {
		if value == searchInt {
			return true
		}
	}
	return false
}

func ContainString(strSlice []string, str string) bool {
	for _, a := range strSlice {
		if a == str {
			return true
		}
	}
	return false
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func MapKeys(mKeys map[string]interface{}) []string {
	keys := make([]string, 0, len(mKeys))
	for mKey := range mKeys {
		keys = append(keys, mKey)
	}
	return keys
}

func Dedup(intSlice []int) []int {
	var returnSlice []int
	for _, value := range intSlice {
		if !ContainInt(returnSlice, value) {
			returnSlice = append(returnSlice, value)
		}
	}
	return returnSlice
}

func DedupInt64(int64Slice []int64) []int64 {
	var returnSlice []int64
	for _, value := range int64Slice {
		if !ContainInt64(returnSlice, value) {
			returnSlice = append(returnSlice, value)
		}
	}
	return returnSlice
}

func ConvertToInterfaces(args []string) []interface{} {
	params := make([]interface{}, len(args))
	for index, arg := range args {
		params[index] = arg
	}
	return params
}

func DataUnmarshal(data []byte) interface{} {
	var params map[string]interface{}
	if err := json.Unmarshal(data, &params); err != nil {
		logs.Error(nil, "[DataUnmarshalException]", zap.String("data", string(data)), zap.Error(err))
		return err
	}
	return params
}

const (
	chars    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLen = len(chars)
	mask     = 1<<6 - 1
)

// [How to generate a random string of a fixed length in Go?](https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986)
// [快速产生一个随机字符串](https://colobu.com/2018/09/02/generate-random-string-in-Go/)
var rng = rand.NewSource(time.Now().UnixNano())

// RandStr 返回指定长度的随机字符串
func RandStr(ln int) string {
	/* chars 38个字符
	 * rng.Int63() 每次产出64bit的随机数,每次我们使用6bit(2^6=64) 可以使用10次
	 */
	buf := make([]byte, ln)
	for idx, cache, remain := ln-1, rng.Int63(), 10; idx >= 0; {
		if remain == 0 {
			cache, remain = rng.Int63(), 10
		}
		buf[idx] = chars[int(cache&mask)%charsLen]
		cache >>= 6
		remain--
		idx--
	}
	return *(*string)(unsafe.Pointer(&buf))
}

func TrySnInfo(ctx context.Context, param []byte) (string, uint64, time.Time, error) {
	type SnValidators struct {
		TraceId   string    `json:"trace_id" validate:"required"`
		NotifyId  uint64    `json:"notify_id" validate:"required"`
		MessageAt time.Time `json:"message_at" validate:"required"`
	}

	var snParams SnValidators
	if err := json.Unmarshal(param, &snParams); err != nil {
		logs.Error(ctx, "SnParamsUnmarshalError", zap.String("param", string(param)), zap.Error(err))
		return "", 0, time.Time{}, errors.New("ParamsUnmarshalError")
	}

	validate := validator.New()
	if err := validate.Struct(snParams); err != nil {
		logs.Error(ctx, "TrySnTask", zap.Any("snParams", snParams),  zap.Error(err))
		return "", 0, time.Time{}, errors.New("参数验证失败")
	}

	return snParams.TraceId, snParams.NotifyId, snParams.MessageAt, nil
}

func TrySnTaskId(ctx context.Context, param []byte) (uint64, error) {
	type TaskValidators struct {
		TaskId uint64 `json:"task_id" validate:"required"`
	}

	var taskParams TaskValidators
	if err := json.Unmarshal(param, &taskParams); err != nil {
		logs.Error(ctx, "SnParamsUnmarshalError", zap.String("param", string(param)), zap.Error(err))
		return 0, errors.New("ParamsUnmarshalError")
	}

	validate := validator.New()
	if err := validate.Struct(taskParams); err != nil {
		logs.Error(ctx, "TrySnTaskId", zap.Any("taskParams", taskParams),  zap.Error(err))
		return 0, errors.New("参数验证失败")
	}

	return taskParams.TaskId, nil
}

func TrySnTraceId(ctx context.Context, param []byte) (string, error) {
	type TraceValidators struct {
		TraceId string `json:"trace_id" validate:"required"`
	}

	var traceParams TraceValidators
	if err := json.Unmarshal(param, &traceParams); err != nil {
		logs.Error(ctx, "SnParamsUnmarshalError", zap.String("param", string(param)), zap.Error(err))
		return "", errors.New("ParamsUnmarshalError")
	}

	validate := validator.New()
	if err := validate.Struct(traceParams); err != nil {
		logs.Error(ctx, "TrySnTraceId", zap.Any("traceParams", traceParams),  zap.Error(err))
		return "", errors.New("参数验证失败")
	}

	return traceParams.TraceId, nil
}

func TrySncCallId(ctx context.Context, param []byte) (uint64, error) {
	type CallValidators struct {
		CallId uint64 `json:"call_id" validate:"required"`
	}

	var callParams CallValidators
	if err := json.Unmarshal(param, &callParams); err != nil {
		logs.Error(ctx, "SnParamsUnmarshalError", zap.String("param", string(param)), zap.Error(err))
		return 0, errors.New("ParamsUnmarshalError")
	}

	validate := validator.New()
	if err := validate.Struct(callParams); err != nil {
		logs.Error(ctx, "TrySncCallId", zap.Any("callParams", callParams),  zap.Error(err))
		return 0, errors.New("参数验证失败")
	}

	return callParams.CallId, nil
}

func TrySneCallId(ctx context.Context, param []byte) (string, error) {
	type CallValidators struct {
		Params struct {
			CallId string `json:"callId"`
		} `json:"params" validate:"required"`
	}

	var callParams CallValidators
	if err := json.Unmarshal(param, &callParams); err != nil {
		logs.Error(ctx, "SnParamsUnmarshalError", zap.String("param", string(param)), zap.Error(err))
		return "", errors.New("ParamsUnmarshalError")
	}

	validate := validator.New()
	if err := validate.Struct(callParams); err != nil {
		logs.Error(ctx, "TrySneCallId", zap.Any("callParams", callParams),  zap.Error(err))
		return "", errors.New("参数验证失败")
	}

	return callParams.Params.CallId, nil
}

// [What is the correct way to find the min between two integers in Go?](https://stackoverflow.com/questions/27516387/what-is-the-correct-way-to-find-the-min-between-two-integers-in-go)
func MinOf(vars ...int64) int64 {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

func MaxOf(vars ...int64) int64 {
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func MaxOfUInt64(vars ...uint64) uint64 {
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

func ArrayItemDelete(s []int64, item int64) []int64 {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

// [聊一聊,Golang “相对”路径问题](https://studygolang.com/articles/12563)
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}
