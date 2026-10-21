package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

// TitleMaxLength es el número máximo de caracteres permitido en el título.
const TitleMaxLength = 200

// Task representa una tarea de la lista.
type Task struct {
	ID        uuid.UUID  `json:"id" db:"id" form:"-"`
	Title     string     `json:"title" db:"title" form:"title"`
	Completed bool       `json:"completed" db:"completed" form:"-"`
	DueDate   *time.Time `json:"due_date" db:"due_date" form:"-"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" form:"-"`
}

// DueDateLayout es el formato del campo type="date" y del atributo datetime.
const DueDateLayout = "2006-01-02"

// ErrDueDateInvalid indica que la fecha límite escrita no es una fecha válida.
var ErrDueDateInvalid = errors.New("La fecha límite no es válida.")

var monthLabels = [12]string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sept", "oct", "nov", "dic"}

// ParseDueDate interpreta la fecha límite del formulario; vacía significa sin fecha.
func ParseDueDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	date, err := time.Parse(DueDateLayout, value)
	if err != nil {
		return nil, ErrDueDateInvalid
	}
	return &date, nil
}

// DueDateValue devuelve la fecha límite como AAAA-MM-DD, o vacío si no tiene.
func (t Task) DueDateValue() string {
	if t.DueDate == nil {
		return ""
	}
	return t.DueDate.UTC().Format(DueDateLayout)
}

// DueDateLabel devuelve la fecha límite en español corto, como «25 oct 2026».
func (t Task) DueDateLabel() string {
	if t.DueDate == nil {
		return ""
	}
	date := t.DueDate.UTC()
	return fmt.Sprintf("%d %s %d", date.Day(), monthLabels[date.Month()-1], date.Year())
}

// Tasks es una colección de Task.
type Tasks []Task

// Validate se ejecuta en cada llamada a los métodos "pop.Validate*".
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	title := strings.TrimSpace(t.Title)
	if title == "" {
		verrs.Add("title", "El título es obligatorio.")
	} else if len([]rune(title)) > TitleMaxLength {
		verrs.Add("title", fmt.Sprintf("El título no puede superar los %d caracteres.", TitleMaxLength))
	}

	return verrs, nil
}
