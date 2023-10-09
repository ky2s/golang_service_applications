package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"snapin-form/config"
	"snapin-form/objects"
	"strings"

	"gorm.io/gorm"
)

type ShortenUrlModels interface {
	ShortenImageURL(link string) (objects.ShortURLResponse, error)
	GetLinkReal(alias string) (objects.ShortURLResponse, error)
}

type shortenurlConnection struct {
	conf config.Configurations
	db   *gorm.DB
}

func NewShortenUrlModels(dbg *gorm.DB) ShortenUrlModels {
	return &shortenurlConnection{
		db: dbg,
	}
}

func (con *shortenurlConnection) ShortenImageURL(link string) (objects.ShortURLResponse, error) {
	url := "https://link.snap-in.co.id/api/v1/links"
	key := "akfKurRqh8LVTfgAontt0Enu61eDe1tO5luJ4S7ByhqYe3H6YvY9wa9Bu82L"
	// CONFIG_SHORTEN_BASE_URL := con.conf.SHORTEN_BASE_URL
	// CONFIG_SHORTEN_KEY := con.conf.SHORTEN_KEY
	method := "POST"
	payload := strings.NewReader("url=" + link)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return objects.ShortURLResponse{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Cookie", "dark_mode=0")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return objects.ShortURLResponse{}, err

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return objects.ShortURLResponse{}, err

	}
	var response objects.ShortURLResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return objects.ShortURLResponse{}, err

	}
	fmt.Println(response)
	return response, nil
}

func (con *shortenurlConnection) GetLinkReal(alias string) (objects.ShortURLResponse, error) {
	url := "https://link.snap-in.co.id/api/v1/alias/" + alias
	key := "akfKurRqh8LVTfgAontt0Enu61eDe1tO5luJ4S7ByhqYe3H6YvY9wa9Bu82L"
	method := "GET"

	// fmt.Println(url)
	// os.Exit(0)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("err 1", err)
		return objects.ShortURLResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "dark_mode=0")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("err 2", err)
		return objects.ShortURLResponse{}, err

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err 3", err)
		return objects.ShortURLResponse{}, err

	}
	fmt.Println("string(body)-------", string(body))
	var response objects.ShortURLResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("err 4", err)
		return objects.ShortURLResponse{}, err

	}
	if response.Status != 200 {
		fmt.Println("err 5", err)
		return objects.ShortURLResponse{}, err
	}

	return response, nil
}
