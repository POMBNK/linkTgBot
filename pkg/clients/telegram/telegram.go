package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/POMBNK/linktgBot/pkg/e"
	"github.com/POMBNK/linktgBot/pkg/logging"
)

const (
	errGet     = "Не удалось выполнить запрос:"
	errSendMsg = "Не удалость отправить сообщение:"
)

// host: host api сервиса телеграма. Пример https://api.telegram.org/<basePath>
// basepath: префикс с которого начинаются все запросы tg-bot.com/bot<TOKEN HERE>
// client: обработчик http запросов, используется стандартный http.Client.
type Client struct {
	logger   *logging.Logger
	host     string
	basePath string
	client   *http.Client
}

// New... Конструктор для Client. Возвращает заполненную структуру для инициализации объекта.
// Пример (client:=telegram.New(host,token)).
func New(logger *logging.Logger, host string, token string) *Client {
	return &Client{
		logger:   logger,
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

// newBasePath... возвращает "bot"+token для basePath.
func newBasePath(token string) string {
	return "bot" + token
}

// Updates... возвращает последние обновления.
// offset(смещение) позволяет получать только новые обновления, а не с первого.
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	// c.logger.Debug("Получаем обновления...")
	data, err := c.doRequest(q, "getUpdates")
	if err != nil {
		c.logger.Error("Не удалось получить обновления с сервера Tg.")
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		c.logger.Error("Не удалось заанмаршелить обновления с сервера Tg.")
		return nil, err
	}
	return res.Result, nil

}

// Большое пояснение для самого себя
/*
	В Updates мы создаем query параметры для Get запроса к методу getUpdate, который необходим
	для получения обновлений с сервера telegram. В данном случае эти query параметры - offset и limit,
	они будут добавлены в конец GET запроса.
	В doRequest мы выполняем Get запрос, сначала создавая структуру для URL затем трансформируя ее поля
	в URL методом String. Получится URL в виде https://api.telegram.org/<bot+token>/getUpdates
	По этому URL выполняется запрос getUpdates?query. где query пары ключ-значение.
		https://api.telegram.org/<bot+token>/getUpdates?query
	В полном формате будет выглядеть так:
		https://api.telegram.org/<bot+token>/getUpdates?offset=<value>&limit=<value>
*/

func (c *Client) doRequest(q url.Values, method string) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	// NewRequest(метод запроса, сборка структуры в URL методом String, Тело запроса если есть, иначе nil).
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errGet, err)
	}
	// RawQuery поле структуры Value, содержащее query значения ПОСЛЕ<?> в Get запросе.
	// Encode приводит пары ключ-значения в правильный вид для отправки запроса.
	req.URL.RawQuery = q.Encode()

	// Отправляем получившийся объект запроса.
	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("Не удалось выполнить GET запрос.")
		return nil, e.Wrap(errGet, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Считываем тело ответа на GET запрос.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// SendMessage отправляет сообщение, хранящееся в text в заданный chatID.
// Отправка также происходит через Get запрос методом /sendMessage/<query>.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.doRequest(q, "sendMessage")
	if err != nil {
		c.logger.Error("Не удалось выполнить запрос Send Message...")
		return e.Wrap(errSendMsg, err)
	}
	return nil
}
