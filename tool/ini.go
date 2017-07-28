package tool

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Pair struct {
	Key   string
	Value string
}

type Config map[string]map[string]string

/*
| read content line by line and send into Vector
*/
func read(filename string) (v Vector, e error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("Open file filed")
	}
	defer f.Close()

	vector := NewVector()
	r := bufio.NewReader(f)

	for {
		if line, _, err := r.ReadLine(); err == nil {
			vector.Add(string(line))
			continue
		}

		return vector, nil
	}
}

/*
| parse the [.+]
| It is the father of config
*/
func parseHead(str string) (string, bool) {
	str = strings.TrimSpace(str)
	exp := `\[(.+?)\]`

	if is, _ := regexp.Match(exp, []byte(str)); !is {
		return "", false
	}

	r := regexp.MustCompile(exp)
	head := r.FindSubmatch([]byte(str))

	return string(head[1]), true
}

/*
| parse foo=joo
| It will be saved into map[head][foo]=joo
*/
func parseBody(str string) (Pair, bool) {
	str = strings.TrimSpace(str)
	exp := `(.+?)=(.+)`

	if is, _ := regexp.Match(exp, []byte(str)); !is {
		return Pair{}, false
	}

	r := regexp.MustCompile(exp)
	body := r.FindSubmatch([]byte(str))

	return Pair{
		Key:   strings.TrimSpace(string(body[1])),
		Value: strings.TrimSpace(string(body[2])),
	}, true

}

/*
|judge whether the string is empty
*/
func isEmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}

func ParseConf(filename string) (Config, error) {
	v, err := read(filename)

	if err != nil {
		return nil, fmt.Errorf("File error:%s", err.Error())
	}

	var (
		head   string
		config = make(Config)
	)

	for pos, val := range v {
		valStr := val.(string)
		//clear the space
		if isEmptyString(valStr) {
			continue
		}
		//judge where it is a note
		if strings.HasPrefix(valStr, ";") {
			continue
		}

		//check the head
		if h, hExist := parseHead(valStr); hExist {
			head = h
			config[head] = make(map[string]string)
			continue
		}
		//check the body
		if b, bExist := parseBody(valStr); bExist {
			//if no head, it must be a error
			if head == "" {
				return nil, fmt.Errorf("Parse error: No Head")
			}
			config[head][b.Key] = b.Value
			continue
		}

		//deal invalid format
		panic(fmt.Errorf("Parse error: invalid format in line:%d, content:%s", pos, valStr))
	}

	return config, nil
}
