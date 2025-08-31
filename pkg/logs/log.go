package logs

import (
	"database/sql"
	"fmt"
	"time"
)

type LogType string

const (
	AccessLog      LogType = "akses_log"
	TransactionLog LogType = "transaction_log"
	ActivityLog    LogType = "activity_log"
	PasienLog      LogType = "pasien_log"
)

type LogEntry struct {
	Tanggal string
	Waktu   string
	User    string
	Pesan   string
	Status  string
}

var db *sql.DB

func InitializeLogger(database *sql.DB) {
	db = database
}

func CreateLog(logType LogType, userID, message, status string) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	now := time.Now()
	logEntry := LogEntry{
		Tanggal: now.Format("2006-01-02"),
		Waktu:   now.Format("15:04:05"),
		User:    userID,
		Pesan:   message,
		Status:  status,
	}

	var query string
	switch logType {
	case AccessLog:
		query = "INSERT INTO akses_log (Tanggal, Waktu, User, Message, Status) VALUES (?, ?, ?, ?, ?)"
	case TransactionLog:
		query = "INSERT INTO transaksi_log (Tanggal, Waktu, User, Message, Status) VALUES (?, ?, ?, ?, ?)"
	case ActivityLog:
		query = "INSERT INTO aktivitas_user_log (Tanggal, Waktu, User, Message, Status) VALUES (?, ?, ?, ?, ?)"
	case PasienLog:
		query = "INSERT INTO pasien_log (Tanggal, Waktu, User, Message, Status) VALUES (?, ?, ?, ?, ?)"
	default:
		return fmt.Errorf("unknown log type: %s", logType)
	}

	_, err := db.Exec(query, logEntry.Tanggal, logEntry.Waktu, logEntry.User, logEntry.Pesan, logEntry.Status)
	if err != nil {
		return fmt.Errorf("failed to create %s log: %w", logType, err)
	}

	return nil
}

func LogAccess(userID, message, status string) error {
	return CreateLog(AccessLog, userID, message, status)
}

func LogTransaction(userID, message, status string) error {
	return CreateLog(TransactionLog, userID, message, status)
}

func LogActivity(userID, message, status string) error {
	return CreateLog(ActivityLog, userID, message, status)
}

func LogPasien(userID, message, status string) error {
	return CreateLog(PasienLog, userID, message, status)
}
