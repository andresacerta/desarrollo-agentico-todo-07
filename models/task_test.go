package models

import (
	"strings"
	"time"
)

func (ms *ModelSuite) Test_Task_Create_Valid() {
	task := &Task{Title: "Comprar café"}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	ms.NotZero(task.ID)
	ms.False(task.Completed)

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Task_Title_Required() {
	task := &Task{Title: ""}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")

	count, err := ms.DB.Count(&Task{})
	ms.NoError(err)
	ms.Equal(0, count)
}

func (ms *ModelSuite) Test_Task_Title_OnlySpaces_Invalid() {
	task := &Task{Title: "   \t  "}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.Contains(verrs.Get("title"), "El título es obligatorio.")
}

func (ms *ModelSuite) Test_Task_Title_MaxLength() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())
}

func (ms *ModelSuite) Test_Task_Title_TooLong_Invalid() {
	task := &Task{Title: strings.Repeat("a", TitleMaxLength+1)}

	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.True(verrs.HasAny())
	ms.NotEmpty(verrs.Get("title"))
}

func (ms *ModelSuite) Test_Task_Toggle_Persists() {
	task := &Task{Title: "Revisar el correo"}
	verrs, err := ms.DB.ValidateAndCreate(task)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	task.Completed = true
	ms.NoError(ms.DB.Update(task))

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, task.ID))
	ms.True(reloaded.Completed)
}

func (ms *ModelSuite) Test_Task_ParseDueDate_Empty() {
	for _, value := range []string{"", "   "} {
		date, err := ParseDueDate(value)
		ms.NoError(err)
		ms.Nil(date)
	}
}

func (ms *ModelSuite) Test_Task_ParseDueDate_Valid() {
	date, err := ParseDueDate("2026-10-25")
	ms.NoError(err)
	ms.Require().NotNil(date)
	ms.Equal(2026, date.Year())
	ms.Equal(time.October, date.Month())
	ms.Equal(25, date.Day())
}

func (ms *ModelSuite) Test_Task_ParseDueDate_Invalid() {
	for _, value := range []string{"2026-02-30", "25/10/2026", "mañana"} {
		date, err := ParseDueDate(value)
		ms.ErrorIs(err, ErrDueDateInvalid, value)
		ms.Nil(date)
	}
}

func (ms *ModelSuite) Test_Task_DueDateLabel() {
	cases := map[string]string{
		"2026-10-25": "25 oct 2026",
		"2027-01-03": "3 ene 2027",
		"2026-09-01": "1 sept 2026",
	}
	for value, label := range cases {
		date, err := ParseDueDate(value)
		ms.NoError(err)
		ms.Equal(label, Task{DueDate: date}.DueDateLabel())
	}
	ms.Equal("", Task{}.DueDateLabel())
}

func (ms *ModelSuite) Test_Task_DueDate_Persists() {
	date, err := ParseDueDate("2026-10-25")
	ms.NoError(err)

	withDate := &Task{Title: "Con fecha", DueDate: date}
	verrs, err := ms.DB.ValidateAndCreate(withDate)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	withoutDate := &Task{Title: "Sin fecha"}
	verrs, err = ms.DB.ValidateAndCreate(withoutDate)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	reloaded := &Task{}
	ms.NoError(ms.DB.Find(reloaded, withDate.ID))
	ms.Equal("2026-10-25", reloaded.DueDateValue())

	reloaded = &Task{}
	ms.NoError(ms.DB.Find(reloaded, withoutDate.ID))
	ms.Nil(reloaded.DueDate)
}
