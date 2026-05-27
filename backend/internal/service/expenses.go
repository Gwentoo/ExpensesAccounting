package service

import (
	"backend/internal/config"
	"backend/internal/jwt"
	myParser "backend/internal/parser"
	"backend/internal/storage/postgres"
	expensesv1 "backend/proto/gen/expenses"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SupersetAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type GuestTokenRequest struct {
	User      SupersetUser       `json:"user"`
	Resources []SupersetResource `json:"resources"`
	RLS       []SupersetRLS      `json:"rls"`
}

type SupersetUser struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SupersetResource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type SupersetRLS struct {
	Clause string `json:"clause"`
}
type ExpensesService struct {
	expensesv1.UnimplementedExpensesServer
	db *postgres.DB
}

func NewExpensesService(db *postgres.DB) *ExpensesService {
	return &ExpensesService{db: db}
}

func (s *ExpensesService) UploadStatement(ctx context.Context, req *expensesv1.UploadStatementRequest) (*expensesv1.UploadStatementResponse, error) {
	userID, err := jwt.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unauthorized: %w", err)
	}

	firstLineReader := csv.NewReader(bytes.NewReader(req.Content))
	firstLineReader.Comma = ';'
	header, err := firstLineReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %w", err)
	}

	cleanHead := strings.TrimPrefix(header[0], "\ufeff")

	var parser myParser.StatementParser
	if cleanHead == "operationDate" {
		parser = &myParser.AlfaParser{}
	} else if cleanHead == "operation_date" {
		parser = &myParser.TBankParser{}
	} else {
		return nil, fmt.Errorf("unknown statement format: %s", cleanHead)
	}

	reader := csv.NewReader(bytes.NewReader(req.Content))
	reader.LazyQuotes = true
	reader.Comma = parser.GetDelimiter()
	reader.FieldsPerRecord = -1
	_, _ = reader.Read()

	var processedRows int32
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("skip invalid line: %v", err)
			continue
		}

		tx, err := parser.Parse(record)
		if err != nil {
			log.Printf("failed to parse row: %v", err)
			continue
		}

		tx.UserID = userID

		err = s.db.SaveTransaction(ctx, *tx)
		if err != nil {
			return nil, fmt.Errorf("database error at row %d: %w", processedRows, err)
		}

		processedRows++
	}

	return &expensesv1.UploadStatementResponse{
		Success:       true,
		ProcessedRows: processedRows,
	}, nil
}

func (s *ExpensesService) GetSupersetToken(ctx context.Context, req *expensesv1.GetSupersetTokenRequest) (*expensesv1.GetSupersetTokenResponse, error) {
	cfg := config.LoadConfig()

	userID, err := jwt.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "ОШИБКА JWT: %v", err)
	}

	authBody, _ := json.Marshal(map[string]string{
		"username": cfg.SuperSet.Username,
		"password": cfg.SuperSet.Password,
		"provider": "db",
	})

	resp, err := http.Post("http://superset:"+strconv.Itoa(cfg.SuperSet.Port)+"/api/v1/security/login", "application/json", bytes.NewBuffer(authBody))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ОШИБКА LOGIN СЕТЬ: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, status.Errorf(codes.Internal, "ОШИБКА LOGIN СТАТУС: %d", resp.StatusCode)
	}

	var authResp SupersetAuthResponse
	json.NewDecoder(resp.Body).Decode(&authResp)

	guestReq := GuestTokenRequest{
		User: SupersetUser{
			Username:  "guest_user",
			FirstName: "Guest",
			LastName:  "User",
		},
		Resources: []SupersetResource{
			{Type: "dashboard", ID: req.DashboardId},
		},
		RLS: []SupersetRLS{
			{Clause: fmt.Sprintf("operation_date >= '%s'", req.StartDate)},
			{Clause: fmt.Sprintf("operation_date <= '%s'", req.EndDate)},
			{Clause: fmt.Sprintf("user_id = %d", userID)},
		},
	}

	guestBody, _ := json.Marshal(guestReq)

	client := &http.Client{}
	httpReq, _ := http.NewRequest("POST", "http://superset:8088/api/v1/security/guest_token/", bytes.NewBuffer(guestBody))
	httpReq.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Host = "localhost:8088"

	guestResp, err := client.Do(httpReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ОШИБКА GUEST СЕТЬ: %v", err)
	}
	defer guestResp.Body.Close()
	if guestResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(guestResp.Body)
		return nil, status.Errorf(codes.Internal, "SUPERSET ОТКЛОНИЛ ЗАПРОС (Код %d). Ответ: %s. ID в запросе был: '%s'",
			guestResp.StatusCode, string(bodyBytes), req.DashboardId)
	}

	var tokenData struct {
		Token string `json:"token"`
	}
	json.NewDecoder(guestResp.Body).Decode(&tokenData)

	return &expensesv1.GetSupersetTokenResponse{
		Token:       tokenData.Token,
		DashboardId: req.DashboardId,
	}, nil
}
