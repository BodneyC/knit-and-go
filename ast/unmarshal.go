package ast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (o *BlockStmt) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		panic(e)
		return e
	}

	var stmtListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["block"], &stmtListRaw); e != nil {
		panic(e)
		return e
	}

	o.Block = make([]Stmt, len(stmtListRaw))

	for i, stmtRaw := range stmtListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*stmtRaw, &m); e != nil {
			panic(e)
			return e
		}
		switch m["type"] {
		case "AliasStmt":
			var p AliasStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Block[i] = &p
		case "AssignStmt":
			var p AssignStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Block[i] = &p
		case "RowStmt":
			var p RowStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Block[i] = &p
		case "GroupStmt":
			var p GroupStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Block[i] = &p
		case "BlockStmt":
			var p BlockStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Block[i] = &p
		default:
			return fmt.Errorf("Unknown type field %s", m["type"])
		}
	}

	if e := json.Unmarshal(*rawMap["start"], &o.Start); e != nil {
		panic(e)
		return e
	}
	if e := json.Unmarshal(*rawMap["end"], &o.End); e != nil {
		panic(e)
		return e
	}

	var ifaceMap map[string]interface{}
	if e := json.Unmarshal(b, &ifaceMap); e != nil {
		panic(e)
		return e
	}
	if val, ok := ifaceMap["length"].(string); !ok {
		return fmt.Errorf("Could not convert %s to string", ifaceMap["length"])
	} else {
		var e error
		if o.Length, e = strconv.ParseInt(val, 10, 64); e != nil {
			return fmt.Errorf("Could not convert %s to int", val)
		}
	}

	if e := json.Unmarshal(*rawMap["desc"], &o.Desc); e != nil {
		panic(e)
		return e
	}

	return nil
}

func (o *Size) UnmarshalJSON(b []byte) error {
	// At     Position        `json:"at"`
	// Ni     int16           `json:"ni"`
	// Nf     float32         `json:"nf"`
	// NId    Ident           `json:"nid"`
	// Before bool            `json:"before"`
	// Unit   MeasurementUnit `json:"unit"`

	var ifaceMap map[string]interface{}
	if e := json.Unmarshal(b, &ifaceMap); e != nil {
		panic(e)
		return e
	}
	var val string
	var e error
	var ok bool
	if val, ok = ifaceMap["ni"].(string); !ok {
		panic(fmt.Errorf("Could not convert \"ni\" (%s) to string", ifaceMap["ni"]))
	}
	if o.Ni, e = strconv.ParseInt(val, 10, 64); e != nil {
		panic(fmt.Errorf("Could parse int: %s : %w", val, e))
	}
	if val, ok = ifaceMap["nf"].(string); !ok {
		panic(fmt.Errorf("Could not convert \"nf\" (%s) to string", ifaceMap["nf"]))
	}
	if o.Nf, e = strconv.ParseFloat(val, 64); e != nil {
		panic(fmt.Errorf("Could parse float: %s : %w", val, e))
	}
	if o.Before, ok = ifaceMap["before"].(bool); !ok {
		panic(fmt.Errorf("Could not convert \"before\" (%s) to string", ifaceMap["before"]))
	}
	// if o.Before, e = strconv.ParseBool(val); e != nil {
	//   panic(fmt.Errorf("Could parse bool: %s : %w", val, e))
	// }
	if val, ok = ifaceMap["unit"].(string); !ok {
		panic(fmt.Errorf("Could not convert \"unit\" (%s) to string", ifaceMap["unit"]))
	}
	var i int64
	if i, e = strconv.ParseInt(val, 10, 64); e != nil {
		panic(fmt.Errorf("Could parse int: %s : %w", val, e))
	}
	o.Unit = MeasurementUnit(i)

	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		panic(e)
		return e
	}
	if e := json.Unmarshal(*rawMap["nid"], &o.NId); e != nil {
		return fmt.Errorf("Could not read \"nid\" field in Size : %w", e)
	}
	if e := json.Unmarshal(*rawMap["at"], &o.At); e != nil {
		return fmt.Errorf("Could not read \"at\" field in Size : %w", e)
	}

	return nil
}

func (o *AssignStmt) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		panic(e)
		return e
	}

	if e := json.Unmarshal(*rawMap["lhs"], &o.Lhs); e != nil {
		panic(e)
		return e
	}

	if e := json.Unmarshal(*rawMap["desc"], &o.Desc); e != nil {
		panic(e)
		return e
	}

	var m map[string]interface{}
	if e := json.Unmarshal(*rawMap["rhs"], &m); e != nil {
		panic(e)
		return e
	}

	switch m["type"] {
	case "Ident":
		var p Ident
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			panic(e)
			return e
		}
		o.Rhs = &p
	case "Size":
		var p Size
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			panic(e)
			return e
		}
		o.Rhs = &p
	case "StitchExpr":
		var p StitchExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			panic(e)
			return e
		}
		o.Rhs = &p
	case "RowExpr":
		var p RowExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			panic(e)
			return e
		}
		o.Rhs = &p
	case "GroupExpr":
		var p GroupExpr
		if e := json.Unmarshal(*rawMap["rhs"], &p); e != nil {
			panic(e)
			return e
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
		panic(e)
		return e
	}

	if e := json.Unmarshal(*rawMap["lbrace"], &o.LBrace); e != nil {
		panic(e)
		return e
	}
	if e := json.Unmarshal(*rawMap["rbrace"], &o.RBrace); e != nil {
		panic(e)
		return e
	}
	if e := json.Unmarshal(*rawMap["args"], &o.Args); e != nil {
		panic(e)
		return e
	}

	var stmtListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["lines"], &stmtListRaw); e != nil {
		panic(e)
		return e
	}

	o.Lines = make([]Stmt, len(stmtListRaw))

	for i, stmtRaw := range stmtListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*stmtRaw, &m); e != nil {
			panic(e)
			return e
		}
		switch m["type"] {
		case "AliasStmt":
			var p AliasStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Lines[i] = &p
		case "AssignStmt":
			var p AssignStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Lines[i] = &p
		case "RowStmt":
			var p RowStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Lines[i] = &p
		case "GroupStmt":
			var p GroupStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Lines[i] = &p
		case "BlockStmt":
			var p BlockStmt
			if e := json.Unmarshal(*stmtRaw, &p); e != nil {
				panic(e)
				return e
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
		panic(e)
		return e
	}

	if e := json.Unmarshal(*rawMap["args"], &o.Args); e != nil {
		panic(e)
		return e
	}

	var exprListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["stitches"], &exprListRaw); e != nil {
		panic(e)
		return e
	}

	o.Stitches = make([]Expr, len(exprListRaw))

	for i, exprRaw := range exprListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*exprRaw, &m); e != nil {
			panic(e)
			return e
		}
		switch m["type"] {
		case "Ident":
			var p Ident
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Stitches[i] = &p
		case "Size":
			var p Size
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Stitches[i] = &p
		case "StitchExpr":
			var p StitchExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Stitches[i] = &p
		case "RowExpr":
			var p RowExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Stitches[i] = &p
		case "GroupExpr":
			var p GroupExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Stitches[i] = &p
		default:
			return fmt.Errorf("Invalid assignment rhs: %s", m["type"])
		}
	}

	return nil
}

func (o *BracketGroup) UnmarshalJSON(b []byte) error {
	var rawMap map[string]*json.RawMessage
	if e := json.Unmarshal(b, &rawMap); e != nil {
		panic(e)
		return e
	}

	if val, ok := rawMap["args"]; !ok {
		panic("\"args\" does not exist in BracketGroup")
		return fmt.Errorf("\"args\" does not exist in BracketGroup")
	} else {
		if val == nil {
			o.Args = make([]Expr, 0)
			return nil
		}
	}

	var exprListRaw []*json.RawMessage
	if e := json.Unmarshal(*rawMap["args"], &exprListRaw); e != nil {
		panic(e)
		return e
	}

	o.Args = make([]Expr, len(exprListRaw))

	for i, exprRaw := range exprListRaw {
		var m map[string]interface{}
		if e := json.Unmarshal(*exprRaw, &m); e != nil {
			panic(e)
			return e
		}
		switch m["type"] {
		case "Ident":
			var p Ident
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Args[i] = &p
		case "Size":
			var p Size
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Args[i] = &p
		case "StitchExpr":
			var p StitchExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Args[i] = &p
		case "RowExpr":
			var p RowExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Args[i] = &p
		case "GroupExpr":
			var p GroupExpr
			if e := json.Unmarshal(*exprRaw, &p); e != nil {
				panic(e)
				return e
			}
			o.Args[i] = &p
		default:
			return fmt.Errorf("Invalid assignment rhs: %s", m["type"])
		}
	}

	return nil
}

// func (o *CommentGroup) UnmarshalJSON(b []byte) error {
// }

// func (o *AliasStmt) UnmarshalJSON(b []byte) error {
// }

// func (o *CallExpr) UnmarshalJSON(b []byte) error {
// }

// func (o *Comment) UnmarshalJSON(b []byte) error {
// }

// func (o *GroupStmt) UnmarshalJSON(b []byte) error {
// }

// func (o *Ident) UnmarshalJSON(b []byte) error {
// }

// func (o *ImportSpec) UnmarshalJSON(b []byte) error {
// }

// func (o *ParenExpr) UnmarshalJSON(b []byte) error {
// }

// func (o *RowStmt) UnmarshalJSON(b []byte) error {
// }

// func (o *Size) UnmarshalJSON(b []byte) error {
// }

// func (o *StitchExpr) UnmarshalJSON(b []byte) error {
// }
