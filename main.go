package main

/*
#include <stdlib.h>
#cgo LDFLAGS: -L ./libs -lvisual_decaf

extern int get_id();
extern char* upload_code(char* code, int id);
extern char* get_token_stream(int id);
extern char* get_ast(int id);
extern char* get_program(int id);
extern char* get_debug_info(int id);
extern void end_use(int id);
*/
import "C"

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"unsafe"
)

func main() {
	server := gin.Default()

	// 设置允许跨域访问
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true

	server.Use(cors.New(config))

	server.GET("/enter/id", enterHandler)
	server.POST("/code", postCodeHandler)
	server.GET("/tokens", tokenStreamHandler)
	server.GET("/ast", astHandler)
	server.GET("/program", programHandler)
	server.GET("/debug", debugHandler)
	server.DELETE("/:id", endUseHandler)

	if err := server.Run(":8080"); err != nil {
		fmt.Println("服务器启动错误！")
	}
}

// 处理请求id
func enterHandler(c *gin.Context) {
	id := int(C.get_id())
	response := "{\"code\":\"1\",\"msg\":\"Success\",\"result\":" + strconv.Itoa(id) + "}"
	c.String(200, response)
}

// 处理提交代码请求
func postCodeHandler(c *gin.Context) {
	code := c.PostForm("code")
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))
	cResponse := C.upload_code(cCode, C.int(nid))
	defer C.free(unsafe.Pointer(cResponse))
	goResponse := C.GoString(cResponse)
	c.String(200, goResponse)
}

// 处理token流请求
func tokenStreamHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cTokenStream := C.get_token_stream(C.int(nid))
	defer C.free(unsafe.Pointer(cTokenStream))
	tokenStream := C.GoString(cTokenStream)
	c.String(200, tokenStream)
}

// 处理ast请求
func astHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cAST := C.get_ast(C.int(nid))
	defer C.free(unsafe.Pointer(cAST))
	ast := C.GoString(cAST)
	c.String(200, ast)
}

// 处理program请求
func programHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cProgram := C.get_program(C.int(nid))
	defer C.free(unsafe.Pointer(cProgram))
	program := C.GoString(cProgram)
	c.String(200, program)
}

// 处理调试请求
func debugHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cDebugInfo := C.get_debug_info(C.int(nid))
	defer C.free(unsafe.Pointer(cDebugInfo))
	debugInfo := C.GoString(cDebugInfo)
	c.String(200, debugInfo)
}

// 处理退出请求
func endUseHandler(c *gin.Context) {
	id := c.Param("id")
	nid, _ := strconv.Atoi(id)
	C.end_use(C.int(nid))
	c.String(200, "Success")
}
