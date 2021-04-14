package ast

import (
	"encoding/json"
)

func (o *BlockStmt) MarshalJSON() ([]byte, error) {
	type Copy BlockStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "BlockStmt",
		Copy: (*Copy)(o),
	})
}

func (o *AliasStmt) MarshalJSON() ([]byte, error) {
	type Copy AliasStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "AliasStmt",
		Copy: (*Copy)(o),
	})
}

func (o *AssignStmt) MarshalJSON() ([]byte, error) {
	type Copy AssignStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "AssignStmt",
		Copy: (*Copy)(o),
	})
}

func (o *CommentExpr) MarshalJSON() ([]byte, error) {
	type Copy CommentExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "CommentExpr",
		Copy: (*Copy)(o),
	})
}

func (o *CommentGroupExpr) MarshalJSON() ([]byte, error) {
	type Copy CommentGroupExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "CommentGroupExpr",
		Copy: (*Copy)(o),
	})
}

func (o *GroupExpr) MarshalJSON() ([]byte, error) {
	type Copy GroupExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "GroupExpr",
		Copy: (*Copy)(o),
	})
}

func (o *GroupStmt) MarshalJSON() ([]byte, error) {
	type Copy GroupStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "GroupStmt",
		Copy: (*Copy)(o),
	})
}

func (o *IdentExpr) MarshalJSON() ([]byte, error) {
	type Copy IdentExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "IdentExpr",
		Copy: (*Copy)(o),
	})
}

func (o *ImportStmt) MarshalJSON() ([]byte, error) {
	type Copy ImportStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "ImportStmt",
		Copy: (*Copy)(o),
	})
}

func (o *RowExpr) MarshalJSON() ([]byte, error) {
	type Copy RowExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "RowExpr",
		Copy: (*Copy)(o),
	})
}

func (o *RowStmt) MarshalJSON() ([]byte, error) {
	type Copy RowStmt
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "RowStmt",
		Copy: (*Copy)(o),
	})
}

func (o *SizeExpr) MarshalJSON() ([]byte, error) {
	type Copy SizeExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "SizeExpr",
		Copy: (*Copy)(o),
	})
}

func (o *StitchExpr) MarshalJSON() ([]byte, error) {
	type Copy StitchExpr
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Copy
	}{
		Type: "StitchExpr",
		Copy: (*Copy)(o),
	})
}
