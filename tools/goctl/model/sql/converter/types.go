package converter

import (
	"fmt"
	"strings"

	"github.com/zeromicro/ddl-parser/parser"
)

var unsignedTypeMap = map[string]string{
	"int":   "uint",
	"int8":  "uint8",
	"int16": "uint16",
	"int32": "uint32",
	"int64": "uint64",
}

var commonMysqlDataTypeMapInt = map[int]string{
	// For consistency, all integer types are converted to int64
	// number
	parser.Bit:       "byte",
	parser.TinyInt:   "int64",
	parser.SmallInt:  "int64",
	parser.MediumInt: "int64",
	parser.Int:       "int64",
	parser.MiddleInt: "int64",
	parser.Int1:      "int64",
	parser.Int2:      "int64",
	parser.Int3:      "int64",
	parser.Int4:      "int64",
	parser.Int8:      "int64",
	parser.Integer:   "int64",
	parser.BigInt:    "int64",
	parser.Float:     "float64",
	parser.Float4:    "float64",
	parser.Float8:    "float64",
	parser.Double:    "float64",
	parser.Decimal:   "float64",
	parser.Dec:       "float64",
	parser.Fixed:     "float64",
	parser.Numeric:   "float64",
	parser.Real:      "float64",
	// date&time
	parser.Date:      literal_0937,
	parser.DateTime:  literal_0937,
	parser.Timestamp: literal_0937,
	parser.Time:      "string",
	parser.Year:      "int64",
	// string
	parser.Char:            "string",
	parser.VarChar:         "string",
	parser.NVarChar:        "string",
	parser.NChar:           "string",
	parser.Character:       "string",
	parser.LongVarChar:     "string",
	parser.LineString:      "string",
	parser.MultiLineString: "string",
	parser.Binary:          "string",
	parser.VarBinary:       "string",
	parser.TinyText:        "string",
	parser.Text:            "string",
	parser.MediumText:      "string",
	parser.LongText:        "string",
	parser.Enum:            "string",
	parser.Set:             "string",
	parser.Json:            "string",
	parser.Blob:            "string",
	parser.LongBlob:        "string",
	parser.MediumBlob:      "string",
	parser.TinyBlob:        "string",
	// bool
	parser.Bool:    "bool",
	parser.Boolean: "bool",
}

var commonMysqlDataTypeMapString = map[string]string{
	// For consistency, all integer types are converted to int64
	// bool
	"bool":    "bool",
	"_bool":   "pq.BoolArray",
	"boolean": "bool",
	// number
	"tinyint":   "int64",
	"smallint":  "int64",
	"mediumint": "int64",
	"int":       "int64",
	"int1":      "int64",
	"int2":      "int64",
	"_int2":     literal_5987,
	"int3":      "int64",
	"int4":      "int64",
	"_int4":     literal_5987,
	"int8":      "int64",
	"_int8":     literal_5987,
	"integer":   "int64",
	"_integer":  literal_5987,
	"bigint":    "int64",
	"float":     "float64",
	"float4":    "float64",
	"_float4":   "pq.Float64Array",
	"float8":    "float64",
	"_float8":   "pq.Float64Array",
	"double":    "float64",
	"decimal":   "float64",
	"dec":       "float64",
	"fixed":     "float64",
	"real":      "float64",
	"bit":       "byte",
	// date & time
	"date":      literal_0937,
	"datetime":  literal_0937,
	"timestamp": literal_0937,
	"time":      "string",
	"year":      "int64",
	// string
	"linestring":      "string",
	"multilinestring": "string",
	"nvarchar":        "string",
	"nchar":           "string",
	"char":            "string",
	"bpchar":          "string",
	"_char":           literal_6813,
	"character":       "string",
	"varchar":         "string",
	"_varchar":        literal_6813,
	"binary":          "string",
	"bytea":           "string",
	"longvarbinary":   "string",
	"varbinary":       "string",
	"tinytext":        "string",
	"text":            "string",
	"_text":           literal_6813,
	"mediumtext":      "string",
	"longtext":        "string",
	"enum":            "string",
	"set":             "string",
	"json":            "string",
	"jsonb":           "string",
	"blob":            "string",
	"longblob":        "string",
	"mediumblob":      "string",
	"tinyblob":        "string",
	"ltree":           "[]byte",
}

// ConvertDataType converts mysql column type into golang type
func ConvertDataType(dataBaseType int, isDefaultNull, unsigned, strict bool) (string, error) {
	tp, ok := commonMysqlDataTypeMapInt[dataBaseType]
	if !ok {
		return "", fmt.Errorf("unsupported database type: %v", dataBaseType)
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned, strict), nil
}

// ConvertStringDataType converts mysql column type into golang type
func ConvertStringDataType(dataBaseType string, isDefaultNull, unsigned, strict bool) (
	goType string, isPQArray bool, err error) {
	tp, ok := commonMysqlDataTypeMapString[strings.ToLower(dataBaseType)]
	if !ok {
		return "", false, fmt.Errorf("unsupported database type: %s", dataBaseType)
	}

	if strings.HasPrefix(dataBaseType, "_") {
		return tp, true, nil
	}

	return mayConvertNullType(tp, isDefaultNull, unsigned, strict), false, nil
}

func mayConvertNullType(goDataType string, isDefaultNull, unsigned, strict bool) string {
	if !isDefaultNull {
		if unsigned && strict {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}

	switch goDataType {
	case "int64":
		return "sql.NullInt64"
	case "int32":
		return "sql.NullInt32"
	case "float64":
		return "sql.NullFloat64"
	case "bool":
		return "sql.NullBool"
	case "string":
		return "sql.NullString"
	case literal_0937:
		return "sql.NullTime"
	default:
		if unsigned {
			ret, ok := unsignedTypeMap[goDataType]
			if ok {
				return ret
			}
		}
		return goDataType
	}
}

const literal_0937 = "time.Time"

const literal_5987 = "pq.Int64Array"

const literal_6813 = "pq.StringArray"
