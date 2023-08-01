// Final_Healthcare_Project/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Patient struct {
	ID         int
	FirstName  string
	LastName   string
	Age        json.Number
	Address    string
	Disease    string
	Modifier   string
	Payer      string
	BillNo     json.Number
	AmountPaid json.Number
}

func main() {
	db, err := sql.Open("mysql", "root:india@123@tcp(localhost:3306)/Patient_information")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("C:/Users/balaji/Desktop/Microservice_project/templates/index.html")

	router.GET("/", func(c *gin.Context) { //context of the current HTTP request.
		patients, err := getPatients(db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"Error": "Failed to retrieve patients.",
			})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Patients": patients,
		})
	})

	router.POST("/add", func(c *gin.Context) {
		firstName := c.PostForm("first_name") //"first_name" parameter from a form submitted via an HTTP POST
		lastName := c.PostForm("last_name")
		age := c.PostForm("age")
		address := c.PostForm("address")
		disease := c.PostForm("disease")
		modifier := c.PostForm("modifier")
		payer := c.PostForm("payer")
		billNo := c.PostForm("bill_no")
		amountPaid := c.PostForm("amount_paid")

		_, err := db.Exec("INSERT INTO patients2 (first_name, last_name, age, address, disease, modifier, payer, bill_no, amount_paid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", firstName, lastName, age, address, disease, modifier, payer, billNo, amountPaid)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"Error": "Failed to add patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	router.POST("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		firstName := c.PostForm("first_name")
		lastName := c.PostForm("last_name")
		age := c.PostForm("age")
		address := c.PostForm("address")
		disease := c.PostForm("disease")
		modifier := c.PostForm("modifier")
		payer := c.PostForm("payer")
		billNo := c.PostForm("bill_no")
		amountPaid := c.PostForm("amount_paid")

		_, err := db.Exec("UPDATE patients2 SET first_name = ?, last_name = ?, age = ?, address = ?, disease = ?, modifier = ?, payer = ?, bill_no = ?, amount_paid = ? WHERE id = ?", firstName, lastName, age, address, disease, modifier, payer, billNo, amountPaid, id)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"Error": "Failed to update patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	router.GET("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM patients2 WHERE id = ?", id)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{
				"Error": "Failed to delete patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	fmt.Println("Server started at http://localhost:8080")
	router.Run(":8080")

}

func getPatients(db *sql.DB) ([]Patient, error) {
	rows, err := db.Query("SELECT id, first_name, last_name, age, address, disease, modifier, payer, bill_no, amount_paid FROM patients2")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient

	for rows.Next() {
		var patient Patient
		err := rows.Scan(&patient.ID, &patient.FirstName, &patient.LastName, &patient.Age, &patient.Address, &patient.Disease, &patient.Modifier, &patient.Payer, &patient.BillNo, &patient.AmountPaid)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}
