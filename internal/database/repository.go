package database

import (
	"time"

	"github.com/LimerDev/worklog/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{db: DB}
}

func (r *Repository) CreateTimeEntry(entry *models.TimeEntry) error {
	return r.db.Create(entry).Error
}

func (r *Repository) GetTimeEntriesByMonth(year int, month time.Month) ([]models.TimeEntry, error) {
	var entries []models.TimeEntry

	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	err := r.db.Preload("Project.Customer").Preload("Consultant").
		Where("date >= ? AND date < ?", startDate, endDate).
		Order("date asc").
		Find(&entries).Error

	return entries, err
}

func (r *Repository) GetAllTimeEntries() ([]models.TimeEntry, error) {
	var entries []models.TimeEntry
	err := r.db.Preload("Project.Customer").Preload("Consultant").Order("date desc").Find(&entries).Error
	return entries, err
}

// GetTimeEntriesByFilters retrieves time entries with flexible filtering
func (r *Repository) GetTimeEntriesByFilters(consultantName, projectName, customerName string, startDate, endDate time.Time) ([]models.TimeEntry, error) {
	var entries []models.TimeEntry
	query := r.db.Preload("Project.Customer").Preload("Consultant")

	// Filter by consultant
	if consultantName != "" {
		query = query.Where("consultant_id IN (SELECT id FROM consultants WHERE name = ?)", consultantName)
	}

	// Filter by project
	if projectName != "" {
		query = query.Where("project_id IN (SELECT id FROM projects WHERE name = ?)", projectName)
	}

	// Filter by customer
	if customerName != "" {
		query = query.Where("project_id IN (SELECT p.id FROM projects p JOIN customers c ON c.id = p.customer_id WHERE c.name = ?)", customerName)
	}

	// Filter by date range
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("date >= ? AND date < ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("date >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("date < ?", endDate)
	}

	err := query.Order("date asc").Find(&entries).Error
	return entries, err
}

func (r *Repository) DeleteTimeEntry(id uint) error {
	return r.db.Delete(&models.TimeEntry{}, id).Error
}

// Customer methods

func (r *Repository) CreateCustomer(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *Repository) GetOrCreateCustomer(name string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("name = ?", name).First(&customer).Error
	if err == gorm.ErrRecordNotFound {
		customer = models.Customer{Name: name, Active: true}
		err = r.db.Create(&customer).Error
	}
	return &customer, err
}

func (r *Repository) GetAllCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.Order("name asc").Find(&customers).Error
	return customers, err
}

func (r *Repository) GetCustomerByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("Projects").First(&customer, id).Error
	return &customer, err
}

// Project methods

func (r *Repository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *Repository) GetOrCreateProject(name string, customerID uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("name = ? AND customer_id = ?", name, customerID).First(&project).Error
	if err == gorm.ErrRecordNotFound {
		project = models.Project{Name: name, CustomerID: customerID, Active: true}
		err = r.db.Create(&project).Error
	}
	return &project, err
}

func (r *Repository) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("Customer").Order("name asc").Find(&projects).Error
	return projects, err
}

func (r *Repository) GetProjectsByCustomer(customerID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("customer_id = ?", customerID).Order("name asc").Find(&projects).Error
	return projects, err
}

func (r *Repository) GetProjectByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Customer").First(&project, id).Error
	return &project, err
}

// Consultant methods

func (r *Repository) CreateConsultant(consultant *models.Consultant) error {
	return r.db.Create(consultant).Error
}

func (r *Repository) GetOrCreateConsultant(name string) (*models.Consultant, error) {
	var consultant models.Consultant
	err := r.db.Where("name = ?", name).First(&consultant).Error
	if err == gorm.ErrRecordNotFound {
		consultant = models.Consultant{Name: name, Active: true}
		err = r.db.Create(&consultant).Error
	}
	return &consultant, err
}

func (r *Repository) GetAllConsultants() ([]models.Consultant, error) {
	var consultants []models.Consultant
	err := r.db.Order("name asc").Find(&consultants).Error
	return consultants, err
}

func (r *Repository) GetConsultantByID(id uint) (*models.Consultant, error) {
	var consultant models.Consultant
	err := r.db.First(&consultant, id).Error
	return &consultant, err
}