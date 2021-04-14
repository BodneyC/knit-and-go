package ast

import (
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/bodneyc/knit-and-go/util"
)

func (o *BlockStmt) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	var stmtListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["block"], &stmtListRaw); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	o.Block = make([]Stmt, len(stmtListRaw))

	for i, stmtRaw := range stmtListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*stmtRaw, &m); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		switch m["type"] {
		case "AliasStmt":
			var p AliasStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Block[i] = &p
		case "AssignStmt":
			var p AssignStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Block[i] = &p
		case "RowStmt":
			var p RowStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Block[i] = &p
		case "GroupStmt":
			var p GroupStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Block[i] = &p
		case "BlockStmt":
			var p BlockStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Block[i] = &p
		default:
			return fmt.Errorf("Unknown type field %s", m["type"])
		}
	}

	if e := json.Unmarshal(*rawMap["start"], &o.Start); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	if e := json.Unmarshal(*rawMap["end"], &o.End); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	var ifaceMap map[string]interface{}
	if e := json.Unmarshal(b, &ifaceMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	// if val, ok := ifaceMap["length"].(string); !ok {
	//   return fmt.Errorf("Could not convert %s to string", ifaceMap["length"])
	// } else {
	//   var e error
	//   if o.Length, e = strconv.ParseInt(val, 10, 64); e != nil {
	//     return fmt.Errorf("Could not convert %s to int", val)
	//   }
	// }

	if e := json.Unmarshal(*rawMap["desc"], &o.Desc); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	return nil
}

// Kill me, please
func (o *SizeExpr) UnmarshalJSON(b []byte) error {
	var ifaceMap map[string]interface{}
	if e := json.Unmarshal(b, &ifaceMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	var val string
	var e error
	var ok bool
	if val, ok = ifaceMap["ni"].(string); !ok {
		return fmt.Errorf("Could not convert \"ni\" (%s) to string", ifaceMap["ni"])
	}
	if o.Ni, e = strconv.ParseInt(val, 10, 64); e != nil {
		return fmt.Errorf("Could parse int: %w%s", e, StackLine())
	}
	if val, ok = ifaceMap["nf"].(string); !ok {
		return fmt.Errorf("Could not convert \"nf\" (%s) to string", ifaceMap["nf"])
	}
	if o.Nf, e = strconv.ParseFloat(val, 64); e != nil {
		return fmt.Errorf("Could parse float: %w%s", e, StackLine())
	}
	if o.Before, ok = ifaceMap["before"].(bool); !ok {
		return fmt.Errorf("Could not convert \"before\" (%s) to string", ifaceMap["before"])
	}
	if val, ok = ifaceMap["unit"].(string); !ok {
		return fmt.Errorf("Could not convert \"unit\" (%s) to string", ifaceMap["unit"])
	}
	var i int64
	if i, e = strconv.ParseInt(val, 10, 64); e != nil {
		return fmt.Errorf("Could parse int: %w%s", e, StackLine())
	}
	o.Unit = MeasurementUnit(i)

	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	if e := json.Unmarshal(*rawMap["nid"], &o.Id); e != nil {
		return fmt.Errorf("Could not read \"nid\" field in SizeExpr : %w", e)
	}
	if e := json.Unmarshal(*rawMap["at"], &o.At); e != nil {
		return fmt.Errorf("Could not read \"at\" field in SizeExpr : %w", e)
	}

	return nil
}

func (o *AssignStmt) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	if e := json.Unmarshal(*rawMap["lhs"], &o.Lhs); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	if e := json.Unmarshal(*rawMap["desc"], &o.Desc); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	var m map[string]interface{}
	if e := json.Unmarshal(*rawMap["rhs"], &m); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	switch m["type"] {
	case "IdentExpr":
		var p IdentExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		o.Rhs = &p
	case "SizeExpr":
		var p SizeExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		o.Rhs = &p
	case "StitchExpr":
		var p StitchExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		o.Rhs = &p
	case "RowExpr":
		var p RowExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		o.Rhs = &p
	case "GroupExpr":
		var p GroupExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		o.Rhs = &p
	default:
		return fmt.Errorf("Invalid assignment rhs: %s", m["type"])
	}

	return nil
}

func (o *GroupExpr) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	if e := json.Unmarshal(*rawMap["lbrace"], &o.LBrace); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	if e := json.Unmarshal(*rawMap["rbrace"], &o.RBrace); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}
	if e := json.Unmarshal(*rawMap["args"], &o.Args); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	var stmtListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["lines"], &stmtListRaw); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	o.Lines = make([]Stmt, len(stmtListRaw))

	for i, stmtRaw := range stmtListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*stmtRaw, &m); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		switch m["type"] {
		case "AliasStmt":
			var p AliasStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Lines[i] = &p
		case "AssignStmt":
			var p AssignStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Lines[i] = &p
		case "RowStmt":
			var p RowStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Lines[i] = &p
		case "GroupStmt":
			var p GroupStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Lines[i] = &p
		case "BlockStmt":
			var p BlockStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Lines[i] = &p
		default:
			return fmt.Errorf("Unknown type field %s", m["type"])
		}
	}

	return nil
}

func (o *RowExpr) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	if e := json.Unmarshal(*rawMap["args"], &o.Args); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	var exprListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["stitches"], &exprListRaw); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	o.Stitches = make([]Expr, len(exprListRaw))

	for i, exprRaw := range exprListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*exprRaw, &m); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		switch m["type"] {
		case "IdentExpr":
			var p IdentExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Stitches[i] = &p
		case "SizeExpr":
			var p SizeExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Stitches[i] = &p
		case "StitchExpr":
			var p StitchExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Stitches[i] = &p
		case "RowExpr":
			var p RowExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Stitches[i] = &p
		case "GroupExpr":
			var p GroupExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Stitches[i] = &p
		default:
			return fmt.Errorf("Invalid assignment rhs: %s", m["type"])
		}
	}

	return nil
}

func (o *Brackets) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	if val, ok := rawMap["args"]; !ok {
		return fmt.Errorf("\"args\" does not exist in Brackets")
	} else {
		if val == nil {
			o.Args = make([]Expr, 0)
			return nil
		}
	}

	var exprListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["args"], &exprListRaw); e != nil {
		return fmt.Errorf("%w%s", e, StackLine())
	}

	o.Args = make([]Expr, len(exprListRaw))

	for i, exprRaw := range exprListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*exprRaw, &m); e != nil {
			return fmt.Errorf("%w%s", e, StackLine())
		}
		switch m["type"] {
		case "IdentExpr":
			var p IdentExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Args[i] = &p
		case "SizeExpr":
			var p SizeExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Args[i] = &p
		case "StitchExpr":
			var p StitchExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Args[i] = &p
		case "RowExpr":
			var p RowExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Args[i] = &p
		case "GroupExpr":
			var p GroupExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				return fmt.Errorf("%w%s", e, StackLine())
			}
			o.Args[i] = &p
		default:
			return fmt.Errorf("Invalid assignment rhs: %s", m["type"])
		}
	}

	return nil
}
