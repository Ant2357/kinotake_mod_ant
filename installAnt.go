package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/zserge/lorca"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func addAntToScenario(path string) (string, error) {
	commonErrorMsg := "Antをシナリオに追加できませんでした"

	fp, err := os.Open(path)
	if err != nil {
		return commonErrorMsg, err
	}
	defer fp.Close()

	// UTF-16LE対応
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	reader := bufio.NewReader(decoder.NewDecoder().Reader(fp))

	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return commonErrorMsg, err
	}

	match, _ := regexp.MatchString("roam = ant,", string(decoded))
	if match {
		return "Antがインストール済みです", errors.New("Antがインストール済みです")
	}

	outputCode := strings.ReplaceAll(string(decoded), "roam = ", "roam = ant,")

	fpW, err := os.Create(path)
	if err != nil {
		return commonErrorMsg, err
	}
	defer fpW.Close()

	writer := transform.NewWriter(fpW, decoder.NewEncoder())
	defer writer.Close()

	// UTF-16LEで出力
	_, err = writer.Write([]byte(outputCode))
	if err != nil {
		return commonErrorMsg, err
	}

	return "Antをシナリオに追加しました", nil
}

func newUi(msg string) {
	templateHtml := `
	<html>
		<head><meta charset="utf-8"><title>きのたけ戦争if Ant MOD</title></head>
		<body>MSG</body>
	</html>
	`

	html := strings.ReplaceAll(templateHtml, "MSG", msg)

	ui, err := lorca.New("data:text/html,"+url.PathEscape(html), "", 320, 240)
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	<-ui.Done()
}

func main() {
	paths := [3]string{"kinotake/script/sc1.dat", "kinotake/script/sc2.dat", "kinotake/script/scenario_rnd/sc_rnd.dat"}

	msg := ""
	var err error
	for _, path := range paths {
		msg, err = addAntToScenario(path)
		if err != nil {
			newUi(msg)
			return
		}
	}

	newUi(msg)
}
