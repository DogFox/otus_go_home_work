package storage

type Event struct {
	ID    string // уникальный идентификатор события (можно воспользоваться UUID)
	Title string // Заголовок - короткий текст
	// Date        time.Time     // Дата и время события
	// Duration    time.Duration // Длительность события (или дата и время окончания);
	// Description string        // Описание события - длинный текст, опционально;
	// User_ID     string        // ID пользователя, владельца события;
	// TimeShift   string        //За сколько времени высылать уведомление, опционально.
}
