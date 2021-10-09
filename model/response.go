package model

import "github.com/CyDrive/consts"

type Response struct {
	StatusCode consts.StatusCode `json:"status_code"`
	Message    string            `json:"message,omitempty"`
	Data       interface{}       `json:"data,omitempty"`
}