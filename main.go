package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)
type Person struct {
	Id int
	First_Name string
	Last_Name string
}
var db *sql.DB
var err error

func main()  {
	srv := Server()
	RegistrodeRutas(srv)
	srv.Run(":8000")
}


func Server() *gin.Engine {
srv := gin.Default()
return srv
}


func RegistrodeRutas(srv *gin.Engine) {
	srv.GET("/personas", personas)
	srv.GET("/persona/:id", persona)
	srv.GET("/person", person)
	srv.POST("/person", in_persona)
	srv.PUT("/person", up_persona)
	srv.PUT("/person/:id", upd_persona)
	srv.DELETE("/person", dl_persona)
	srv.DELETE("/person/:id", dlt_persona)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func personas(c *gin.Context)  {
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")
	checkErr(err)
	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)

	var (
		persona Person
		personas     []Person
	)
	rows,err :=db.Query("select id, first_name, last_name from person;")
	checkErr(err)
	if rows!=nil{
		for rows.Next(){
			err:=rows.Scan(&persona.Id,&persona.First_Name,&persona.Last_Name)
			personas=append(personas, persona)
			checkErr(err)
		}
		defer rows.Close()
	}
	c.JSON(http.StatusOK,gin.H{"result":personas,"cantidad":len(personas)})
}
func persona(c *gin.Context)  {
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)
	var (
		persona Person
		result gin.H
	)
	id:=c.Param("id")
	row:=db.QueryRow("select id, first_name, last_name from person where id = ?;",id)
	err=row.Scan(&persona.Id,&persona.First_Name,&persona.Last_Name)
	if err != nil{
		result=gin.H{"result":nil,"cantidad":0}
	}else{
		result=gin.H{"result":persona,"cantidad":1}
	}
	c.JSON(http.StatusOK, result)
}
func person(c *gin.Context)  {
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)
	var (
		persona Person
		result gin.H
	)
	//id:=c.Param("id")
	id := c.Query("id")
	row:=db.QueryRow("select id, first_name, last_name from person where id = ?;",id)
	err=row.Scan(&persona.Id,&persona.First_Name,&persona.Last_Name)
	if err != nil{
		result=gin.H{"result":nil,"cantidad":0}
	}else{
		result=gin.H{"result":persona,"cantidad":1}
	}
	c.JSON(http.StatusOK, result)
}
func in_persona(c *gin.Context){
	c.Header("Content-Type", "application/json; charset=utf-8")
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)
	var persona Person
	var first_name string
	var last_name string
	var buffer bytes.Buffer
	if c.ShouldBind(&persona) == nil {
		log.Println(persona.First_Name)
		log.Println(persona.Last_Name)
		first_name =persona.First_Name
		last_name =persona.Last_Name
	}else{
		first_name = c.PostForm("first_name")
		last_name = c.PostForm("last_name")
	}
	stmt, err := db.Prepare("insert into person (first_name, last_name) values(?,?);")
	checkErr(err)
	if (first_name!="" && last_name!="" ){
		_, err = stmt.Exec(first_name, last_name)
		checkErr(err)
		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)

		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"message": fmt.Sprintf("Problemas en su OBJ enviado"),
		})
	}

	defer stmt.Close()


}
func up_persona(c *gin.Context){
	c.Header("Content-Type", "application/json; charset=utf-8")
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)


	var persona Person
	var first_name string
	var last_name string
	var buffer bytes.Buffer
	//127.0.0.1:8000/person?id=18
	id := c.Query("id")
	//id:=c.Param("id")
	fmt.Println(id)
	if c.ShouldBind(&persona) == nil {
		log.Println(persona.First_Name)
		log.Println(persona.Last_Name)
		first_name =persona.First_Name
		last_name =persona.Last_Name
	}else{
		first_name = c.PostForm("first_name")
		last_name = c.PostForm("last_name")
	}

	stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
	checkErr(err)
	if (first_name!="" && last_name!="" && id!=""){
		_, err = stmt.Exec(first_name,last_name,id)
		if err != nil {
			fmt.Print(err.Error())
		}

		defer stmt.Close()
		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"message": fmt.Sprintf("Problemas al UPDATE %s",id),
		})
	}
	defer stmt.Close()


}
func upd_persona(c *gin.Context){
	c.Header("Content-Type", "application/json; charset=utf-8")
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)
	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)


	var persona Person
	var first_name string
	var last_name string
	var buffer bytes.Buffer
	//id := c.Query("id")
	//127.0.0.1:8000/person/18
	id:=c.Param("id")
	fmt.Println(id)
	if c.ShouldBind(&persona) == nil {
		log.Println(persona.First_Name)
		log.Println(persona.Last_Name)
		first_name =persona.First_Name
		last_name =persona.Last_Name
	}else{
		first_name = c.PostForm("first_name")
		last_name = c.PostForm("last_name")
	}

	stmt, err := db.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
	checkErr(err)
	if (first_name!="" && last_name!="" && id!=""){
		_, err = stmt.Exec(first_name,last_name,id)
		checkErr(err)

		defer stmt.Close()
		// Fastest way to append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"message": fmt.Sprintf("Problemas al UPDATE %s",id),
		})
	}
	defer stmt.Close()


}
func dl_persona(c *gin.Context){
	c.Header("Content-Type", "application/json; charset=utf-8")
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)

	//127.0.0.1:8000/person?id=18
	id := c.Query("id")
	//id:=c.Param("id")
	fmt.Println(id)

	stmt, err := db.Prepare("delete from person where id= ?;")
	checkErr(err)
	if (id!=""){
		_, err = stmt.Exec(id)
		checkErr(err)

		defer stmt.Close()
		// Fastest way to append strings

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully delete to %s", id),
		})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"message": fmt.Sprintf("Problemas al Delete %s",id),
		})
	}
	defer stmt.Close()


}
func dlt_persona(c *gin.Context){
	c.Header("Content-Type", "application/json; charset=utf-8")
	db,err = sql.Open("mysql","root:psw@tcp(ip:3306)/DB")

	checkErr(err)

	defer db.Close()
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	err=db.Ping()
	checkErr(err)

	//127.0.0.1:8000/person/18
	//id := c.Query("id")
	id:=c.Param("id")
	fmt.Println(id)

	stmt, err := db.Prepare("delete from person where id= ?;")
	checkErr(err)
	if (id!=""){
		_, err = stmt.Exec(id)
		checkErr(err)

		defer stmt.Close()
		// Fastest way to append strings

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully delete to %s", id),
		})
	}else{
		c.JSON(http.StatusBadRequest,gin.H{
			"message": fmt.Sprintf("Problemas al Delete %s",id),
		})
	}
	defer stmt.Close()


}
