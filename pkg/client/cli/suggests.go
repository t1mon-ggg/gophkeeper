package cli

import prompt "github.com/c-bata/go-prompt"

var (
	initSuggest = map[string][]prompt.Suggest{
		"": []prompt.Suggest{
			{Text: "cmd", Description: "working area"},
			{Text: "config", Description: "configuration area"},
			{Text: "status", Description: "get current connection state"},
			{Text: "quit", Description: "save changes and exit"},
			{Text: "save", Description: "save changes"},
		},
		">>> ": []prompt.Suggest{
			{Text: "cmd", Description: "working area"},
			{Text: "config", Description: "configuration area"},
			{Text: "status", Description: "get current connection state"},
			{Text: "quit", Description: "save changes and exit"},
			{Text: "save", Description: "save changes"},
		},
		"quit> ": []prompt.Suggest{},
		"cmd> ": []prompt.Suggest{
			{Text: "list", Description: "get list of secrets"},
			{Text: "get", Description: "get secret value"},
			{Text: "insert", Description: "insert new secret to database"},
			{Text: "delete", Description: "remove secret from database"},
			{Text: "quit", Description: "save changes and exit"},
			{Text: "save", Description: "save changes"},
			{Text: "..", Description: "go to up level"},
		},
		"config> ": []prompt.Suggest{
			{Text: "view", Description: "view current config"},
			{Text: "quit", Description: "save changes and exit"},
			{Text: "save", Description: "save changes"},
			{Text: "..", Description: "go to up level"},
		},
		"user> ": []prompt.Suggest{
			{Text: "roster", Description: "hash sum of publickey and confirmation status"},
			{Text: "revoke", Description: "revoke public key"},
			{Text: "confirm", Description: "confirm public key"},
			{Text: "quit", Description: "save changes and exit"},
			{Text: "save", Description: "save changes"},
			{Text: "..", Description: "go to up level"},
		},
	}
	insertSuggest = []prompt.Suggest{
		{Text: "otp", Description: "add new otp secret"},
		{Text: "userpass", Description: "add new user-password secret"},
		{Text: "creditcard", Description: "add new creditcard secret"},
		{Text: "anytext", Description: "add new text secret"},
		{Text: "anybinary", Description: "add new binary secret"},
	}
	otpSuggest = []prompt.Suggest{
		{Text: "method", Description: "set up method"},
		{Text: "issuer", Description: "set up issuer"},
		{Text: "account", Description: "set up account"},
		{Text: "secret", Description: "set up otp secret"},
		{Text: "recoverycodes", Description: "recovery codes"},
		{Text: "description", Description: "define secret description"},
	}
	userpassSuggest = []prompt.Suggest{
		{Text: "username", Description: "set up password"},
		{Text: "password", Description: "set up username"},
		{Text: "description", Description: "define secret description"},
	}
	creditcardSuggest = []prompt.Suggest{
		{Text: "number", Description: "set up credit card number"},
		{Text: "holder", Description: "set up credit hard holder name"},
		{Text: "expire", Description: "set up expire date"},
		{Text: "cvv", Description: "set up cvv code"},
		{Text: "description", Description: "define secret description"},
	}
	anytextSuggest = []prompt.Suggest{
		{Text: "text", Description: "set up text"},
		{Text: "description", Description: "define secret description"},
	}
	anybinarySuggest = []prompt.Suggest{
		{Text: "path", Description: "set up path to binary"},
		{Text: "description", Description: "define secret description"},
	}
	modifySuggest = map[string][]prompt.Suggest{
		"otp":        otpSuggest,
		"userpass":   userpassSuggest,
		"creditcard": creditcardSuggest,
		"anytext":    anytextSuggest,
		"anybinary":  anybinarySuggest,
	}
)
