package main

import (
	"flag"
	"reflect"
)

type MODE int

const (
	Server_Mode MODE = iota
	Client_Mode
)

type TRANSCODER int

const (
	Ascii_Transcoder TRANSCODER = iota
	Hex_Transcoder
	Octal_Transcoder
)

type Config struct {
	Mode       MODE       `flag:"m" default:"1" message:"service mode (client or server)"`
	Port       uint16     `flag:"p" default:"0" message:"server port"`
	Host       string     `flag:"h" default:"0.0.0.0" message:"server host or ip"`
	Transcoder TRANSCODER `flag:"t" default:"0" message:"transcoder format like(ascii hex octal)"`
}

func ParseConfig() {
	config := Config{}
	configType := reflect.TypeOf(config)

	flags := make([]string, configType.NumField())
	for i := 0; i < configType.NumField(); i++ {
		tag := configType.Field(i).Tag

		flagVal := tag.Get("flag")
		defaultVal := tag.Get("default")
		messageVal := tag.Get("message")

		flag.StringVar(&flags[i], flagVal, defaultVal, messageVal)
	}

    flag.Parse()

	for i := 0; i < configType.NumField(); i++ {
		_ = configType.Field(i).Type.Kind()

	}
}
