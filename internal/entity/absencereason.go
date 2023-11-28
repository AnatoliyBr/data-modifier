package entity

// AbsenceReason wraps info about user absence reason.
type AbsenceReason struct {
	Description string
	Emoji       string
	ID          int
}

// ReasonList stores AbsenceReason in-memory.
var ReasonList = map[int]*AbsenceReason{
	1: {
		ID:          1,
		Description: "Личные дела",
		Emoji:       "🏠",
	},
	2: {
		ID:          2,
		Description: "Гостевой пропуск",
		Emoji:       "✈",
	},
	3: {
		ID:          3,
		Description: "Командировка",
		Emoji:       "✈",
	},
	4: {
		ID:          4,
		Description: "Местная командировка",
		Emoji:       "🚙",
	},
	5: {
		ID:          5,
		Description: "Болезнь",
		Emoji:       "🌡",
	},
	6: {
		ID:          6,
		Description: "Больничный лист",
		Emoji:       "🏥",
	},
	7: {
		ID:          7,
		Description: "Ночные работы",
		Emoji:       "🌃",
	},
	8: {
		ID:          8,
		Description: "Дежурство",
		Emoji:       "👨‍💻",
	},
	9: {
		ID:          9,
		Description: "Учеба",
		Emoji:       "📚",
	},
	10: {
		ID:          10,
		Description: "Удаленная работа",
		Emoji:       "🏠",
	},
	11: {
		ID:          11,
		Description: "Отпуск",
		Emoji:       "🌅",
	},
	12: {
		ID:          12,
		Description: "Отпуск за свой счет",
		Emoji:       "⛺",
	},
	13: {
		ID:          13,
		Description: "Отсутствие с отработкой",
		Emoji:       "⌛",
	},
}
