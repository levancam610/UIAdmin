package controllers

import (
	"../DTO"
	"../database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gin-contrib/sessions"


)

type Post struct {
	Id int `json:"id"`
	FirstName string `json:"title"`
	LastName string
	Email string
	UserName string
	PassWord string
}

type Clothes struct {
	Id int
	Name string
	CategoryId int
	Gender string
	Amount int
	Price int
}
type DataClothes struct {
	TotalPage int
	TotalCount int
	Data []DTO.QuanDTO
}
func FindById(c * gin.Context){
/*
	db := database.DBConn()
	rows, err := db.Query("SELECT * FROM users WHERE id = " + c.Param("id"))
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
	}

	 := Post{}

	for rows.Next(){
		var id, categoryId int
		var amount, price float32
		var name, gender string

		err = rows.Scan(&id, &name, &categoryId, &gender, &amount, &price)
		if err != nil {
			panic(err.Error())
		}

		post.Id = id
		post.FirstName = firstName
		post.LastName = lastName
		post.Email = email
		post.UserName = userName
		post.PassWord = passWord
	}

	c.JSON(200, post)
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong*/
}

func GetList(c * gin.Context){
	page:= c.Param("page")
	tmp,err := strconv.Atoi(page)
	if(err!=nil){
		panic("loi page ko convert dc integer")
	}
	from := strconv.Itoa(tmp*12)
	/* ORDER BY created desc limit */
	db := database.DBConn()
	rows, err := db.Query("SELECT id, name, categoryId, gender, amount, price FROM clothes ORDER BY created desc limit "+from+",12")
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
	}
	post := DTO.QuanDTO{}
	list := [] DTO.QuanDTO{}
	for rows.Next(){
		var id, categoryId int
		var amount, price int
		var name, gender string

		err = rows.Scan(&id, &name, &categoryId, &gender, &amount, &price)
		if err != nil {
			panic(err.Error())
		}

		post.Id = id
		post.Name = name
		post.Amount = amount
		post.CategoryId = categoryId
		post.Gender = gender
		post.Price = price
		list = append(list,post);

	}
	c.JSON(200, list)
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong
}


func GetList2(c * gin.Context){
	currentPage:= c.Query("currentPage")
	pageSize:= c.Query("pageSize")
	orderBy:= c.Query("orderBy")
	search:= c.Query("search")
	categoryId := c.Query("categoryId")

	currentPage2,err1 := strconv.Atoi(currentPage)
	currentPage2 -= 1;
	pageSize2,err2 := strconv.Atoi(pageSize)
	if(err1!=nil || err2!=nil){
		panic("loi page ko convert dc integer")
	}
	from := strconv.Itoa(currentPage2*pageSize2)
	/* ORDER BY created desc limit */
	db := database.DBConn()
	sql := "SELECT id, name, categoryId, gender, amount, price, image FROM clothes where categoryId="+categoryId+" and name like '%"+search+"%' ORDER BY "+orderBy+" desc limit "+from+","+pageSize;
	fmt.Sprintf(sql)
	rows, err := db.Query(sql)
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
	}

	data := DataClothes{}
	post := DTO.QuanDTO{}
	//list := [] DTO.QuanDTO{}
	for rows.Next(){
		var id, categoryId int
		var amount, price int
		var name, gender, image string
		err = rows.Scan(&id, &name, &categoryId, &gender, &amount, &price, &image)
		if err != nil {
			panic(err.Error())
		}
		post.Id = id
		post.Name = name
		post.Amount = amount
		post.CategoryId = categoryId
		post.Gender = gender
		post.Price = price
		post.Image = image
		data.Data = append(data.Data, post)
	}
	rows2, err := db.Query("Select count(*) FROM clothes where categoryId="+categoryId)
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
		panic(err)
	}
	var rs int
	for rows2.Next() {
		err = rows2.Scan(&rs)
	}
	data.TotalCount = rs
	if rs%pageSize2==0 {
		rs = rs/pageSize2
	} else{
		rs = rs/pageSize2 +1
	}
	data.TotalPage = rs
	c.JSON(200, data)
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong
}
func CreateClothes(c * gin.Context){
	db := database.DBConn()
	var json DTO.QuanDTO
	toDay := time.Now().Format("02-01-2006")
	if err := c.ShouldBindJSON(&json); err == nil {
		stm, err := db.Prepare("INSERT INTO clothes SET name=?, categoryId=?, gender=?, amount=?, price=?,created=?")
		if err != nil {
			panic(err.Error())
		}
		stm.Exec(json.Name,json.CategoryId, json.Gender, json.Amount, json.Price, toDay)
		c.JSON(200, gin.H{
			"messages": "inserted",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	defer db.Close()
}

/* find category by id */
func GetAllCategory(c * gin.Context){
	db := database.DBConn()
	rows, err := db.Query("SELECT * FROM category")
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
	}
	post := DTO.CategoryDTO{}
	list := [] DTO.CategoryDTO{}
	for rows.Next(){
		var id, types int
		var name string
		err = rows.Scan(&id, &name, &types)
		if err != nil {
			panic(err.Error())
		}
		post.Id = id
		post.Name = name
		post.Type = types
		list = append(list,post);
	}
	c.JSON(200, list)
	defer db.Close()
}
/*----------------------------- */

/* Delete clothes*/
func DeleteClothes(c * gin.Context){
	db := database.DBConn()
	id:= c.Param("id")
	_, err := db.Query("Delete FROM clothes WHERE id = " + id)
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
		panic("error delte clothes")
	}
	c.JSON(200, gin.H{
		"messages": "deleted",
	})
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong*/
}
/*---------------------------*/
func UploadImage(c * gin.Context){
	db := database.DBConn()
	var json DTO.ImageDTO
	if err := c.ShouldBindJSON(&json); err == nil {
		stm, err := db.Prepare("INSERT INTO images SET link=?, clothesId=?")
		if err != nil {
			panic(err.Error())
		}
		stm.Exec(json.Link,json.ClothesId)
		c.JSON(200, gin.H{
			"messages": "inserted",
		})

	} else {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	defer db.Close()
}

/* Load image from clothes*/
func GetImageByClothesId(c * gin.Context){
	db := database.DBConn()
	rows, err := db.Query("SELECT * FROM images where clothesId = "+c.Param("id"))
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
	}
	post := DTO.ImageDTO{}
	list := [] DTO.ImageDTO{}
	for rows.Next(){
		var id, clothesId int
		var link string
		err = rows.Scan(&id, &link, &clothesId)
		if err != nil {
			panic(err.Error())
		}
		post.Id = id
		post.ClothesId = clothesId
		post.Link = link
		list = append(list,post)
	}
	c.JSON(200, list)
	defer db.Close()
}
/*--------------------------*/

/*--------delete image from clothes list--------------------*/
func DeleteImage(c * gin.Context){
	db := database.DBConn()
	id:= c.Param("id")
	_, err := db.Query("Delete FROM images WHERE id = " + id)
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
		panic("error delte clothes")
	}
	c.JSON(200, gin.H{
		"messages": "deleted",
	})
	defer db.Close()
}
/*----------------------------*/
/*-------- count page for list clothes --------*/

/* Delete clothes*/
func CountPageClothes(c * gin.Context){
	db := database.DBConn()
	rows, err := db.Query("Select count(*) FROM clothes")
	if err != nil{
		c.JSON(500, gin.H{
			"messages" : "Story not found",
		});
		panic(err)
	}
	var rs int
	for rows.Next() {
		err = rows.Scan(&rs)
	}
	if rs%12==0 {
		rs = rs/12
	} else{
		rs = rs/12 +1
	}
	c.JSON(200, rs)
	defer db.Close() // Hoãn lại việc close database connect cho đến khi hàm Read() thực hiệc xong*/
}
/*------------------------------ */
/* ----------- Login ---------------*/
func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	var user, pass string
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}
	db := database.DBConn()
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&user, &pass)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username no valid"})
		return
	}
	if strings.Trim(pass, " ") == strings.Trim(password, " "){
		session.Set("user", username) //In real world usage you'd set this to the users ID
		session.Save()
		c.JSON(http.StatusOK, gin.H{"status": "OK"}) //Successfully authenticated user
		return
	}else{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password wrong"})
	}
}

func GetSession(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user")!=nil{
		c.JSON(http.StatusOK, session.Get("user"))
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user")!=nil{
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"status":"OK"})
	}
}

/* ----------------------------------*/