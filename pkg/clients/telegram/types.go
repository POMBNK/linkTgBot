package telegram

// Ответ с апи приходит в JSON формате. Если поле ok == True, то ответ содержится в поле result.
// Поле result - JSON сериализованные Update объекты. Result: []Update.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update - объект приходящий в результате успешного запроса. Содержит 2 поля ID обновления и тело Message.
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
