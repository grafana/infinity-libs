package xmlframer

import (
	"errors"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/jsonframer"
)

type FramerOptions struct {
	FramerType   string
	FrameName    string
	RootSelector string
	Columns      []jsonframer.ColumnSelector
}

func ToFrame(xmlString string, options FramerOptions) (*data.Frame, error) {
	xml := strings.NewReader(xmlString)
	jsonStr, err := xj.Convert(xml)
	if err != nil {
		return nil, errors.Join(errors.New("error converting xml to grafana data frame"), err)
	}
	framerOptions := jsonframer.FramerOptions{
		FramerType:   jsonframer.FramerType(options.FramerType),
		FrameName:    options.FrameName,
		RootSelector: options.RootSelector,
		Columns:      options.Columns,
	}
	if framerOptions.FramerType == "" {
		framerOptions.FramerType = jsonframer.FramerTypeGJSON
	}
	return jsonframer.ToFrame(jsonStr.String(), framerOptions)
}
