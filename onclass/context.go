package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

func (c *Context) ReadJson(req interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}

	// 反序列化
	err = json.Unmarshal(body, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, data interface{}) error {
	// 序列化
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.W.WriteHeader(code)
	_, err = c.W.Write(bytes)
	return err
}

func (c *Context) OkJson(data interface{}) error {
	return c.WriteJson(http.StatusOK, data)
}

func (c *Context) SystemErrorJson(data interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, data)
}

func (c *Context) BadRequestJson(data interface{}) error {
	return c.WriteJson(http.StatusBadRequest, data)
}
