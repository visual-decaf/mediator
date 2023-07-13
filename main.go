package main

/*
#include <stdlib.h>
#cgo LDFLAGS: -L ./libs -lvisual_decaf

extern int get_id();
extern void compile(char* code, int id);
extern char* get_token_stream(int id);
extern char* get_ast(int id);
extern char* get_program(int id);
*/
import "C"

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"unsafe"
)

func main() {
	server := gin.Default()

	server.GET("/enter", enterHandler)
	server.POST("/code", postCodeHandler)
	server.GET("/tokens", tokenStreamHandler)
	server.GET("/ast", astHandler)
	server.GET("/program", programHandler)
	server.POST("/next", nextStepHandler)

	if err := server.Run(":8080"); err != nil {
		fmt.Println("服务器启动错误！")
	}
}

func enterHandler(c *gin.Context) {
	id := int(C.get_id())
	writeSuccessResult(c, strconv.Itoa(id))
}

// 处理提交代码请求
func postCodeHandler(c *gin.Context) {
	code := c.PostForm("code")
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))
	C.compile(cCode, C.int(nid))
	writeSuccessResult(c, "成功")
}

// 处理token流请求
func tokenStreamHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cTokenStream := C.get_token_stream(C.int(nid))
	defer C.free(unsafe.Pointer(cTokenStream))
	tokenStream := C.GoString(cTokenStream)
	writeSuccessResult(c, tokenStream)
}

// 处理ast请求
func astHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cAST := C.get_ast(C.int(nid))
	defer C.free(unsafe.Pointer(cAST))
	ast := C.GoString(cAST)
	writeSuccessResult(c, ast)
}

// 处理program请求
func programHandler(c *gin.Context) {
	id := c.Query("id")
	nid, _ := strconv.Atoi(id)
	cProgram := C.get_program(C.int(nid))
	defer C.free(unsafe.Pointer(cProgram))
	program := C.GoString(cProgram)
	writeSuccessResult(c, program)
}

// 处理调试请求
func nextStepHandler(c *gin.Context) {
	writeSuccessResult(c, "nothing")
}

// 返回成功结果
func writeSuccessResult(c *gin.Context, result string) {
	var form = "{\"code\": 1,\n\"msg\": \"操作成功\",\n\"result\": %s\n}"
	c.String(200, form, result)
}
