package common

import (
	"shortlinks/models"
	"fmt"
)

var apps = make(map[string]app)

type app struct {
	id uint32
	key string
	secret string
}

func init() {
	engine := RegisterMysql()
	
	appModels := make([]models.Application, 0)
	err := engine.Where("status=?", 0).Cols("id", "app_key", "app_secret").Find(&appModels)
	if err != nil {
		ERROR.Println(err)
		panic(err)
	}
	
	for _, appModel := range appModels {
		tmp := app{appModel.Id, appModel.AppKey, appModel.AppSecret}
		apps[appModel.AppKey] = tmp
	}
}

func GetAppSecret(key string) (string, error) {
	tmp, ok := apps[key]
	if !ok {
		return "", fmt.Errorf("the app does not exist")
	}
	
	return tmp.secret, nil
}

func GetAppId(key string) (uint32, error) {
	tmp, ok := apps[key]
	if !ok {
		return 0, fmt.Errorf("the app does not exist")
	}
	
	return tmp.id, nil
}