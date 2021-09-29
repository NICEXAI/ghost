package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	// such as: @Lazy range:data_list
	easyMode = `^\w+$`
	// such as: @Lazy range:data_list|>ghost
	generalMode = `^(\w+)|>(\w+)$`
	// such as: @Lazy range:data_list|>(name>ghost,age>123)
	complexMode = `^(\w+)|>\((\S+)\)$`
)

// ParseAndExecuteRangeExpr parse and execute range expr
func ParseAndExecuteRangeExpr(order RangeCommand, attrs map[string]interface{}) (RangeCommand, error) {
	var (
		reg   *regexp.Regexp
		oList [][]string
	)

	reg = regexp.MustCompile(easyMode)

	if reg.MatchString(order.Expr) {
		value := attrs[order.Expr]

		switch list := value.(type) {
		case []map[string]interface{}:
			order.Loop = len(list)
			return order, nil
		case []string:
			order.Loop = len(list)
			return order, nil
		case []int:
			order.Loop = len(list)
			return order, nil
		}

		return order, errors.New(fmt.Sprintf("%v's %v type is not support", order.Expr, reflect.TypeOf(value)))
	}

	reg = regexp.MustCompile(generalMode)
	oList = reg.FindAllStringSubmatch(order.Expr, -1)

	if len(oList) == 2 {
		variable := oList[0][1]
		target := oList[1][2]
		value := attrs[variable]

		switch dataList := value.(type) {
		case []map[string]interface{}:
			order.Loop = len(dataList)
			actionList := make([]RangeAction, 0)

			for _, data := range dataList {
				bData, err := json.Marshal(data)
				if err != nil {
					return order, err
				}
				actionList = append(actionList, RangeAction{
					Target: target,
					Value:  string(bData),
				})
			}

			order.Action = append(order.Action, actionList)
			return order, nil
		case []string:
			order.Loop = len(dataList)
			actionList := make([]RangeAction, 0)

			for _, data := range dataList {
				actionList = append(actionList, RangeAction{
					Target: target,
					Value:  data,
				})
			}

			order.Action = append(order.Action, actionList)
			return order, nil
		case []int:
			order.Loop = len(dataList)
			actionList := make([]RangeAction, 0)

			for _, data := range dataList {
				actionList = append(actionList, RangeAction{
					Target: target,
					Value:  strconv.Itoa(data),
				})
			}

			order.Action = append(order.Action, actionList)
			return order, nil
		}

		return order, errors.New(fmt.Sprintf("%v's %v type is not support", order.Expr, reflect.TypeOf(oList)))
	}

	reg = regexp.MustCompile(complexMode)
	oList = reg.FindAllStringSubmatch(order.Expr, -1)

	if len(oList) == 2 {
		variable := oList[0][1]
		subExpr := oList[1][2]
		value := attrs[variable]

		dataList, ok := value.([]map[string]interface{})
		if !ok {
			return order, errors.New(fmt.Sprintf("%v data type must be []map[string]interface{}, not allow %v", variable, reflect.TypeOf(value)))
		}

		order.Loop = len(dataList)

		for _, oTag := range strings.Split(subExpr, ",") {
			tagInfo := strings.Split(oTag, ">")
			if len(tagInfo) != 2 {
				return order, errors.New(fmt.Sprintf("%v is invalid", subExpr))
			}
			actionList := make([]RangeAction, 0)
			for _, data := range dataList {
				actionList = append(actionList, RangeAction{
					Target: tagInfo[1],
					Value:  fmt.Sprintf("%v", data[tagInfo[0]]),
				})
			}
			order.Action = append(order.Action, actionList)
		}

		return order, nil
	}

	return order, errors.New("expr " + order.Expr + " not support")
}
