package datastore

import (
	"fmt"
)

func (r *ouranosRepository) DeletePartsStructure(traceID string) error {
	result := r.db.Unscoped().Table("parts_structures").Where("trace_id = ?", traceID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table parts_structures: %v", result.Error)
	}
	return nil
}
