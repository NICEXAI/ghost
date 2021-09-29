package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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
		reg *regexp.Regexp
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

	if reg.MatchString(order.Expr) {
		oList := reg.FindAllStringSubmatch(order.Expr, -1)
		variable := oList[0][1]
		target := oList[1][2]
		value := attrs[variable]

		switch dataList := value.(type) {
		case []map[string]interface{}:
			order.Loop = len(dataList)

			for _, data := range dataList {
				bData, err := json.Marshal(data)
				if err != nil {
					return order, err
				}
				order.Action = append(order.Action, RangeAction{
					Target: target,
					Value:  string(bData),
				})
			}

			return order, nil
		case []string:
			order.Loop = len(dataList)

			for _, data := range dataList {
				order.Action = append(order.Action, RangeAction{
					Target: target,
					Value:  data,
				})
			}

			return order, nil
		case []int:
			for _, data := range dataList {
				order.Action = append(order.Action, RangeAction{
					Target: target,
					Value:  strconv.Itoa(data),
				})
			}

			return order, nil
		}

		return order, errors.New(fmt.Sprintf("%v's %v type is not support", order.Expr, reflect.TypeOf(oList)))
	}

	return order, errors.New("expr " + order.Expr + " not support")
}
