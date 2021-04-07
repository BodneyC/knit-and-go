package ast

import (
	"encoding/json"
	"reflect"
)

func (o *BlockStmt) MarshalJSON() ([]byte, error) {
	type Copy BlockStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type:      reflect.TypeOf(*o).Name(),
		Copy: (*Copy)(o),
	})
}

func (o *AliasStmt) MarshalJSON() ([]byte, error) {
	type Copy AliasStmt
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:      reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *AssignStmt) MarshalJSON() ([]byte, error) {
	type Copy AssignStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type:       reflect.TypeOf(*o).Name(),
		Copy: (*Copy)(o),
	})
}

// func (o *CallExpr) MarshalJSON() ([]byte, error) {
//   type Copy CallExpr
//   return json.Marshal(&struct {
//     Type string `json:"type"`
//     *Copy
//   }{
//     Type:     reflect.TypeOf(*o).Name(),
//     Copy: (*Copy)(o),
//   })
// }

func (o *Comment) MarshalJSON() ([]byte, error) {
	type Copy Comment
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:    reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *CommentGroup) MarshalJSON() ([]byte, error) {
	type Copy CommentGroup
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:         reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *GroupExpr) MarshalJSON() ([]byte, error) {
	type Copy GroupExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type:      reflect.TypeOf(*o).Name(),
		Copy: (*Copy)(o),
	})
}

func (o *GroupStmt) MarshalJSON() ([]byte, error) {
	type Copy GroupStmt
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:      reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *Ident) MarshalJSON() ([]byte, error) {
	type Copy Ident
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:  reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *ImportSpec) MarshalJSON() ([]byte, error) {
	type Copy ImportSpec
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:       reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

// func (o *ParenExpr) MarshalJSON() ([]byte, error) {
//   type Copy ParenExpr
//   return json.Marshal(&struct {
//     Type string `json:"type"`
//     *Copy
//   }{
//     Type:      reflect.TypeOf(*o).Name(),
//     Copy: (*Copy)(o),
//   })
// }

func (o *RowExpr) MarshalJSON() ([]byte, error) {
	type Copy RowExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type:    reflect.TypeOf(*o).Name(),
		Copy: (*Copy)(o),
	})
}

func (o *RowStmt) MarshalJSON() ([]byte, error) {
	type Copy RowStmt
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:    reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *Size) MarshalJSON() ([]byte, error) {
	type Copy Size
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type: reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}

func (o *StitchExpr) MarshalJSON() ([]byte, error) {
	type Copy StitchExpr
  return json.Marshal(&struct {
    Type string `json:"type"`
    *Copy
  }{
    Type:       reflect.TypeOf(*o).Name(),
    Copy: (*Copy)(o),
  })
}
