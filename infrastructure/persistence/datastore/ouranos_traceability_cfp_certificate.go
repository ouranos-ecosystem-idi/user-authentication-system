package datastore

import "fmt"

func (r *ouranosRepository) DeleteCFPCertificateByCFPID(cfpID string) error {
	result := r.db.Unscoped().Table("cfp_certificates").Where("cfp_id = ?", cfpID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table cfp_certificates: %v", result.Error)
	}
	return nil
}
