package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
)

// Table .
type Table struct {
	Header []string `json:"header"`
	Data   []Data   `json:"data"`
}

// Data .
type Data []string

// Options .
type Options struct {
	Header bool
	CSV    bool
	JSON   bool
	Pretty bool
	Debug  bool
}

// Result .
type Result struct {
	ResultData map[string]interface{} `json:"data,omitempty"`
	Result     string                 `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Info       string                 `json:"info,omitempty"`
}

// RenderJSON .
func RenderJSON(data interface{}, opts Options) {
	var jsonBytes []byte
	var err error
	if opts.Pretty {
		jsonBytes, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
	} else {
		jsonBytes, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(string(jsonBytes))
}

// RenderError .
func RenderError(errorMsg string, opts Options) {
	if opts.JSON {
		jsonData := Result{
			Error: trimQuotes(errorMsg),
		}
		RenderJSON(jsonData, opts)
	} else {
		//fmt.Println(fmt.Sprintf("Error: %s", aurora.Yellow(trimQuotes(errorMsg))))
		fmt.Println("Error:", trimQuotes(errorMsg))
	}
}

// RenderInfo .
func RenderInfo(infoMsg string, opts Options) {
	if opts.JSON {
		jsonData := Result{
			Info: trimQuotes(infoMsg),
		}
		RenderJSON(jsonData, opts)
	} else {
		fmt.Println("Info:", trimQuotes(infoMsg))
	}
}

// RenderResult .
func RenderResult(result Result, opts Options) {
	if opts.JSON {
		RenderJSON(result, opts)
	} else {
		if trimQuotes(result.Result) == "success" {
			fmt.Println(fmt.Sprintf("Result: %s", aurora.Green(trimQuotes(result.Result))))
			if len(result.ResultData) != 0 {
				for k, v := range result.ResultData {
					fmt.Println(fmt.Sprintf("%s: %v", k, v))
				}
			}
		} else {
			fmt.Println(fmt.Sprintf("Result: %s", aurora.Yellow(trimQuotes(result.Result))))
			if len(result.ResultData) != 0 {
				for k, v := range result.ResultData {
					fmt.Println(fmt.Sprintf("%s: %v", k, v))
				}
			}
		}
	}

}

// RenderOutput .
func RenderOutput(data Table, opts Options) {
	if opts.Debug {
		fmt.Println(fmt.Sprintf("%s", aurora.Yellow("Final result:")))
	}
	if opts.JSON {
		// really basic tabledata to json implementation
		var rawData []interface{}
		for _, dataValues := range data.Data {
			jsonData := make(map[string]interface{})
			for indexID, dataValue := range dataValues {
				dataHeader := strings.Replace(strings.ToLower(data.Header[indexID]), " ", "-", -1)
				jsonData[dataHeader] = dataValue
			}
			rawData = append(rawData, jsonData)
		}
		returnedData := map[string]interface{}{
			"data": rawData,
		}
		RenderJSON(returnedData, opts)
	} else {
		// otherwise render a table
		table := tablewriter.NewWriter(os.Stdout)
		opts.Header = !opts.Header
		if opts.Header {
			table.SetHeader(data.Header)
		}
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		if opts.CSV {
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetCenterSeparator("")
			table.SetRowSeparator("")
			table.SetColumnSeparator(",")
		} else {
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetCenterSeparator("")
			table.SetRowSeparator("")
			table.SetColumnSeparator(",")
			table.SetTablePadding("\t") // pad with tabs
			table.SetNoWhiteSpace(true)
		}
		for _, rowData := range data.Data {
			table.Append(rowData)
		}
		table.Render()
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
