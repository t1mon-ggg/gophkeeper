package cli

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func (c *CLI) completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	line := in.Text
	switch livePrefixState.livePrefix {
	case "cmd> ":
		if strings.Contains(line, "get") {
			list := c.storage.ListSecrets()
			if len(list) != 0 {
				ss := []prompt.Suggest{}
				for k, v := range list {
					ss = append(ss, prompt.Suggest{Text: k, Description: v})
				}
				return prompt.FilterHasPrefix(ss, w, true)
			}
		}
		if strings.Contains(line, "delete") {
			list := c.storage.ListSecrets()
			if len(list) != 0 {
				ss := []prompt.Suggest{}
				for k, v := range list {
					ss = append(ss, prompt.Suggest{Text: k, Description: v})
				}
				return prompt.FilterHasPrefix(ss, w, true)
			}
		}
		if strings.Contains(line, "insert") {
			args := strings.Split(line, " ")
			if len(args) == 3 {
				return prompt.FilterHasPrefix(insertSuggest, w, true)
			}
			if len(args) > 3 {
				switch args[2] {
				case "otp":
					return prompt.FilterHasPrefix(modifySuggest["otp"], w, true)
				case "userpass":
					return prompt.FilterHasPrefix(modifySuggest["userpass"], w, true)
				case "creditcard":
					return prompt.FilterHasPrefix(modifySuggest["creditcard"], w, true)
				case "anytext":
					return prompt.FilterHasPrefix(modifySuggest["anytext"], w, true)
				case "anybinary":
					return prompt.FilterHasPrefix(modifySuggest["anybinary"], w, true)
				}
				return prompt.FilterHasPrefix(insertSuggest, w, true)
			}
		}
	}
	if _, ok := suggests[livePrefixState.livePrefix]; ok {
		s = suggests[livePrefixState.livePrefix]
		return prompt.FilterHasPrefix(s, w, true)
	}
	return prompt.FilterHasPrefix(s, w, true)
}

func changelivePrefix() (string, bool) {
	return livePrefixState.livePrefix, livePrefixState.isEnable
}
