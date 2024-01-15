package main

/*import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ViewNum   int        `json:"view_num"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/articles/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		article, err := getArticle(db, id)
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
		id := insertArticle(db, article)
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
		updateArticle(db, id, article)
		c.JSON(http.StatusOK, gin.H{"success": "Article updated"})
	})

	r.DELETE("/articles/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		deleteArticle(db, id)
		c.JSON(http.StatusOK, gin.H{"success": "Article deleted"})
	})

	r.Run(":8080")
}

// 后续是各种数据库操作函数，如之前定义的 insertArticle, getArticle, updateArticle, deleteArticle
func insertArticle(db *sql.DB, article Article) int {
	query := `INSERT INTO articles (title, content, view_num,created_at,updated_at) VALUES (?, ?, ?,now(),now())`
	result, err := db.Exec(query, article.Title, article.Content, article.ViewNum)
	if err != nil {
		log.Fatal(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(insertId)
}
func getArticle(db *sql.DB, id int) (Article, error) {
	var article Article
	query := `SELECT id, title, content, view_num, created_at, updated_at, deleted_at FROM articles WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Content, &article.ViewNum, &article.CreatedAt, &article.UpdatedAt, &article.DeletedAt)
	if err != nil {
		return article, err
	}

	return article, nil
}

func updateArticle(db *sql.DB, id int, article Article) error {
	query := `UPDATE articles SET title = ?, content = ?, view_num = ? WHERE id = ?`
	_, err := db.Exec(query, article.Title, article.Content, article.ViewNum, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteArticle(db *sql.DB, id int) {
	query := `DELETE FROM articles WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		log.Fatal(err)
	}
}
*/
