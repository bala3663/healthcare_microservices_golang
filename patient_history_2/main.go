// Patient_History_2/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Patient1 struct {
	ID          int
	PatientName string
	PatientType string
	POC         string
	Plan        string
	Treatment   string
}

func main() {
	db, err := sql.Open("mysql", "root:india@123@tcp(localhost:3306)/Patient_information")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("C:/Users/balaji/Desktop/Microservice_project/templates/index1.html")

	router.GET("/", func(c *gin.Context) {
		patients, err := getPatients(db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index1.html", gin.H{
				"Error": "Failed to retrieve patients.",
			})
			return
		}

		c.HTML(http.StatusOK, "index1.html", gin.H{
			"Patients": patients,
		})
	})

	router.POST("/add", func(c *gin.Context) {

		patientName := c.PostForm("patient_name")
		patientType := c.PostForm("patient_type")
		POC := c.PostForm("poc")
		plan := c.PostForm("plan")
		treatment := c.PostForm("treatment")

		_, err := db.Exec("INSERT INTO patients_history (patient_name, patient_type, poc, plan, treatment) VALUES (?, ?, ?, ?, ?)", patientName, patientType, POC, plan, treatment)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index1.html", gin.H{
				"Error": "Failed to add patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	router.POST("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		patientName := c.PostForm("patient_name")
		patientType := c.PostForm("patient_type")
		POC := c.PostForm("poc")
		plan := c.PostForm("plan")
		treatment := c.PostForm("treatment")

		_, err := db.Exec("UPDATE patients_history SET patient_name = ?, patient_type = ?, poc = ?, plan = ?, treatment = ? WHERE id = ?", patientName, patientType, POC, plan, treatment, id)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index1.html", gin.H{
				"Error": "Failed to update patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	router.GET("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM patients_history WHERE id = ?", id)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "index1.html", gin.H{
				"Error": "Failed to delete patient.",
			})
			return
		}

		c.Redirect(http.StatusFound, "/")
	})

	fmt.Println("Server started at http://localhost:8081")
	router.Run(":8081")

}

func getPatients(db *sql.DB) ([]Patient1, error) {
	rows, err := db.Query("SELECT id, patient_name, patient_type, poc, plan, treatment FROM patients_history")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient1

	for rows.Next() {
		var patient Patient1
		err := rows.Scan(&patient.ID, &patient.PatientName, &patient.PatientType, &patient.POC, &patient.Plan, &patient.Treatment)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}
