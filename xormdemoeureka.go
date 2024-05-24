package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hudl/fargo"
	"log"
	"net/http"
	"strconv"
	"time"
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

	// 设置Eureka连接
	eureka := fargo.NewConn("http://120.46.54.84:8761/eureka")
	// 定义要注册的实例
	instance := fargo.Instance{
		HostName:         "172.20.10.2",
		Port:             9090,
		App:              "MyGoApp",
		IPAddr:           "127.0.0.1",
		VipAddress:       "MyGoApp",
		SecureVipAddress: "MyGoApp",
		HealthCheckUrl:   "http:// 172.20.10.2:9090/health",
		StatusPageUrl:    "http:// 172.20.10.2:9090/info",
		HomePageUrl:      "http:// 172.20.10.2:9090/",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
	}
	// 注册实例到Eureka
	err := eureka.RegisterInstance(&instance)
	if err != nil {
		log.Fatal("Eureka Registration Failed: ", err)
	}
	// 启动心跳机制保持注册状态
	go func() {
		for {
			err := eureka.HeartBeatInstance(&instance)
			if err != nil {
				log.Println("Eureka Heartbeat Failed: ", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

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
	r.GET("/articleslist", func(c *gin.Context) {
		articles, err := getAllArticles(engine)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve articles"})
			return
		}
		c.JSON(http.StatusOK, articles)
	})
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
	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		// 检查您的服务是否健康，例如数据库连接是否正常等。
		// 如果一切正常，返回HTTP 状态200。
		c.Status(http.StatusOK)
	})
	r.GET("/info", func(c *gin.Context) {
		// 提供服务的相关信息，如版本、描述等。
		// 以下只是一个示例，您应根据实际情况返回相应的状态信息。
		c.JSON(http.StatusOK, gin.H{
			"version":     "1.0.0",
			"description": "MyGoApp service status information",
		})
	})
	r.Run(":9090")
	// 在程序退出前注销服务
	defer func() {
		err := eureka.DeregisterInstance(&instance)
		if err != nil {
			log.Fatal("Eureka Deregistration Failed: ", err)
		}
	}()
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
func getAllArticles(engine *xorm.Engine) ([]Article, error) {
	var articles []Article
	err := engine.Find(&articles)
	return articles, err
}
