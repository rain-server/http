package main

import (
	"bytes"
	"context"
	"fmt"
	xhttp "github.com/goclub/http"
	"io/ioutil"
	"log"
	"net/http"
)

func main () {
	router := NewRouter()
	// request
	RequestBindQuery(router)
	RequestBindFormUrlencoded(router)
	RequestBindFormData(router)
	RequestBindJSON(router)
	RequestBindQueryAndJSON(router)
	RequestBindQueryAndJSON(router)
	RequestBindParam(router)
	RenderFormFile(router)
	RequestFile(router)
	// response
	ResponseBytes(router)
	addr := ":3000"
	serve := http.Server{
		Handler: router,
		Addr: addr,
	}
	log.Print("http://127.0.0.1" + addr)
	router.LogPatterns()
	go func() {
		listenErr := serve.ListenAndServe() ; if listenErr !=nil {
			if listenErr != http.ErrServerClosed {
				panic(listenErr)
			}
		}
	}()
	xhttp.GracefulClose(func() {
		log.Print("Shuting down server...")
		if err := serve.Shutdown(context.Background()); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
		log.Println("Server exiting")
	})
}

func NewRouter() *xhttp.Router {
	router := xhttp.NewRouter(xhttp.RouterOption{
		OnCatchError: func(c *xhttp.Context, errInterface interface{}) error {
			if errInterface != nil {
				panic(errInterface)
			}
			return nil
		},
	})
	return router
}



// 打开 http://127.0.0.1:3000/request/query?name=nimoc&age=18
func RequestBindQuery(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/request/query",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			Name string `query:"name"`
			Age int `query:"age"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}

func RequestBindFormUrlencoded(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.POST, "/request/form_urlencoded",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			Name string `form:"name"`
			Age int `form:"age"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}
/*
使用 curl 发起请求
curl --location --request POST 'http://127.0.0.1:3000/request/form_data' \
--form 'name="nimoc"' \
--form 'age="18"'
*/
func RequestBindFormData(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.POST, "/request/form_data",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			Name string `form:"name"`
			Age int `form:"age"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}
/*
使用 curl 发起请求
curl --location --request GET 'http://127.0.0.1:3000/request/json' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"nimoc",
    "age": 18
}'
*/
func RequestBindJSON(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/request/json",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			Name string `json:"name"`
			Age int `json:"age"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}


/*
curl --location --request GET 'http://127.0.0.1:3000/request/query_and_json?id=11' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"nimoc",
    "age": 18
}'
*/
func RequestBindQueryAndJSON(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/request/query_and_json",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			ID string `query:"id"`
			Name string `json:"name"`
			Age int `json:"age"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}

// http://127.0.0.1:3000/request/param/11
func RequestBindParam(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/request/param/{userID}",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		req := struct {
			UserID string `param:"userID"`
		}{}
		reject = c.BindRequest(&req) ; if reject != nil {return}
		dump := fmt.Sprintf("%+v", req)
		return c.Bytes([]byte(dump))
	})
}
func RenderFormFile(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/request/file",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		html := []byte(`
		<form action="/request/file" method="post" enctype="multipart/form-data" >
　　　		<input type="file" name="file" /> <br />
　　　		<button type="submit" >上传</button>
		</form>`)
		return c.Render(func(buffer *bytes.Buffer) {
			buffer.Write(html)
		})
	})
}
// 打开 http://127.0.0.1:3000/request/file 上传文件
func RequestFile(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.POST, "/request/file",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		file, fileHeader, reject := c.Request.FormFile("file") ; if reject != nil {return}
		defer file.Close()
		data, reject := ioutil.ReadAll(file) ; if reject != nil {return}
		body :=  append([]byte(fileHeader.Filename + ":"), data...)
		return c.Bytes(body)
	})
}
func ResponseBytes(router *xhttp.Router) {
	pattern := xhttp.Pattern{
		xhttp.GET, "/response/bytes",
	}
	router.HandleFunc(pattern, func(c *xhttp.Context) (reject error) {
		return c.Bytes([]byte("goclub"))
	})
}