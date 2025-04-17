package middleware

import (
	"bytes"
	"errors"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"app/app/controller/activitylog"
	"app/app/helper"
	"app/app/model"
	"app/app/response"
	"app/config"
)

const (
	LocalOrigin  = "LC-Origin"
	LocalCountry = "LC-COUNTRY"
	LocalCFRay   = "LC-CP-RAY"
	LocalIP      = "LC-IP"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type LogResponseInfo struct {
	Method    string
	Path      string
	IP        string
	UserAgent string
	Header    any
	Query     any
	Request   string
	Response  string
}

func NewLogResponse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newLog := new(LogResponseInfo)
		newLog.Method = ctx.Request.Method
		newLog.UserAgent = ctx.Request.UserAgent()
		newLog.Path = ctx.FullPath()
		newLog.IP = ctx.ClientIP()
		// newLog.Header = ctx.Request.Header
		newLog.Query = ctx.Request.URL.Query()

		// logger.Infof("Request: %s ", newLog.Header)

		// Set Header value
		ctx.Set(LocalOrigin, GetHeader(ctx, `Origin`))
		ctx.Set(LocalCountry, GetHeader(ctx, `CF-IPCountry`))
		ctx.Set(LocalCFRay, GetHeader(ctx, `CF-RAY`))
		ctx.Set(LocalIP, GetHeader(ctx, `CF-Connecting-IP`))

		// GET Data Body
		body, err := io.ReadAll(ctx.Request.Body)
		if errors.Is(err, io.EOF) {

		} else if err != nil {
			response.InternalError(ctx, err.Error())
			ctx.Abort()
			return
		} else {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			newLog.Request = string(body)
		}

		// Set struct Resposne
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		// Next Process
		ctx.Next()

		// Get Response Body
		resBody := blw.body.String()
		newLog.Response = string(resBody)

		statusCode := ctx.Writer.Status()
		// Check 404 and redirect to 403
		if statusCode == 404 {
			response.Forbidden(ctx, nil)
			ctx.Abort()
			return
		}

		user, err := helper.GetUserByToken(ctx)
		if err != nil {
			response.Unauthorized(ctx, nil)
			ctx.Abort()
			return
		}
		db := config.GetDB()
		acsv := activitylog.NewService(db)
		logTemp := model.ActivityLog{
			Section:    newLog.Path,
			EventType:  newLog.Method,
			StatusCode: statusCode,
			Responses:  newLog.Response,
			Parameters: newLog.Request,
			Query:      newLog.Query,
			IpAddress:  newLog.IP,
			UserAgent:  newLog.UserAgent,
			CreatedBy:  user,
			CreatedAt:  time.Now().Unix(),
		}
		_, err = acsv.Create(ctx, logTemp)
		if err != nil {
			response.InternalError(ctx, err.Error())
			ctx.Abort()
			return
		}
	}
}

func GetHeader(ctx *gin.Context, key string) string {
	val, ok := ctx.Get(LocalIP)
	if !ok {
		return `not-found`
	}
	return val.(string)
}
