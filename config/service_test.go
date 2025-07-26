package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	. "github.com/bhorvath/ddclient/config"
	"github.com/bhorvath/ddclient/mock"
)

var configFilename = "/tmp/ddclient_test.json"

// Expect passed in args to be included in returned configs.
func TestBuildsConfigsFromArgs(t *testing.T) {
	testConfig := mock.GetAppConfig()
	builtConfig, err := NewService(mock.GetAppArgs()).BuildConfig()

	if err != nil {
		t.Errorf("Got error: %v", err.Error())
	}
	if !reflect.DeepEqual(testConfig, builtConfig) {
		t.Errorf("Expected: %v; got: %v", testConfig, builtConfig)
	}
}

// Except config file to be parsed and returned as configs.
func TestBuildsConfigsFromFile(t *testing.T) {
	testCfg := mock.GetAppConfig()
	d, _ := json.Marshal(testCfg)
	ioutil.WriteFile(configFilename, d, 0644)
	defer func() { os.Remove(configFilename) }()

	a := &Args{ConfigFilePath: configFilename}
	builtCfg, err := NewService(a).BuildConfig()

	if err != nil {
		t.Errorf("Got error: %v", err.Error())
	}
	if !reflect.DeepEqual(testCfg, builtCfg) {
		t.Errorf("Expected: %v; got: %v", testCfg, builtCfg)
	}
}

// Expect error if a config file path is provided, but the file can't be found.
func TestErrorIfConfigFileNotFound(t *testing.T) {
	a := &Args{ConfigFilePath: configFilename}
	_, err := NewService(a).BuildConfig()

	if err == nil {
		t.Errorf("Expected error; got: %v", err)
	}
}

// Expect config file to be saved with given config.
// Existing saved values are kept unless a new value is explicitly specified.
func TestSavesConfigFile(t *testing.T) {
	d, _ := json.Marshal(mock.GetAppConfig())
	ioutil.WriteFile(configFilename, d, 0644)
	defer func() { os.Remove(configFilename) }()

	want := App{
		Record: Record{
			Domain:		"internet.com",
			Name: 		"test",
			Type:			"A",
		},
		Porkbun: Porkbun{
			APIKey:			"api-key",
			SecretKey:	"secret-key",
		},
	}

	a := Args{Record: want.Record, Porkbun: want.Porkbun}
	a.ConfigFilePath = configFilename
	NewService(&a).SaveConfig()
	f, err := ioutil.ReadFile(configFilename)
	if err != nil {
		t.Fatalf("Failed reading config file '%s: %s", configFilename, err.Error())
	}
	got := App{}
	json.Unmarshal(f, &got)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected: %v; got: %v", want, got)
	}
}

// Expect an error if config save is attempted without a specified filename.
func TestSaveErrorIfNoConfigFileSpecified(t *testing.T) {
	a := &Args{}
	err := NewService(a).SaveConfig()

	if !ErrorContains(err, "No config filename specified") {
		t.Errorf("Expected error; got: %v", err)
	}
}

func TestValidatesConfigReturnsNoErrorIfValid(t *testing.T) {
	testCfg := mock.GetAppConfig()
	d, _ := json.Marshal(testCfg)
	ioutil.WriteFile(configFilename, d, 0644)
	defer func() { os.Remove(configFilename) }()

	a := &Args{ConfigFilePath: configFilename}
	_, err := NewService(a).BuildConfig()

	if err != nil {
		if ErrorContains(err, "Validation failed") {
			t.Errorf("Expected no validation error; got: %v", err.Error())
		} else {
			t.Errorf("Got unexpected error: %v", err.Error())
		}
	}
}

func TestValidatesConfigReturnsErrorIfInvalid(t *testing.T) {
	a := &Args{}
	_, err := NewService(a).BuildConfig()

	if err == nil {
		t.Error("Expected validation error; got nil")
	} else {
		if !ErrorContains(err, "Validation failed") {
			t.Errorf("Got unexpected error: %v", err.Error())
		}
	}
}

// ErrorContains checks if the error message in got contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(got error, want string) bool {
	if got == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(got.Error(), want)
}
