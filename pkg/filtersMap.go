package filters

import (
	"strconv"

	"gorm.io/gorm"
)

type FilterFunc func(string, *gorm.DB) *gorm.DB

var QuestionFilters = map[string]FilterFunc{
	"statement": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("statement ILIKE ?", "%"+value+"%")
	},
	"answer": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("id IN (SELECT question_id FROM alternatives WHERE text ILIKE ?)", "%"+value+"%")
	},
	"subject": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.
			Joins("INNER JOIN subject ON subject.id = question.subject_id").
			Where("subject.id = ?", value) // Assumindo que subjects.name armazena "Direito Constitucional"
	},
	"topic": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Where("topic = ?", value)
	},
	"source": func(value string, qb *gorm.DB) *gorm.DB {
		return qb.Joins("JOIN question_source ON question_source.question_id = question.id").
			Where("question_source.source_id = ?", value)
	},

	"year": func(value string, qb *gorm.DB) *gorm.DB {
		if yearInt, err := strconv.Atoi(value); err == nil {
			return qb.
				Joins("LEFT JOIN question_source ON question.id = question_source.question_id").
				Joins("LEFT JOIN source ON source.id = question_source.source_id").
				Where("(source.metadata->>'year')::int = ? OR EXTRACT(YEAR FROM question.created_at) = ?",
					yearInt,
					yearInt)
		}
		return qb
	},
	"list": func(value string, qb *gorm.DB) *gorm.DB {
		if isList, err := strconv.ParseBool(value); err == nil {
			if isList {
				return qb.Where("question_set_id IS NOT NULL")
			}
			return qb.Where("question_set_id IS NULL")
		}
		return qb
	},
}
