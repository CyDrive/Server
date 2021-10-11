package models

import "github.com/CyDrive/consts"

type Response struct {
	StatusCode consts.StatusCode `json:"status_code"`
	Message    string            `json:"message,omitempty"`
	Data       string            `json:"data,omitempty"` // json format
}
