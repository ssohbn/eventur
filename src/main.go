package main

import (
	"net/http"

  "github.com/gin-gonic/gin"

  "html/template"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  //front end routes
  templ := template.Must(template.ParseFiles(
    "src/templates/layout.html",
    "src/templates/index.html",
    "src/templates/partials/navbar.html", 
  ))
  r.SetHTMLTemplate(templ)

  r.GET("/", func(c *gin.Context) {
      c.HTML(http.StatusOK, "layout.html", gin.H{
          "Title": "Home Page",
      })
  })
  // run the server
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
