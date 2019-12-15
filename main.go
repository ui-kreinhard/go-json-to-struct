package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/segmentio/ksuid"
	"log"
	"os"
	"reflect"
	"strings"
)

func getStructName() string {
	id := ksuid.New()
	return "gen_" + id.String()
}

func getIdentifierPart(attributeName string, jsonPartValue interface{}) (*jen.Statement) {
	upperedFirst := strings.Title(attributeName)
	switch jsonPartValue.(type) {
	case []interface{}:
		return jen.Id(upperedFirst).Index()
	default:
		return jen.Id(upperedFirst)
	}
}

func getType(jsonPartValue interface{}) interface{}{
	switch jsonPartValue.(type) {
	case []interface{}:
		casted := jsonPartValue.([]interface{})
		return casted[0]
	default:
		return jsonPartValue
	}
}

func getCodePart(attributeName string, jsonPartValue interface{}, parent *jen.File) (jen.Code, string) {
	tag := map[string]string{"json": attributeName}
	identifierPart := getIdentifierPart(attributeName, jsonPartValue)
	modulatedType := getType(jsonPartValue)
	switch v := modulatedType.(type) {

	case int:
		return identifierPart.Int().Tag(tag), ""
	case float64:
		if float64(v)-float64(int64(v)) == 0.0 {
			return identifierPart.Int().Tag(tag), ""
		} else {
			return identifierPart.Float64().Tag(tag), ""
		}

	case string:
		return identifierPart.String().Tag(tag), ""
	case bool:
		return identifierPart.Bool().Tag(tag), ""
	case map[string]interface{}:
		_, subStructName := gen(modulatedType.(map[string]interface{}),  attributeName, parent)
		return identifierPart.Id(subStructName).Tag(tag), subStructName
	default:
		fmt.Println("I don't know, ask stackoverflow.", reflect.TypeOf(jsonPartValue))
	}
	return nil, ""
}

func gen(jsonPart map[string]interface{}, structName string, parent *jen.File) (jen.Code, string) {
	fields := []jen.Code{}
	for k, _ := range jsonPart {
		codePart, _ := getCodePart(k, jsonPart[k], parent)
		fields = append(fields, codePart)
	}
	return parent.Type().Id("Gen_"+structName).Struct(fields...), "Gen_" + structName
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(text), &jsonMap)
	if err != nil {
		log.Println(err)
	}
	f := jen.NewFile("types")
	gen(jsonMap, "root", f)
	fmt.Printf("%#v\n", f)
}
