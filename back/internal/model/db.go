package model

import "gorm.io/gorm"

type MailDB interface {
	Model(value interface{}) (tx MailDB)
	Select(query interface{}, args ...interface{}) (tx MailDB)
	Create(value interface{}) (tx MailDB)
	Update(column string, value interface{}) (tx MailDB)
	Delete(value interface{}, conds ...interface{}) (tx MailDB)
	Where(query interface{}, args ...interface{}) (tx MailDB)
	Find(dest interface{}, conds ...interface{}) (tx MailDB)
	First(dest interface{}, conds ...interface{}) (tx MailDB)
	Error() error
}

type mailDB struct {
	*gorm.DB
}

func NewMailDB(db *gorm.DB) MailDB {
	return &mailDB{DB: db}
}

func (m *mailDB) Model(value interface{}) (tx MailDB) {
	return &mailDB{m.DB.Model(value)}
}

func (m *mailDB) Select(query interface{}, args ...interface{}) (tx MailDB) {
	return &mailDB{m.DB.Select(query, args...)}
}

func (m *mailDB) Create(value interface{}) (tx MailDB) {
	return &mailDB{m.DB.Create(value)}
}

func (m *mailDB) Update(column string, value interface{}) (tx MailDB) {
	return &mailDB{m.DB.Update(column, value)}
}

func (m *mailDB) Delete(value interface{}, conds ...interface{}) (tx MailDB) {
	return &mailDB{m.DB.Delete(value, conds...)}
}

func (m *mailDB) Where(query interface{}, args ...interface{}) (tx MailDB) {
	return &mailDB{m.DB.Where(query, args...)}
}

func (m *mailDB) Find(dest interface{}, conds ...interface{}) (tx MailDB) {
	return &mailDB{m.DB.Find(dest, conds...)}
}

func (m *mailDB) First(dest interface{}, conds ...interface{}) (tx MailDB) {
	return &mailDB{m.DB.First(dest, conds...)}
}

func (m *mailDB) Error() error {
	return m.DB.Error
}
