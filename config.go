package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"os/user"
	"strings"
)

const (
	pathSeparator   = string(os.PathSeparator)
	configName      = pathSeparator + "tgsend.ini"
	botsSection     = "bots"
	chatsSection    = "chats"
	defaultsSection = "defaults"
	defaultBotKey   = "bot"
	defaultChatKey  = "chat"
)

type (
	aliasName       string
	aliasValue      string
	aliasCollection map[aliasName]aliasValue

	config struct {
		defaultBot  aliasName
		defaultChat aliasName
		bots        aliasCollection
		chats       aliasCollection
	}
)

func (a *aliasName) String() string {
	return string(*a)
}

func (a *aliasName) Set(value string) error {
	*a = aliasName(value)
	return nil
}

func (c *config) Init(file *ini.File) error {
	var err error

	if c.bots, err = getAliasCollectionFromSection(file, botsSection); err != nil {
		return err
	}

	if c.chats, err = getAliasCollectionFromSection(file, chatsSection); err != nil {
		return err
	}

	var defaultSection *ini.Section
	if defaultSection, err = file.GetSection(defaultsSection); err != nil {
		return err
	}

	if c.defaultBot, err = getDefaultAliasFromCollection(defaultSection, defaultBotKey, c.bots); err != nil {
		return err
	}

	if c.defaultChat, err = getDefaultAliasFromCollection(defaultSection, defaultChatKey, c.chats); err != nil {
		return err
	}

	return nil
}

func (c *config) GetBotId(name aliasName) aliasValue {
	return getAliasValue(c.bots, name)
}

func (c *config) GetChatId(name aliasName) aliasValue {
	return getAliasValue(c.chats, name)
}

func getAliasValue(collection aliasCollection, name aliasName) aliasValue {
	if value, ok := collection[name]; ok {
		return value
	}

	return ""
}

func getAliasCollectionFromSection(file *ini.File, sectionName string) (aliasCollection, error) {
	section, err := file.GetSection(sectionName)
	if err != nil {
		return nil, err
	}

	collection := aliasCollection{}

	for _, option := range section.Keys() {
		collection[aliasName(option.Name())] = aliasValue(option.Value())
	}

	return collection, nil
}

func getDefaultAliasFromCollection(section *ini.Section, keyName string, collection aliasCollection) (aliasName, error) {
	key, err := section.GetKey(keyName)
	if err != nil {
		return "", err
	}

	alias := aliasName(key.Value())
	if _, ok := collection[alias]; !ok {
		return "", fmt.Errorf("unknown alias in defaults section: %s", alias)
	}

	return alias, nil
}

func getPossibleConfigs() []string {
	var paths = []string{fmt.Sprintf("%setc%s", pathSeparator, configName)}
	if currentUser, err := user.Current(); err == nil {
		paths = append(paths, fmt.Sprintf("%s%s.config%s", currentUser.HomeDir, pathSeparator, configName))
	}

	return paths
}

func readConfig(configFile string) (*config, error) {
	file, err := ini.Load(configFile)
	if err != nil {
		return nil, err
	}

	var configObject config
	if err := configObject.Init(file); err != nil {
		return nil, err
	}

	return &configObject, nil
}

func getConfig() (*config, error) {
	paths := getPossibleConfigs()
	var foundPath string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			foundPath = path
			break
		}
	}

	if foundPath == "" {
		return nil, fmt.Errorf("config file is not found in following paths: %s", strings.Join(paths, ", "))
	}

	return readConfig(foundPath)
}
