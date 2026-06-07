package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


const geminiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

const geminiSystemPrompt = `Ты помощник салона красоты Beauty Dana (г. Есик, Казахстан). Отвечай коротко, на "Вы". Только про красоту, уход, стиль. Максимум 3 предложения.`


type ChatMsgHistory struct {
	Role    string `json:"role"`
	Message string `json:"message"`
}

type ChatRequestGemini struct {
	Message   string           `json:"message" binding:"required"`
	SessionID string           `json:"session_id"`
	UserID    string           `json:"user_id"`
	History   []ChatMsgHistory `json:"history"`
	Session   *BotSession      `json:"session"`
}

type BotSession struct {
	Step        string  `json:"step"`
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	MasterID    string  `json:"master_id"`
	MasterName  string  `json:"master_name"`
	Date        string  `json:"date"`
	SlotStart   string  `json:"slot_start"`
	Price       float64 `json:"price"`
}

type ServiceData struct {
	ID          string
	Name        string
	Category    string
	DurationMin int
	Price       float64
	Description string
}

type MasterData struct {
	ID       string
	Name     string
	Bio      string
	Exp      int
	Services []string
}

type BotData struct {
	Services []ServiceData
	Masters  []MasterData
}


type GeminiMsg struct {
	Role  string          `json:"role"`
	Parts []GeminiMsgPart `json:"parts"`
}

type GeminiMsgPart struct {
	Text string `json:"text"`
}

type GeminiSimpleReq struct {
	Contents          []GeminiMsg    `json:"contents"`
	SystemInstruction *GeminiMsg     `json:"systemInstruction,omitempty"`
	GenerationConfig  *GeminiGenConf `json:"generationConfig,omitempty"`
}

type GeminiGenConf struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type GeminiSimpleResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}


func ChatGemini(pool *pgxpool.Pool, geminiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ChatRequestGemini
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		lower := strings.ToLower(normalizeChatMsg(req.Message))

		sess := req.Session
		if sess == nil {
			sess = &BotSession{Step: "idle"}
		}

		
		data, err := loadData(ctx, pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		
		if sess.Step != "idle" && sess.Step != "" {
			resp, quickReplies, newSess := handleBookingStep(ctx, pool, lower, req.Message, req.UserID, sess, data)
			saveHistorySimple(ctx, pool, req.SessionID, req.UserID, req.Message, resp)
			c.JSON(http.StatusOK, gin.H{
				"message":       resp,
				"quick_replies": quickReplies,
				"session":       newSess,
				"source":        "script",
			})
			return
		}


		if resp, qr, newSess := handleScriptIntent(ctx, pool, lower, req.Message, req.UserID, sess, data); resp != "" {
			saveHistorySimple(ctx, pool, req.SessionID, req.UserID, req.Message, resp)
			c.JSON(http.StatusOK, gin.H{
				"message":       resp,
				"quick_replies": qr,
				"session":       newSess,
				"source":        "script",
			})
			return
		}


		if geminiKey == "" {
			resp := buildDefaultMsg()
			c.JSON(http.StatusOK, gin.H{"message": resp, "source": "script"})
			return
		}

		geminiResp, err := callGeminiSimple(req.Message, geminiKey)
		if err != nil {
			fmt.Printf("[GEMINI ERROR] %v\n", err)
			resp := buildDefaultMsg()
			c.JSON(http.StatusOK, gin.H{"message": resp, "source": "fallback"})
			return
		}

		saveHistorySimple(ctx, pool, req.SessionID, req.UserID, req.Message, geminiResp)
		c.JSON(http.StatusOK, gin.H{
			"message":       geminiResp,
			"quick_replies": []string{"Записаться", "Услуги и цены", "Наши мастера"},
			"session":       sess,
			"source":        "gemini",
		})
	}
}


func handleScriptIntent(ctx context.Context, pool *pgxpool.Pool, lower, original, userID string, sess *BotSession, data *BotData) (string, []string, *BotSession) {

	if hasAny(lower, []string{"привет", "здравствуй", "добрый день", "добрый вечер", "добрый утро", "салам", "сәлем", "hello", "hi"}) {
		name := ""
		if userID != "" {
			pool.QueryRow(ctx, `select full_name from public.profiles where id=$1::uuid`, userID).Scan(&name)
			if name != "" {
				name = ", " + strings.Split(name, " ")[0]
			}
		}
		return fmt.Sprintf("Привет%s! Я помощник Beauty Dana. Могу рассказать об услугах, ценах, мастерах или помочь с записью.", name),
			[]string{"Записаться", "Услуги и цены", "Наши мастера", "Скидки"},
			sess
	}

	if hasAny(lower, []string{"пока", "до свидания", "до встречи", "сау бол", "bye", "goodbye"}) {
		return "До свидания! Будем рады видеть вас в Beauty Dana!", nil, sess
	}

	if hasAny(lower, []string{"спасибо", "благодар", "рахмет", "thanks", "thank you"}) {
		return "Пожалуйста! Если появятся вопросы — обращайтесь.", []string{"Записаться", "Услуги и цены"}, sess
	}

	if hasAny(lower, []string{"кто вы", "что это", "о салоне", "расскажите о себе", "что за салон", "what is"}) {
		return "Мы - **Beauty Dana**, студия красоты в городе Есик \n\nРаботаем с 2012 года. Стрижки, окрашивание, маникюр, педикюр, брови, макияж.\nЗапись онлайн круглосуточно — без звонков и ожидания.",
			[]string{"Услуги и цены", "Наши мастера", "Записаться"},
			sess
	}

	if hasAny(lower, []string{"часы", "режим", "график", "когда работаете", "до скольки", "открыт", "закрыт", "жұмыс"}) {
		return "Работаем **Пн-Сб с 10:00 до 19:00**\nВоскресенье - выходной.\n\nЗапись онлайн доступна круглосуточно!", []string{"Записаться"}, sess
	}

	if hasAny(lower, []string{"адрес", "где находитесь", "где вы", "как доехать", "как добраться", "мекенжай", "location"}) {
		return "**г. Есик, Алматинская область**\n\nТочный адрес и маршрут - на странице «О нас».", []string{"Записаться"}, sess
	}

	if hasAny(lower, []string{"телефон", "номер", "позвонить", "контакт", "связаться", "whatsapp"}) {
		return "Телефон: **+7 777 000 00 00**\n📱 WhatsApp: тот же номер\n\nИли просто запишитесь онлайн — это быстрее!", []string{"Записаться"}, sess
	}

	if hasAny(lower, []string{"скидк", "акци", "бонус", "промокод", "жеңілдік", "discount", "offer"}) {
		return "**Скидки Beauty Dana:**\n\n• Первая запись - **скидка 30%**\n• Несколько услуг — **+5%** за каждую доп. услугу\n• Перенос мастером — скидка 30%\n• Баллы лояльности — 10 баллов = 100 ₸\n• Beauty Streak — скидка до 10% за регулярные визиты",
			[]string{"Записаться", "Программа лояльности"},
			sess
	}

	if hasAny(lower, []string{"лояльност", "балл", "streak", "бонусн", "накопит", "ұпай"}) {
		return "**Beauty Streak — программа лояльности:**\n\n• 10 баллов за каждый завершённый визит\n• 10 баллов = 100 ₸ скидки\n• Уровни: Bronze - Silver - Gold - Platinum\n• Чем регулярнее визиты — тем выше уровень и скидка (до 10%)",
			[]string{"Записаться", "Мои баллы"},
			sess
	}

	if hasAny(lower, []string{"мои балл", "мой баланс", "сколько баллов", "мой streak"}) {
		if userID == "" {
			return "Для просмотра баллов нужно **войти в аккаунт**.", []string{"Войти"}, sess
		}
		var points int
		var level string
		pool.QueryRow(ctx, `select coalesce(points_balance,0) from public.loyalty_accounts where user_id=$1::uuid`, userID).Scan(&points)
		pool.QueryRow(ctx, `select coalesce(current_level,'bronze') from public.beauty_streaks where user_id=$1::uuid`, userID).Scan(&level)
		return fmt.Sprintf("**Ваши баллы:** %d (= %d ₸)\n **Уровень:** %s", points, points*10, strings.Title(level)),
			[]string{"Записаться"},
			sess
	}

	if hasAny(lower, []string{"услуг", "прайс", "что делаете", "чем занимает", "что у вас", "services", "қызмет"}) {
		return buildServicesMsg(data), []string{"Записаться", "Наши мастера"}, sess
	}

	if hasAny(lower, []string{"стрижк", "постричь", "укладк", "окрашивани", "балаяж", "мелирован", "завивк", "шаш"}) {
		return buildCategoryMsg("hair", data), []string{"Записаться на стрижку"}, sess
	}
	if hasAny(lower, []string{"маникюр", "педикюр", "ногт", "покрыти", "гель", "nail"}) {
		return buildCategoryMsg("nails", data), []string{"Записаться на маникюр"}, sess
	}
	if hasAny(lower, []string{"брови", "ресниц", "макияж", "лицо", "визаж", "face"}) {
		return buildCategoryMsg("face", data), []string{"Записаться"}, sess
	}

	if hasAny(lower, []string{"цен", "сколько стоит", "стоимост", "прейскурант", "баға", "price", "cost"}) {
		return buildPricesMsg(data), []string{"Записаться"}, sess
	}

	if hasAny(lower, []string{"мастер", "специалист", "кто работает", "команда", "шебер", "master"}) {
		return buildMastersMsg(data), []string{"Записаться"}, sess
	}

	for _, m := range data.Masters {
		firstName := strings.ToLower(strings.Split(m.Name, " ")[0])
		if strings.Contains(lower, firstName) || strings.Contains(lower, strings.ToLower(m.Name)) {
			return buildMasterDetailMsg(m), []string{fmt.Sprintf("Записаться к %s", strings.Split(m.Name, " ")[0])}, sess
		}
	}

	if hasAny(lower, []string{"записат", "запись", "хочу записат", "забронир", "book", "appointment", "жазыл"}) {
		if userID == "" {
			return "Для записи нужно **войти в аккаунт**", []string{"Войти"}, sess
		}
		newSess := &BotSession{Step: "service"}
		return "Отлично! Выберите услугу:", buildServiceButtons(data), newSess
	}

	if hasAny(lower, []string{"мои записи", "мои визиты", "моя запись", "когда я записан", "мое расписание"}) {
		if userID == "" {
			return "Для просмотра записей нужно **войти в аккаунт**.", []string{"Войти"}, sess
		}
		return buildMyBookingsMsg(ctx, pool, userID), []string{"Записаться", "Отменить запись"}, sess
	}

	if hasAny(lower, []string{"отменить", "отмена записи", "хочу отменить", "cancel"}) {
		if userID == "" {
			return "Для отмены записи нужно **войти в аккаунт**.", []string{"Войти"}, sess
		}
		return "Для отмены записи перейдите в раздел **«Мои записи»** в личном кабинете — там кнопка «Отменить» рядом с каждой активной записью.", []string{"Мои записи"}, sess
	}

	return "", nil, sess
}


func handleBookingStep(ctx context.Context, pool *pgxpool.Pool, lower, original, userID string, sess *BotSession, data *BotData) (string, []string, *BotSession) {

	switch sess.Step {

	case "service":
		for _, svc := range data.Services {
			if strings.Contains(lower, strings.ToLower(svc.Name)) || original == svc.ID {
				sess.ServiceID   = svc.ID
				sess.ServiceName = svc.Name
				sess.Price       = svc.Price
				sess.Step        = "master"
				masterBtns := buildMasterButtons(data)
				masterBtns = append(masterBtns, "Любой мастер")
				return fmt.Sprintf("**%s** выбрана\n\nК какому мастеру хотите записаться?", svc.Name),
					masterBtns, sess
			}
		}
		return "Не нашёл такую услугу. Выберите из списка:", buildServiceButtons(data), sess

	case "master":
		if hasAny(lower, []string{"любой", "не важно", "без разницы", "любого"}) {
			sess.MasterID   = "any"
			sess.MasterName = "Любой мастер"
		} else {
			for _, m := range data.Masters {
				firstName := strings.ToLower(strings.Split(m.Name, " ")[0])
				if strings.Contains(lower, firstName) || strings.Contains(lower, strings.ToLower(m.Name)) {
					sess.MasterID   = m.ID
					sess.MasterName = m.Name
					break
				}
			}
		}
		if sess.MasterID == "" {
			btns := buildMasterButtons(data)
			btns = append(btns, "Любой мастер")
			return "Не нашёл такого мастера. Выберите из списка:", btns, sess
		}
		sess.Step = "date"
		return fmt.Sprintf("Мастер **%s**\n\nНа какую дату? Напишите, например: **сегодня**, **завтра** или **25 мая**", sess.MasterName),
			[]string{"Сегодня", "Завтра", "Послезавтра"},
			sess

	case "date":
		parsed := parseChatDateSimple(original)
		if parsed == "" {
			return "Не понял дату. Попробуйте: **завтра**, **25 мая** или **2026-05-25**",
				[]string{"Сегодня", "Завтра", "Послезавтра"}, sess
		}
		sess.Date = parsed
		sess.Step = "time"
		slots := getSlotsSimple(ctx, pool, sess.MasterID, sess.ServiceID, parsed, data)
		if len(slots) == 0 {
			sess.Step = "date"
			return fmt.Sprintf("На **%s** нет свободного времени \n\nВыберите другую дату:", formatDateRu(parsed)),
				[]string{"Завтра", "Послезавтра"}, sess
		}
		return fmt.Sprintf("Свободное время на **%s**:", formatDateRu(parsed)), slots, sess

	case "time":
		if len(original) == 5 && original[2] == ':' {
			sess.SlotStart = fmt.Sprintf("%sT%s:00+05:00", sess.Date, original)
			sess.Step = "confirm"
			return fmt.Sprintf(
				"Проверьте запись:\n\n **%s**\n **%s**\n **%s в %s**\n **%.0f ₸**\n\nПодтвердить запись?",
				sess.ServiceName, sess.MasterName, formatDateRu(sess.Date), original, sess.Price,
			), []string{"Да, записаться", "Отмена"}, sess
		}
		return "Выберите время из списка выше", nil, sess

	case "confirm":
		if hasAny(lower, []string{"да", "записат", "подтвер", "окей", "ок"}) {
			if userID == "" {
				return "Для записи нужно **войти в аккаунт**.", []string{"Войти"}, &BotSession{Step: "idle"}
			}
			err := createBookingSimple(ctx, pool, userID, sess)
			if err != nil {
				if strings.Contains(err.Error(), "slot_unavailable") {
					sess.Step = "date"
					return "Этот слот только что заняли. Выберите другое время:",
						[]string{"Сегодня", "Завтра"}, sess
				}
				return "Ошибка при создании записи. Попробуйте ещё раз или запишитесь через форму на сайте.",
					nil, &BotSession{Step: "idle"}
			}
			timeStr := ""
			if len(sess.SlotStart) >= 16 {
				timeStr = sess.SlotStart[11:16]
			}
			return fmt.Sprintf("**Запись создана!**\n\n %s\n %s\n %s в %s\n %.0f ₸\n\nДо встречи! ",
				sess.ServiceName, sess.MasterName, formatDateRu(sess.Date), timeStr, sess.Price),
				[]string{"Мои записи", "Записаться ещё"},
				&BotSession{Step: "idle"}
		}
		return "Запись отменена. Чем ещё могу помочь?",
			[]string{"Записаться", "Услуги и цены"},
			&BotSession{Step: "idle"}
	}

	return buildDefaultMsg(), []string{"Записаться", "Услуги и цены"}, &BotSession{Step: "idle"}
}


func callGeminiSimple(message, apiKey string) (string, error) {
	req := GeminiSimpleReq{
		Contents: []GeminiMsg{
			{Role: "user", Parts: []GeminiMsgPart{{Text: message}}},
		},
		SystemInstruction: &GeminiMsg{
			Role:  "user",
			Parts: []GeminiMsgPart{{Text: geminiSystemPrompt}},
		},
		GenerationConfig: &GeminiGenConf{
			Temperature:     0.7,
			MaxOutputTokens: 300,
		},
	}

	body, _ := json.Marshal(req)
	url := geminiURL + "?key=" + apiKey
	httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil { return "", err }
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("gemini error %d: %s", resp.StatusCode, string(respBody))
	}

	var gr GeminiSimpleResp
	if err := json.Unmarshal(respBody, &gr); err != nil { return "", err }
	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return gr.Candidates[0].Content.Parts[0].Text, nil
}


func loadData(ctx context.Context, pool *pgxpool.Pool) (*BotData, error) {
	data := &BotData{}

	rows, err := pool.Query(ctx, `
		select id::text, name, category, duration_min, price, coalesce(description,'')
		from public.services where is_active=true order by sort_order
	`)
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var s ServiceData
		rows.Scan(&s.ID, &s.Name, &s.Category, &s.DurationMin, &s.Price, &s.Description)
		data.Services = append(data.Services, s)
	}

	mrows, err := pool.Query(ctx, `
		select m.id::text, p.full_name, coalesce(m.bio,''), m.experience_years
		from public.masters m join public.profiles p on p.id=m.id
		where m.is_active=true
	`)
	if err != nil { return nil, err }
	defer mrows.Close()
	for mrows.Next() {
		var m MasterData
		mrows.Scan(&m.ID, &m.Name, &m.Bio, &m.Exp)
		srows, _ := pool.Query(ctx, `
			select s.name from public.master_services ms
			join public.services s on s.id=ms.service_id
			where ms.master_id=$1::uuid
		`, m.ID)
		for srows.Next() { var sn string; srows.Scan(&sn); m.Services = append(m.Services, sn) }
		srows.Close()
		data.Masters = append(data.Masters, m)
	}
	return data, nil
}


func buildServicesMsg(data *BotData) string {
	if len(data.Services) == 0 {
		return "Информация об услугах временно недоступна."
	}
	var sb strings.Builder
	sb.WriteString("**Наши услуги:**\n\n")
	for _, s := range data.Services {
		sb.WriteString(fmt.Sprintf("• **%s** — %d мин, **%.0f ₸**\n", s.Name, s.DurationMin, s.Price))
	}
	sb.WriteString("\nЧтобы записаться — нажмите «Записаться» ниже")
	return sb.String()
}

func buildCategoryMsg(cat string, data *BotData) string {
	names := map[string]string{"hair": "Волосы", "nails": "Ногти", "face": "Лицо"}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s — наши услуги:**\n\n", names[cat]))
	found := false
	for _, s := range data.Services {
		if s.Category == cat {
			sb.WriteString(fmt.Sprintf("• **%s** — %d мин, %.0f ₸\n", s.Name, s.DurationMin, s.Price))
			found = true
		}
	}
	if !found { return "По этому направлению пока нет услуг." }
	return sb.String()
}

func buildPricesMsg(data *BotData) string {
	var sb strings.Builder
	sb.WriteString("**Цены на услуги:**\n\n")
	for _, s := range data.Services {
		sb.WriteString(fmt.Sprintf("• %s — **%.0f ₸** (%d мин)\n", s.Name, s.Price, s.DurationMin))
	}
	sb.WriteString("\n Первая запись — скидка **30%**!")
	return sb.String()
}

func buildMastersMsg(data *BotData) string {
	if len(data.Masters) == 0 { return "Информация о мастерах временно недоступна." }
	var sb strings.Builder
	sb.WriteString("**Наши мастера:**\n\n")
	for _, m := range data.Masters {
		sb.WriteString(fmt.Sprintf("**%s** — %d лет опыта\n", m.Name, m.Exp))
		if len(m.Services) > 0 {
			sb.WriteString(fmt.Sprintf("   %s\n", strings.Join(m.Services, ", ")))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func buildMasterDetailMsg(m MasterData) string {
	resp := fmt.Sprintf("**%s** \n%d лет опыта\n\n", m.Name, m.Exp)
	if m.Bio != "" { resp += m.Bio + "\n\n" }
	if len(m.Services) > 0 {
		resp += " Услуги: " + strings.Join(m.Services, ", ")
	}
	return resp
}

func buildMyBookingsMsg(ctx context.Context, pool *pgxpool.Pool, userID string) string {
	rows, err := pool.Query(ctx, `
		select s.name, p.full_name, b.starts_at, b.status, b.price_paid
		from public.bookings b
		join public.services s on s.id=b.service_id
		join public.profiles p on p.id=b.master_id
		where b.client_id=$1::uuid and b.starts_at>=now()-interval '7 days'
		order by b.starts_at desc limit 5
	`, userID)
	if err != nil { return "Не удалось загрузить записи." }
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString("**Ваши записи:**\n\n")
	count := 0
	for rows.Next() {
		var svcName, masterName, status string
		var startsAt time.Time
		var price float64
		rows.Scan(&svcName, &masterName, &startsAt, &status, &price)
		statusLabel := map[string]string{
    "confirmed": "[Подтверждено]",
    "pending":   "[Ожидает]",
    "completed": "[Завершено]",
    "cancelled": "[Отменено]",
	}
	label := statusLabel[status]
	if label == "" { label = "[Запись]" }
	sb.WriteString(fmt.Sprintf("%s **%s**\n   %s · %s · %.0f ₸\n\n",
	    label, svcName, masterName, startsAt.Format("02.01 15:04"), price))
	count++
	}
	if count == 0 { return "У вас пока нет записей. Хотите записаться?" }
	return sb.String()
}

func buildDefaultMsg() string {
	return "Не совсем понял вопрос \n\nМогу помочь с:\n• услугами и ценами\n• записью к мастеру\n• информацией о мастерах\n• скидками и баллами\n• часами работы\n\nСпросите что-нибудь из этого!"
}

func buildServiceButtons(data *BotData) []string {
	var btns []string
	for _, s := range data.Services { btns = append(btns, s.Name) }
	return btns
}

func buildMasterButtons(data *BotData) []string {
	var btns []string
	for _, m := range data.Masters { btns = append(btns, m.Name) }
	return btns
}


func getSlotsSimple(ctx context.Context, pool *pgxpool.Pool, masterID, serviceID, date string, data *BotData) []string {
	if masterID == "any" {
		for _, m := range data.Masters {
			slots := getSlotsMaster(ctx, pool, m.ID, serviceID, date)
			if len(slots) > 0 { return slots }
		}
		return nil
	}
	return getSlotsMaster(ctx, pool, masterID, serviceID, date)
}

func getSlotsMaster(ctx context.Context, pool *pgxpool.Pool, masterID, serviceID, date string) []string {
	var startTime, endTime string
	t, _ := time.Parse("2006-01-02", date)
	dow := int(t.Weekday())

	err := pool.QueryRow(ctx, `
		select start_time::text, end_time::text from public.master_schedules
		where master_id=$1::uuid and day_of_week=$2
	`, masterID, dow).Scan(&startTime, &endTime)
	if err != nil { return nil }

	var duration int
	pool.QueryRow(ctx, `select duration_min from public.services where id=$1::uuid`, serviceID).Scan(&duration)
	if duration == 0 { duration = 60 }

	rows, err := pool.Query(ctx, `
		select starts_at, ends_at from public.bookings
		where master_id=$1::uuid and starts_at::date=$2::date and status not in ('cancelled','no_show')
	`, masterID, date)
	if err != nil { return nil }
	defer rows.Close()

	type slot struct{ s, e time.Time }
	var busy []slot
	for rows.Next() {
		var s, e time.Time; rows.Scan(&s, &e)
		busy = append(busy, slot{s, e})
	}

	sp, _ := time.Parse("15:04:05", startTime)
	ep, _ := time.Parse("15:04:05", endTime)
	base, _ := time.Parse("2006-01-02", date)
	now := time.Now().Add(30 * time.Minute)

	cur := base.Add(time.Duration(sp.Hour())*time.Hour + time.Duration(sp.Minute())*time.Minute)
	end := base.Add(time.Duration(ep.Hour())*time.Hour + time.Duration(ep.Minute())*time.Minute)

	var result []string
	for !cur.Add(time.Duration(duration)*time.Minute).After(end) {
		slotEnd := cur.Add(time.Duration(duration) * time.Minute)
		if cur.After(now) {
			overlap := false
			for _, b := range busy {
				if cur.Before(b.e) && slotEnd.After(b.s) { overlap = true; break }
			}
			if !overlap { result = append(result, cur.Format("15:04")) }
		}
		cur = cur.Add(30 * time.Minute)
	}
	return result
}


func createBookingSimple(ctx context.Context, pool *pgxpool.Pool, clientID string, sess *BotSession) error {
	startsAt, err := time.Parse("2006-01-02T15:04:05-07:00", sess.SlotStart)
	if err != nil {
		startsAt, err = time.Parse("2006-01-02T15:04:05Z", sess.SlotStart)
		if err != nil { return fmt.Errorf("invalid time: %w", err) }
	}

	var duration int
	pool.QueryRow(ctx, `select duration_min from public.services where id=$1::uuid`, sess.ServiceID).Scan(&duration)
	if duration == 0 { duration = 60 }
	endsAt := startsAt.Add(time.Duration(duration) * time.Minute)

	masterID := sess.MasterID
	if masterID == "any" {
		rows, _ := pool.Query(ctx, `select id::text from public.masters where is_active=true`)
		var mids []string
		for rows.Next() { var id string; rows.Scan(&id); mids = append(mids, id) }
		rows.Close()
		for _, mid := range mids {
			var cnt int
			pool.QueryRow(ctx, `
				select count(*) from public.bookings
				where master_id=$1::uuid and status not in ('cancelled','no_show')
				and starts_at < $3 and ends_at > $2
			`, mid, startsAt, endsAt).Scan(&cnt)
			if cnt == 0 { masterID = mid; break }
		}
		if masterID == "any" { return fmt.Errorf("slot_unavailable") }
	}

	var cnt int
	pool.QueryRow(ctx, `
		select count(*) from public.bookings
		where master_id=$1::uuid and status not in ('cancelled','no_show')
		and starts_at < $3 and ends_at > $2
	`, masterID, startsAt, endsAt).Scan(&cnt)
	if cnt > 0 { return fmt.Errorf("slot_unavailable") }

	price := sess.Price
	var bookCount int
	pool.QueryRow(ctx, `select count(*) from public.bookings where client_id=$1::uuid and status not in ('cancelled')`, clientID).Scan(&bookCount)
	if bookCount == 0 { price = price * 0.70 }

	_, err = pool.Exec(ctx, `
		insert into public.bookings
		  (client_id, master_id, service_id, starts_at, ends_at, status, price_paid, notes)
		values ($1::uuid, $2::uuid, $3::uuid, $4, $5, 'confirmed', $6, 'Запись через чат-бот')
	`, clientID, masterID, sess.ServiceID, startsAt, endsAt, price)
	return err
}


func parseChatDateSimple(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	now   := time.Now()
	switch {
	case hasAny(lower, []string{"сегодня", "today", "бүгін"}):
		return now.Format("2006-01-02")
	case hasAny(lower, []string{"завтра", "tomorrow", "ертең"}):
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	case hasAny(lower, []string{"послезавтра"}):
		return now.AddDate(0, 0, 2).Format("2006-01-02")
	}
	if len(lower) == 10 && lower[4] == '-' { return lower }
	months := map[string]int{
		"января":1,"февраля":2,"марта":3,"апреля":4,"мая":5,"июня":6,
		"июля":7,"августа":8,"сентября":9,"октября":10,"ноября":11,"декабря":12,
	}
	parts := strings.Fields(lower)
	if len(parts) >= 2 {
		var day int
		fmt.Sscanf(parts[0], "%d", &day)
		if m, ok := months[parts[1]]; ok && day > 0 {
			d := time.Date(now.Year(), time.Month(m), day, 0, 0, 0, 0, time.Local)
			if d.Before(now) { d = d.AddDate(1, 0, 0) }
			return d.Format("2006-01-02")
		}
	}
	return ""
}

func formatDateRu(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil { return dateStr }
	months := []string{"","января","февраля","марта","апреля","мая","июня","июля","августа","сентября","октября","ноября","декабря"}
	days   := []string{"воскресенье","понедельник","вторник","среда","четверг","пятница","суббота"}
	return fmt.Sprintf("%s, %d %s", days[t.Weekday()], t.Day(), months[t.Month()])
}

func saveHistorySimple(ctx context.Context, pool *pgxpool.Pool, sessionID, userID, userMsg, botMsg string) {
	if sessionID == "" { return }
	var uid interface{} = nil
	if userID != "" { uid = userID }
	pool.Exec(ctx, `
		insert into public.chat_history (session_id, user_id, role, message)
		values ($1,$2,'user',$3),($1,$2,'bot',$4)
	`, sessionID, uid, userMsg, botMsg)
}

func hasAny(s string, kws []string) bool {
	for _, k := range kws {
		if strings.Contains(s, k) { return true }
	}
	return false
}

func normalizeChatMsg(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(' ')
		}
	}
	return b.String()
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
