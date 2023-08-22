package accrual

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ANiWarlock/gophermart/cmd/gophermart/config"
	"net/http"
	"time"
)

type Accrual struct {
	Order   string
	Status  string
	Accrual float64
}

var cfg config.AppConfig

func Init(conf *config.AppConfig) {
	cfg = *conf
}

func Get(order string) (*Accrual, error) {
	var accrual Accrual
	var buf bytes.Buffer
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	path := fmt.Sprintf("%s/api/orders/%s", cfg.AccrualSystemAddress, order)
	resp, err := client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get accrual: %w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot process body: %v", err)
	}
	defer resp.Body.Close()

	if err = json.Unmarshal(buf.Bytes(), &accrual); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}

	return &accrual, nil
}
