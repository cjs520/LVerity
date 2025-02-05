package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"time"
)

// Logger 中间件，用于记录请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后 body 会被清空
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 获取响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		// 执行时间
		latency := end.Sub(start)

		// 请求信息
		method := c.Request.Method
		uri := c.Request.RequestURI
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		queryParams := c.Request.URL.Query()

		// 构建日志数据
		logData := map[string]interface{}{
			"timestamp":    end.Format("2006/01/02 - 15:04:05"),
			"status":      status,
			"latency":     latency.String(),
			"client_ip":   clientIP,
			"method":      method,
			"uri":         uri,
			"query":       queryParams,
			"request_id":  c.GetHeader("X-Request-ID"),
		}

		// 添加请求体（如果存在且是 JSON）
		if len(requestBody) > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, requestBody, "", "  "); err == nil {
				logData["request_body"] = prettyJSON.String()
			}
		}

		// 添加响应体（如果存在且是 JSON）
		if blw.body.Len() > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, blw.body.Bytes(), "", "  "); err == nil {
				logData["response_body"] = prettyJSON.String()
			}
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logData["errors"] = c.Errors.String()
		}

		// 转换为 JSON 并输出
		jsonLog, _ := json.Marshal(logData)
		log.Println(string(jsonLog))
	}
}

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
