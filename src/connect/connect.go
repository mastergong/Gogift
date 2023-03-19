package connect

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	member "giftcard/json"
	"log"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

var serverHost = ""

var db *sql.DB

func getConnectionString(uid, password, hostName string) string {
	query := getCommonQueryString("master")
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(uid, password),
		Host:     hostName,
		RawQuery: query.Encode(),
	}
	return u.String()
}

func getCommonQueryString(dbName string) url.Values {
	query := url.Values{}
	query.Add("connection timeout", strconv.Itoa(30))
	query.Add("Encrypt", "disable")
	query.Add("TrustServerCertificate", "true")
	//query.Add("log", "63")
	return query
}

func Connect() error {

	var serverUser string
	fmt.Print("Enter User_server: ")
	fmt.Scanln(&serverUser)

	var serverPwd string
	fmt.Print("Enter Pwd_server: ")
	fmt.Print("\033") // Hide inpu
	fmt.Scanln(&serverPwd)
	fmt.Println("\033") // Show input

	connString := getConnectionString(serverUser, serverPwd, serverHost) //fmt.Sprintf("server=%s;user id=%s;password=%s;", serverHost, serverUser, serverPwd)
	//fmt.Println(connString)
	var err error
	db, err = sql.Open("sqlserver", connString)

	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	cmd = exec.Command("cmd", "/c", "Color 6")
	cmd.Stdout = os.Stdout
	cmd.Run()

	log.Printf("Connected!\n")

	return err

}

func Disconnect() {

	if db != nil {
		defer db.Close()
	}
}

func GetUserMembName(membname string) string {

	tsql := " "
	return getUser(tsql)
}

func GetUserMembId(membno string) string {

	tsql := ""
	return getUser(tsql)
}

func GetUserMembIdLoan(membno string) string {

	tsql := " "
	return getUser(tsql)
}

func GetUserCardid(CardId string) string {
	tsql := fmt.Sprintf(" '%s'  ", CardId)
	return getUser(tsql)
}

func getUserloan(queryWhere string) string {

	ctx := context.Background()
	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return ""
	}

	tsql := " "
	tsql += " "
	tsql += queryWhere

	rows, err := getJSON(tsql)

	if err != nil {
		return ""
	}
	return rows
}

func getUser(queryWhere string) string {

	ctx := context.Background()
	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return ""
	}

	tsql := ""
	tsql += ""
	tsql += " "
	tsql += queryWhere

	rows, err := getJSON(tsql)

	if err != nil {
		return ""
	}

	//*******************
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(rows)

	return rows
}

func UpdateUser(member []member.MemberArray) (err error) {

	fmt.Println("-------------------------------")
	//fmt.Println(member)

	for _, elem := range member {

		v := reflect.ValueOf(elem)
		typeOfS := v.Type()
		STMN_MEMB_NO2 := "xx"
		STMN_MEMB_NO := "xx"
		STMN_RECV_STS := "xx"
		for i := 0; i < v.NumField(); i++ {
			if typeOfS.Field(i).Name == "" {
				STMN_MEMB_NO2 = fmt.Sprintf("%v", v.Field(i).Interface())
			}

			if typeOfS.Field(i).Name == "" {
				STMN_MEMB_NO = fmt.Sprintf("%v", v.Field(i).Interface())
			}

			if typeOfS.Field(i).Name == "" {
				STMN_RECV_STS = fmt.Sprintf("%v", v.Field(i).Interface())
			}
		}

		fmt.Println(v)

		fmt.Println(" " + STMN_RECV_STS)

		if STMN_RECV_STS == "0" || STMN_RECV_STS == "" {

			STMN_RECV_STS = "1"
			if STMN_MEMB_NO2 != "" {
				STMN_RECV_STS = "2"
			}
			var sqlUpdate = ""
			fmt.Println(sqlUpdate)
			if _, err = db.Exec(sqlUpdate); err != nil {
				panic(err)
			} else {
				fmt.Println("Update Success!")
			}

			//*******************
			f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()

			log.SetOutput(f)
			log.Println(sqlUpdate)
		}
	}
	return err
}

func getJSON(sqlString string) (string, error) {
	rows, err := db.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}
