package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gomail.v2"
)

var db *sql.DB

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

func main() {
	var err error

	// Update this with your actual DB config
	// dsn := "root:root@tcp(127.0.0.1:8889)/Ecomm"
	dsn := "root:shyam@tcp(13.200.235.187:3306)/Ecomm"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error opening DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("❌ DB connection failed:", err)
	}
	fmt.Println("✅ Connected to MySQL successfully")

	http.HandleFunc("/users", usersHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		insertUser(w, r)
	case http.MethodGet:
		getUsers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO Users(name, email, mobile) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare insert query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.Name, user.Email, user.Mobile)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	user.ID = int(id)

	// Send email in background
	go sendWelcomeEmail(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, mobile FROM Users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Mobile); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func sendWelcomeEmail(user User) {
	fmt.Print(user)
	m := gomail.NewMessage()
	m.SetHeader("From", "varun.singhal78@gmail.com")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Welcome to Goapp!")

	// Improved HTML content with better formatting and style
	htmlBody := fmt.Sprintf(`
		<html>
			<head>
				<style>
					body {
						font-family: Arial, sans-serif;
						color: #333;
						line-height: 1.6;
					}
					.container {
						max-width: 600px;
						margin: 0 auto;
						padding: 20px;
						background-color: #f9f9f9;
						border-radius: 8px;
					}
					.header {
						text-align: center;
						padding-bottom: 20px;
					}
					.header h2 {
						color: #4CAF50;
					}
					.content {
						font-size: 16px;
						padding: 10px 0;
					}
					.footer {
						text-align: center;
						font-size: 12px;
						color: #777;
						padding-top: 20px;
					}
					.button {
						display: inline-block;
						background-color: #4CAF50;
						color: #fff;
						padding: 10px 20px;
						text-decoration: none;
						border-radius: 5px;
						margin-top: 20px;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<h2>Welcome to Goapp, %s!</h2>
					</div>
					<div class="content">
						<p>Thanks for registering on <strong>Goapp</strong>!</p>
						<p><strong>Your mobile number:</strong> %s</p>
						<a href="https://www.appslure.com" class="button">Visit Goapp</a>
					</div>
					<div class="footer">
						<p>Regards,<br><i>Ecomm Team</i></p>
						<p>If you have any questions, feel free to reach out to our support team.</p>
					</div>
				</div>
			</body>
		</html>`,
		user.Name, user.Mobile)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "varun.singhal78@gmail.com", "hrnk zfpf alrl uscj")

	if err := d.DialAndSend(m); err != nil {
		log.Println("❌ Email send failed:", err)
		return
	}

	log.Println("✅ HTML welcome email sent to", user.Email)
}
