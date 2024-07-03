package xmlframer

import (
	"errors"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/jsonframer"
)

type FramerOptions struct {
	FrameName       string
	RootSelector    string
	Columns         []jsonframer.ColumnSelector
	OverrideColumns []jsonframer.ColumnSelector
}

func xmlToJson(xmlString string) (jsonString string, err error) {
	xml := strings.NewReader(xmlString)
	jsonStr, err := xj.Convert(xml)
	if err != nil {
		return "", errors.Join(errors.New("error converting xml to grafana data frame"), err)
	}
	if jsonStr == nil {
		return "", errors.New("invalid xml content")
	}
	return jsonStr.String(), err
}

func ToFrame(xmlString string, options FramerOptions) (*data.Frame, error) {
	jsonStr, err := xmlToJson(xmlString)
	if err != nil {
		return nil, err
	}
	framerOptions := jsonframer.FramerOptions{
		FramerType:      jsonframer.FramerTypeGJSON,
		FrameName:       options.FrameName,
		RootSelector:    options.RootSelector,
		Columns:         options.Columns,
		OverrideColumns: options.OverrideColumns,
	}
	return jsonframer.ToFrame(jsonStr, framerOptions)
}
