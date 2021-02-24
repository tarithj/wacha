package utils

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"regexp"
	"strings"
)

func GetParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func CanBanMembers(user *discordgo.User) bool {
	list := strings.Split(os.Getenv("BanPermIdList"), " ")
	for i := 0; i < len(list); i++ {
		if user.ID == list[i] {
			return true
		}
	}
	return false
}
