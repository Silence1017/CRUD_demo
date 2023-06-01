package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/crud_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 解决建表时会自动添加复数的问题，比如list变成lists
			SingularTable: true,
		},
	})

	fmt.Println(db)
	fmt.Println(err)

	// 连接池
	sqlDB, err := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	// 结构体
	type List struct {
		gorm.Model        // 自动添加ID、CreatedAt、UpdatedAt、DeletedAt
		Name       string `grom:"type:varchar(20); not null" json:"name" binding:"required"`
		State      string `grom:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone      string `grom:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email      string `grom:"type:varchar(40); not null" json:"email" binding:"required"`
		Address    string `grom:"type:varchar(200); not null" json:"address" binding:"required"`
	}

	/* 注意点：
	1. 结构体里面的变量（Name）必须是首字符大写
	gorm 指定类型
	json 表示json接受时的名称
	binding required 表示必须传入
	*/

	// 1.没有主键 答：给结构体添加gorm.Model
	// 2.表名变成复数的问题 答：在connection的Config中添加NamingStrategy

	// 迁移：用于使用代码建表
	db.AutoMigrate(&List{})

	// 启动并监听一个HTTP请求
	PORT := "3001"
	r := gin.Default()

	// 测试
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "请求成功",
	//	})
	//})

	// 增加
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		// 接收输入数据进行绑定
		err := c.ShouldBindJSON(&data)

		// 判断绑定是否有错误
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			// 数据库操作
			db.Create(&data) // 创建一条数据

			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})

	// 删除
	// 1. 先找到对应的id所对应的条目
	// 2. 判断id是否存在
	// 3. 返回id没有找到
	// restful编码规范
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List

		// url中:id的形式用Param接收
		id := c.Param("id")
		// url中id="id"的形式用Query接收
		// c.Query()

		// 判断id是否存在
		db.Where("id = ?", id).Find(&data)

		// 判断id存在的情况，则删除，不存在则报错
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id不存在，删除失败",
				"code": 400,
			})
		} else {
			// 操作数据库删除
			// Gorm的删除是软删除，只会记录删除的时间，并不会真的把数据删除了
			db.Where("id = ?", id).Delete(&data)

			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 修改
	r.PUT("/user/update/:id", func(c *gin.Context) {
		// 1. 找到对应的id所对应的条目
		// 2. 判断id是否存在
		// 3. 修改对应条目
		// 4. 返回id没有找到
		var data List

		// url中:id的形式用Param接收
		id := c.Param("id")
		// url中id="id"的形式用Query接收
		// c.Query()

		// 判断id是否存在
		db.Select("id").Where("id = ?", id).Find(&data)
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "修改失败",
					"code": 400,
				})
			} else {
				// 修改数据库内容
				db.Where("id = ?", id).Updates(&data)

				c.JSON(200, gin.H{
					"msg":  "修改成功",
					"code": 200,
				})
			}
		}
	})

	// 条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {
		// 获取路径参数
		name := c.Param("name")

		var dataList []List

		// 查询数据库
		db.Where("name = ?", name).Find(&dataList)
		// 判断是否查询到数据
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"data": gin.H{},
				"code": 400,
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"data": dataList,
				"code": 200,
			})
		}
	})

	// 全部查询/分页查询
	r.GET("/user/list", func(c *gin.Context) {
		var dataList []List

		// 1. 查询全部数据
		// 2. 查询分页数据
		var total int64
		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}

		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		fmt.Println(pageSize)
		fmt.Println(pageNum)

		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)

		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}

	})

	r.Run(":" + PORT)

}
