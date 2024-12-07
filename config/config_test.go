package config

import "testing"

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}

func TestInitConfiguration_WithCorrectFileName(t *testing.T) {
	cfg, err := initConfiguration("config")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}

func TestInitConfiguration_WithEmptyFileName(t *testing.T) {
	cfg, err := initConfiguration("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}

func TestInitConfiguration_WithWrongFileName(t *testing.T) {
	_, err := initConfiguration("configuration")
	if err == nil {
		t.Fatal("Didn't get any error")
	}
}

func TestInitConfiguration_WithAlteredFileName(t *testing.T) {
	_, err := initConfiguration("config_err")
	if err == nil {
		t.Fatal("Didn't get any error")
	}
}
