// context 中不能使用 global 中的方法打印日志, global 会调用 context 的方法,会陷入循环

package logs

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var (
	// CtxLoggerName define the ctx logger name
	CtxLoggerName = "ctx_logger"
	// CtxTraceKey define the trace id key
	CtxTraceKey = "trace"
)

// CtxLogger get the ctxLogger in context
func CtxLogger(c context.Context, fields ...zap.Field) *zap.Logger {
	// context
	if c == nil {
		c = context.Background()
		//fmt.Println("[ctx.go][CtxLogger][c == nil]")
	}

	// ctxLogger
	var ctxLogger interface{}
	if gc, ok := c.(*gin.Context); ok {
		ctxLogger, _ = gc.Get(CtxLoggerName)
		//fmt.Printf("[ctx.go][CtxLogger][detected gin.Context]CtxLoggerName: %s\n", CtxLoggerName)
	} else {
		ctxLogger = c.Value(CtxLoggerName)
		//fmt.Printf("[ctx.go][CtxLogger][detected context.Context]CtxLoggerName: %s\n", CtxLoggerName)
	}

	// zapLogger
	var zapLogger *zap.Logger
	if ctxLogger != nil {
		zapLogger = ctxLogger.(*zap.Logger)
		//fmt.Printf("[ctx.go][get zapLogger][ctxLogger != nil][from ctxLogger]\n")
	} else {
		_, zapLogger = NewCtxLogger(c, CloneLogger(CtxLoggerName), CtxTraceId(c))
		//fmt.Printf("[ctx.go][get zapLogger][ctxLogger == nil][new ctxLogger]\n")
	}

	// addition fields
	if len(fields) > 0 {
		zapLogger = zapLogger.With(fields...)
	}

	return zapLogger
}

// NewCtxLogger return a context with logger and trace id and a logger with trace id
func NewCtxLogger(c context.Context, logger *zap.Logger, traceId string) (context.Context, *zap.Logger) {
	//fmt.Printf("[NewCtxLogger][traceId: %s]\n", traceId)

	// context
	if c == nil {
		c = context.Background()
		//fmt.Println("[NewCtxLogger][c == nil]")
	}

	// detect trace id
	if traceId == "" {
		traceId = CtxTraceId(c)
	}

	// create ctx logger with trace id field
	ctxLogger := logger.With(zap.String(CtxTraceKey, traceId))

	// set data in gin.Context
	if gc, ok := c.(*gin.Context); ok {
		// set ctxLogger in gin.Context
		gc.Set(CtxLoggerName, ctxLogger)
		// set traceId in gin.Context
		gc.Set(CtxTraceKey, traceId)
	}

	// set ctxLogger in context.Context
	c = context.WithValue(c, CtxLoggerName, ctxLogger)
	// set traceId in context.Context
	c = context.WithValue(c, CtxTraceKey, traceId)

	return c, ctxLogger
}

// CtxTraceId get trace id from context
func CtxTraceId(c context.Context) string {
	// context
	if c == nil {
		c = context.Background()
	}

	// first get from gin context
	if gc, ok := c.(*gin.Context); ok {
		if traceId := gc.GetString(CtxTraceKey); traceId != "" {
			return traceId
		}
	}

	// get from go context
	traceId := c.Value(CtxTraceKey)
	if traceId != nil {
		return traceId.(string)
	}

	// uuid
	uuid4, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("uuid generate error: %#v\n", err)
	}

	// return default value
	return uuid4.String()
}
