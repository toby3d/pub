package domain

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// Action represent a single micropub action.
type Action struct {
	action string
}

var (
	ActionUnd      Action = Action{action: ""}         // "und"
	ActionCreate   Action = Action{action: "create"}   // "create"
	ActionUpdate   Action = Action{action: "update"}   // "update"
	ActionDelete   Action = Action{action: "delete"}   // "delete"
	ActionUndelete Action = Action{action: "undelete"} // "undelete"
)

var ErrActionSyntax error = Error{
	Description: "unsupported action emun",
	Frame:       xerrors.Caller(1),
	Code:        http.StatusBadRequest,
}

var stringsActions = map[string]Action{
	ActionCreate.action:   ActionCreate,
	ActionUpdate.action:   ActionUpdate,
	ActionDelete.action:   ActionDelete,
	ActionUndelete.action: ActionUndelete,
}

func ParseAction(raw string) (Action, error) {
	if a, ok := stringsActions[raw]; ok {
		return a, nil
	}

	return ActionUnd, fmt.Errorf("cannot parse string as Action enum: %w", ErrActionSyntax)
}

func (a Action) String() string {
	if a.action == "" {
		return "und"
	}

	return a.action
}

func (a Action) GoString() string {
	return "domain.Action(" + a.String() + ")"
}
