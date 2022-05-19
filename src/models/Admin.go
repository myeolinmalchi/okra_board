package models

type Admin struct {
    ID          string      `json:"id"`
    Password    string      `json:"pw"`
    Name        string      `json:"name,omitempty"`
    Email       string      `json:"email,omitempty"`
    Phone       string      `json:"phone,omitemtpy"`
}

type AdminValidtionResult struct {
    ID          *string      `json:"id,omitempty"`
    Password    *string      `json:"pw,omitempty"`
    Name        *string      `json:"name,omitempty"`
    Email       *string      `json:"email,omitempty"`
    Phone       *string      `json:"phone,omitempty"`
}