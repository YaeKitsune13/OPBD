package backend

import (
	"example/project/backend/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=1234 dbname=cursavoi_project port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Doctor{},
		&models.Service{},
		&models.Medication{},
		&models.Pet{},
		&models.Appointment{},
		&models.Visit{},
		&models.VisitPrescription{},
		&models.WeightHistory{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %v", err)
	}

	// Добавляем FK вручную
	fkQueries := []string{
		// pets -> owners
		`ALTER TABLE pets ADD CONSTRAINT fk_pets_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// appointments -> pets
		`ALTER TABLE appointments ADD CONSTRAINT fk_appointments_pet FOREIGN KEY (pet_id) REFERENCES pets(pet_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// appointments -> doctors
		`ALTER TABLE appointments ADD CONSTRAINT fk_appointments_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(doctor_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// visits -> appointments
		`ALTER TABLE visits ADD CONSTRAINT fk_visits_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(appointment_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// visit_prescriptions -> visits
		`ALTER TABLE visit_prescriptions ADD CONSTRAINT fk_visit_prescriptions_visit FOREIGN KEY (visit_id) REFERENCES visits(visit_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// visit_prescriptions -> services
		`ALTER TABLE visit_prescriptions ADD CONSTRAINT fk_visit_prescriptions_service FOREIGN KEY (service_id) REFERENCES services(service_id) ON UPDATE CASCADE ON DELETE SET NULL`,
		// visit_prescriptions -> medications
		`ALTER TABLE visit_prescriptions ADD CONSTRAINT fk_visit_prescriptions_medication FOREIGN KEY (medication_id) REFERENCES medications(medication_id) ON UPDATE CASCADE ON DELETE SET NULL`,
		// weight_histories -> pets
		`ALTER TABLE weight_histories ADD CONSTRAINT fk_weight_histories_pet FOREIGN KEY (pet_id) REFERENCES pets(pet_id) ON UPDATE CASCADE ON DELETE CASCADE`,
		// weight_histories -> visits
		`ALTER TABLE weight_histories ADD CONSTRAINT fk_weight_histories_visit FOREIGN KEY (visit_id) REFERENCES visits(visit_id) ON UPDATE CASCADE ON DELETE SET NULL`,
		// weight_histories -> doctors
		`ALTER TABLE weight_histories ADD CONSTRAINT fk_weight_histories_doctor FOREIGN KEY (doctor_id) REFERENCES doctors(doctor_id) ON UPDATE CASCADE ON DELETE SET NULL`,
		// doctor -> users
		`ALTER TABLE doctors ADD CONSTRAINT fk_doctors_user FOREIGN KEY (user_id) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE CASCADE`,
	}

	for _, q := range fkQueries {
		if err := db.Exec(q).Error; err != nil {
			fmt.Printf("FK warning (already exists?): %v\n", err)
		}
	}

	return db, nil
}
