package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func printOutput(v interface{}) {
	switch outputFormat {
	case "json":
		printJSON(v)
	case "yaml":
		printYAML(v)
	default:
		printTable(v)
	}
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json error: %s\n", err)
		return
	}
	fmt.Println(string(data))
}

func printYAML(v interface{}) {
	data, err := yaml.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "yaml error: %s\n", err)
		return
	}
	fmt.Print(string(data))
}

func printTable(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			fmt.Println("(no results)")
			return
		}
		printSliceTable(val)
	default:
		if val.Kind() == reflect.Struct {
			printKeyValue(val)
		} else {
			printJSON(v)
		}
	}
}

func printSliceTable(slice reflect.Value) {
	if slice.Len() == 0 {
		return
	}

	first := slice.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}
	if first.Kind() != reflect.Struct {
		printJSON(slice.Interface())
		return
	}

	fields := getTableFields(first)
	if len(fields) == 0 {
		printJSON(slice.Interface())
		return
	}

	var rows [][]string
	var header []string
	for _, f := range fields {
		header = append(header, f.label)
	}

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		var row []string
		for _, f := range fields {
			row = append(row, formatField(elem, f))
		}
		rows = append(rows, row)
	}

	colWidths := make([]int, len(header))
	for i, h := range header {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	for i, h := range header {
		fmt.Printf("%-*s  ", colWidths[i], h)
	}
	fmt.Println()

	for _, row := range rows {
		for i, cell := range row {
			fmt.Printf("%-*s  ", colWidths[i], cell)
		}
		fmt.Println()
	}
}

type tableField struct {
	label string
	name  string
}

func getTableFields(v reflect.Value) []tableField {
	t := v.Type()
	var fields []tableField

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]

		if name == "id" || name == "name" || name == "email" ||
			name == "ip" || name == "role" || name == "status" ||
			name == "connected" || name == "domain" ||
			name == "enabled" || name == "description" ||
			name == "hostname" || name == "os" ||
			name == "peers_count" || name == "resources_count" ||
			name == "valid" || name == "state" ||
			name == "used_times" || name == "last_used" ||
			name == "network" || name == "metric" ||
			name == "type" || name == "address" ||
			name == "peers" || name == "groups" ||
			name == "issued" ||
			name == "version" || name == "dns_label" ||
			name == "user_id" || name == "city_name" ||
			name == "country_code" || name == "geoname_id" ||
			name == "last_login" || name == "created_at" ||
			name == "expires" || name == "mode" ||
			name == "activity" || name == "activity_code" ||
			name == "timestamp" || name == "target_id" ||
			name == "initiator_email" || name == "initiator_name" ||
			name == "initiator_id" || name == "method" ||
			name == "host" || name == "path" ||
			name == "duration_ms" || name == "status_code" ||
			name == "source_ip" || name == "reason" ||
			name == "flow_id" || name == "reporter_id" ||
			name == "direction" || name == "protocol" ||
			name == "rx_bytes" || name == "tx_bytes" ||
			name == "rx_packets" || name == "tx_packets" ||
			name == "service_id" || name == "page" ||
			name == "page_size" || name == "total_records" ||
			name == "total_pages" || name == "data" ||
			name == "bytes_upload" || name == "bytes_download" ||
			name == "auth_method_used" || name == "subdivision_code" ||
			name == "metadata" || name == "geo_location" {
			fields = append(fields, tableField{label: name, name: name})
		}
	}

	if len(fields) == 0 {
		for i := 0; i < t.NumField() && i < 4; i++ {
			f := t.Field(i)
			tag := f.Tag.Get("json")
			if tag != "" && tag != "-" {
				name := strings.Split(tag, ",")[0]
				fields = append(fields, tableField{label: name, name: name})
			}
		}
	}

	return fields
}

func formatField(v reflect.Value, f tableField) string {
	field := v.FieldByNameFunc(func(name string) bool {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			tag := sf.Tag.Get("json")
			if tag == "" {
				continue
			}
			tagName := strings.Split(tag, ",")[0]
			if tagName == f.name && sf.Name == name {
				return true
			}
		}
		return false
	})

	if !field.IsValid() {
		return ""
	}

	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Bool:
		if field.Bool() {
			return "yes"
		}
		return "no"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", field.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.1f", field.Float())
	case reflect.Ptr:
		if field.IsNil() {
			return ""
		}
		return formatField(field.Elem(), f)
	case reflect.Slice:
		var parts []string
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			switch elem.Kind() {
			case reflect.String:
				parts = append(parts, elem.String())
			case reflect.Struct:
				if nameField := elem.FieldByName("Name"); nameField.IsValid() {
					parts = append(parts, nameField.String())
				} else if idField := elem.FieldByName("ID"); idField.IsValid() {
					parts = append(parts, idField.String())
				}
			}
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprintf("%v", field.Interface())
	}
}

func printKeyValue(v reflect.Value) {
	t := v.Type()
	maxWidth := 0
	var rows [][2]string

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		field := v.Field(i)

		val := formatFieldValue(field)
		if val == "" || val == "nil" {
			continue
		}

		if len(name) > maxWidth {
			maxWidth = len(name)
		}
		rows = append(rows, [2]string{name, val})
	}

	for _, row := range rows {
		fmt.Printf("  %-*s  %s\n", maxWidth, row[0]+":", row[1])
	}
}

func formatFieldValue(field reflect.Value) string {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return ""
		}
		return formatFieldValue(field.Elem())
	}

	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Bool:
		if field.Bool() {
			return "yes"
		}
		return "no"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", field.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.1f", field.Float())
	case reflect.Slice:
		if field.Len() == 0 {
			return ""
		}
		var parts []string
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			switch elem.Kind() {
			case reflect.String:
				parts = append(parts, elem.String())
			case reflect.Struct:
				nameField := elem.FieldByName("Name")
				idField := elem.FieldByName("ID")
				typeField := elem.FieldByName("Type")
				if nameField.IsValid() && nameField.String() != "" {
					parts = append(parts, nameField.String())
				} else if idField.IsValid() {
					entry := idField.String()
					if typeField.IsValid() {
						entry += " (" + typeField.String() + ")"
					}
					parts = append(parts, entry)
				}
			}
		}
		if len(parts) == 0 {
			return fmt.Sprintf("[%d items]", field.Len())
		}
		return strings.Join(parts, ", ")
	case reflect.Map:
		return "(map)"
	default:
		if field.CanInterface() {
			return fmt.Sprintf("%v", field.Interface())
		}
		return ""
	}
}
