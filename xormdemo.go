package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type Article struct {
	ID        int        `xorm:"'id' autoincr pk" json:"id"`
	Title     string     `xorm:"varchar(255)" json:"title"`
	Content   string     `xorm:"text" json:"content"`
	ViewNum   int        `xorm:"int" json:"view_num"`
	CreatedAt time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted" json:"deleted_at"`
}

func (a *Article) TableName() string {
	return "articles"
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Close()
	// 同步数据库结构
	err = engine.Sync2(new(Article))
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	r.GET("/articles/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		article, err := getArticle(engine, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		fmt.Println(article)
		c.JSON(http.StatusOK, article)
	})

	r.POST("/articles", func(c *gin.Context) {
		var article Article
		if err := c.BindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id := insertArticle(engine, article)
		article.ID = id
		c.JSON(http.StatusCreated, article)
	})

	r.PUT("/articles/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var article Article
		if err := c.BindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updateArticle(engine, id, article)
		c.JSON(http.StatusOK, gin.H{"success": "Article updated"})
	})

	r.DELETE("/articles/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		deleteArticle(engine, id)
		c.JSON(http.StatusOK, gin.H{"success": "Article deleted"})
	})

	r.Run(":9090")
}

// 使用xorm重写数据库操作函数
func insertArticle(engine *xorm.Engine, article Article) int {
	_, err := engine.Insert(&article)
	if err != nil {
		log.Println("Insert error:", err)
		return 0
	}
	return article.ID // 返回插入记录的ID
}

func getArticle(engine *xorm.Engine, id int) (Article, error) {
	var article Article
	has, err := engine.ID(id).Get(&article)
	if err != nil {
		return article, err
	}
	if !has {
		return article, sql.ErrNoRows
	}
	return article, nil
}

func updateArticle(engine *xorm.Engine, id int, article Article) error {
	_, err := engine.ID(id).Update(&article)
	return err
}

func deleteArticle(engine *xorm.Engine, id int) {
	_, err := engine.ID(id).Delete(&Article{})
	if err != nil {
		log.Fatal(err)
	}
}
