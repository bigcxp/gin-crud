package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//定义User结构体
type User struct {
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

//查询全部信息
func getAll() (users []User, err error) {
	//1.操作数据库
	db, _ := sql.Open("mysql", "root:pxichen@tcp(localhost:3306)/gin?charset=utf8")
	//处理错误
	if err != nil {
		log.Fatal(err.Error())
	}
	//关闭数据库连接
	defer db.Close()

	//2.查询
	rows, err := db.Query("SELECT id,username,password FROM user_tb")
	//处理错误
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		var user User
		//遍历表中所有行的信息
		rows.Scan(&user.Id, &user.Username, &user.Password)
		//将user添加到users切片中
		users = append(users, user)
	}
	//关闭连接
	defer rows.Close()
	return
}

//添加数据
func addUser(user User) (Id int, err error) {
	//1.操作数据库
	db, err := sql.Open("mysql", "root:pxichen@tcp(localhost:3306)/gin?charset=utf8")
	//处理错误
	if err != nil {
		log.Fatal(err.Error())
	}
	//关闭数据库连接
	defer db.Close()
	//预sql处理
	stmt, err := db.Prepare("INSERT INTO user_tb(username,password) VALUES(?,?)")
	if err != nil {
		return
	}
	//执行插入操作
	fmt.Println(user)
	rs, err := stmt.Exec(user.Username, user.Password)
	if err != nil {
		return
	}
	//返回插入的id
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	//将id类型转换
	Id = int(id)
	defer stmt.Close()
	return
}

//修改数据
func update(user User) (rowsAffected int64, err error) {

	//1.操作数据库
	db, err := sql.Open("mysql", "root:pxichen@tcp(localhost:3306)/gin?charset=utf8")
	//错误检查
	if err != nil {
		log.Fatal(err.Error())
	}
	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("UPDATE  user_tb SET username=?, password=? WHERE id=?")
	if err != nil {
		return
	}
	//执行修改操作
	rs, err := stmt.Exec(user.Username, user.Password, user.Id)
	if err != nil {
		return
	}
	//返回插入的id
	rowsAffected, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	return
}

//通过id删除
func del(id int) (rows int, err error) {
	//1.操作数据库
	db, err := sql.Open("mysql", "root:pxichen@tcp(localhost:3306)/gin?charset=utf8")
	//错误检查
	if err != nil {
		log.Fatal(err.Error())
	}
	//推迟数据库连接的关闭
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM user_tb WHERE id=?")
	if err != nil {
		log.Fatalln(err)
	}

	rs, err := stmt.Exec(id)
	if err != nil {
		log.Fatalln(err)
	}
	//删除的行数
	row, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	rows = int(row)
	return
}

func main() {
	//创建路由
	router := gin.Default()
	//查询所有
	router.GET("/user", func(c *gin.Context) {
		//调用查询方法
		users, err := getAll()
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"result":  users,
			"account": len(users),
		})
	})

	//添加数据
	router.POST("/add", func(c *gin.Context) {
		var user User
		//接收并绑定参数
		err := c.Bind(&user)
		if err != nil {
			log.Fatal(err)
		}
		//调用添加方法
		Id, err := addUser(user)
		fmt.Print("id=", Id)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%s插入成功", user.Username),
		})
	})

	//修改数据
	router.PUT("/update", func(c *gin.Context) {
		var u User
		err := c.Bind(&u)
		if err != nil {
			log.Fatal(err)
		}
		num, err := update(u)
		fmt.Print("num=", num)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("修改id: %d 成功", u.Id),
		})
	})

	//利用DELETE请求方法通过id删除
	router.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		Id, err := strconv.Atoi(id)

		if err != nil {
			log.Fatalln(err)
		}
		rows, err := del(Id)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("delete rows ", rows)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})

	router.Run(":8080")

}
