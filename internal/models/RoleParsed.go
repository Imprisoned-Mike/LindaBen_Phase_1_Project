package models

import (
	"strconv"
	"strings"
)

type RoleParsed struct {
	Role     string //admin, school_admin, vendor_admin
	EntityID *uint
}

func ParseRoles(roles string) []RoleParsed {
	var parsedRoles []RoleParsed

	for role := range strings.SplitSeq(roles, ",") {
		if strings.Contains(role, ":") {
			eid := strings.SplitN(role, ":", 2)[1]
			eidParsed, _ := strconv.Atoi(eid)
			eidParsedUint := uint(eidParsed)

			parsedRoles = append(parsedRoles, RoleParsed{
				Role:     strings.SplitN(role, ":", 2)[0],
				EntityID: &eidParsedUint,
			})
		} else {
			parsedRoles = append(parsedRoles, RoleParsed{
				Role:     role,
				EntityID: nil,
			})
		}
	}
	return parsedRoles
}
