package api

import (
	"ginBlog/sql"

	"github.com/gin-gonic/gin"
)

type Route func(*gin.Engine)
type Group func(*gin.Engine)
type GroupRoute func(*gin.RouterGroup)

var Routes = []Route{}
var Groups = []Group{}
var DB *sql.DB

// var A chan int

func init() {
	// A = make(chan int)
	// A <- 5
}
