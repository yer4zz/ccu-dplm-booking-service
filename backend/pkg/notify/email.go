package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"booking-service/internal/model"
)

type EmailNotifier struct {
	apiKey string
	from   string
}

func NewEmailNotifier(apiKey, from string) *EmailNotifier {
	if apiKey == "" {
		fmt.Println("[notify] ⚠️  RESEND_API_KEY не задан")
	}
	return &EmailNotifier{apiKey: apiKey, from: from}
}

type emailPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

func (n *EmailNotifier) send(ctx context.Context, p emailPayload) error {
	if n.apiKey == "" {
		fmt.Printf("[notify] пропуск (нет ключа) - %v\n", p.To)
		return nil
	}
	body, _ := json.Marshal(p)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend %d: %s", resp.StatusCode, string(respBody))
	}
	fmt.Printf("[notify] - %v (status %d)\n", p.To, resp.StatusCode)
	return nil
}

func (n *EmailNotifier) BookingCreated(ctx context.Context, b *model.BookingNotification) {
	if err := n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{b.ClientEmail},
		Subject: "✂️ Запись подтверждена — Beauty Dana",
		Html:    bookingCreatedHTML(b),
	}); err != nil {
		fmt.Printf("[notify] BookingCreated error: %v\n", err)
	}
}

func (n *EmailNotifier) BookingCancelled(ctx context.Context, b *model.BookingNotification) {
	if err := n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{b.ClientEmail},
		Subject: "Запись отменена — Beauty Dana",
		Html:    bookingCancelledHTML(b),
	}); err != nil {
		fmt.Printf("[notify] BookingCancelled error: %v\n", err)
	}
}

func (n *EmailNotifier) BookingRescheduled(ctx context.Context, b *model.BookingNotification) {
	if err := n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{b.ClientEmail},
		Subject: "📅 Запись перенесена — Beauty Dana",
		Html:    bookingRescheduledHTML(b),
	}); err != nil {
		fmt.Printf("[notify] BookingRescheduled error: %v\n", err)
	}
}

func (n *EmailNotifier) BookingAutoRescheduled(ctx context.Context, b *model.BookingNotification) {
	if err := n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{b.ClientEmail},
		Subject: "🔄 Ваша запись перенесена к другому мастеру — Beauty Dana",
		Html:    bookingAutoRescheduledHTML(b),
	}); err != nil {
		fmt.Printf("[notify] BookingAutoRescheduled error: %v\n", err)
	}
}

func (n *EmailNotifier) SendUpcomingReminder(ctx context.Context, d *ReminderData) error {
	return n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{d.ClientEmail},
		Subject: "⏰ Напоминание о записи завтра — Beauty Dana",
		Html:    upcomingReminderHTML(d),
	})
}

func (n *EmailNotifier) SendPersonalizedReminder(ctx context.Context, d *PersonalizedData) error {
	return n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{d.ClientEmail},
		Subject: "💇 Не пора ли снова? — Beauty Dana",
		Html:    personalizedReminderHTML(d),
	})
}

func (n *EmailNotifier) SendSOSToMaster(ctx context.Context, d *SOSNotification) error {
	return n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{d.MasterEmail},
		Subject: fmt.Sprintf("🆘 SOS запрос от %s — %.0f ₸", d.ClientName, d.SOSPrice),
		Html:    sosToMasterHTML(d),
	})
}

func (n *EmailNotifier) SendSOSResponse(ctx context.Context, d *SOSResponse) error {
	subject := "✅ Мастер принял ваш SOS запрос"
	if !d.Accepted {
		subject = "К сожалению, мастер не смог принять SOS"
	}
	return n.send(ctx, emailPayload{
		From:    n.from,
		To:      []string{d.ClientEmail},
		Subject: subject,
		Html:    sosResponseHTML(d),
	})
}

func bookingCreatedHTML(b *model.BookingNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,.08)}
.h{background:#1A1714;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0;font-weight:600}
.b{padding:32px 40px}.r{display:flex;justify-content:space-between;padding:12px 0;border-bottom:1px solid #F4F2EF;font-size:14px}
.l{color:#6B6560}.v{font-weight:500}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
.btn{display:inline-block;background:#C9956A;color:#fff;padding:14px 32px;border-radius:99px;text-decoration:none;font-weight:600;font-size:15px;margin-top:20px}
</style></head><body><div class="w">
<div class="h"><h1>Запись подтверждена</h1></div>
<div class="b">
<p style="color:#6B6560;font-size:15px;margin:0 0 20px">Здравствуйте, <b style="color:#1A1714">%s</b>!</p>
<div class="r"><span class="l">Услуга</span><span class="v">%s</span></div>
<div class="r"><span class="l">Мастер</span><span class="v">%s</span></div>
<div class="r"><span class="l">Дата и время</span><span class="v">%s</span></div>
<div class="r" style="border:none"><span class="l">Стоимость</span><span class="v">%.0f ₸</span></div>
<div style="text-align:center"><a class="btn" href="#">Посмотреть запись</a></div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		b.ClientName, b.ServiceName, b.MasterName,
		b.StartsAt.Format("02.01.2006 в 15:04"), b.Price)
}

func bookingCancelledHTML(b *model.BookingNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:#EF4444;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
.btn{display:inline-block;background:#1A1714;color:#fff;padding:14px 32px;border-radius:99px;text-decoration:none;font-weight:600;font-size:15px;margin-top:20px}
</style></head><body><div class="w">
<div class="h"><h1>Запись отменена</h1></div>
<div class="b">
<p style="color:#6B6560;font-size:15px">Здравствуйте, <b style="color:#1A1714">%s</b>!<br><br>
Ваша запись на <b>%s</b> (%s) была отменена.<br>
Вы можете записаться снова в удобное время.</p>
<div style="text-align:center"><a class="btn" href="#">Записаться снова</a></div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		b.ClientName, b.ServiceName,
		b.StartsAt.Format("02.01.2006 в 15:04"))
}

func bookingRescheduledHTML(b *model.BookingNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:#3B82F6;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.r{display:flex;justify-content:space-between;padding:12px 0;border-bottom:1px solid #F4F2EF;font-size:14px}
.l{color:#6B6560}.v{font-weight:500}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
.disc{background:#EAF5EF;border-radius:10px;padding:12px 16px;margin-top:16px;color:#3D8C5F;font-weight:600;font-size:14px}
</style></head><body><div class="w">
<div class="h"><h1>Запись перенесена</h1></div>
<div class="b">
<p style="color:#6B6560;font-size:15px;margin:0 0 20px">Здравствуйте, <b style="color:#1A1714">%s</b>!</p>
<div class="r"><span class="l">Услуга</span><span class="v">%s</span></div>
<div class="r"><span class="l">Мастер</span><span class="v">%s</span></div>
<div class="r"><span class="l">Новая дата</span><span class="v">%s</span></div>
<div class="r" style="border:none"><span class="l">Стоимость</span><span class="v">%.0f ₸</span></div>
<div class="disc">Скидка 30%% применена за перенос</div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		b.ClientName, b.ServiceName, b.MasterName,
		b.StartsAt.Format("02.01.2006 в 15:04"), b.Price)
}

func bookingAutoRescheduledHTML(b *model.BookingNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,.08)}
.h{background:#3B82F6;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0;font-weight:600}
.h p{color:rgba(255,255,255,.85);font-size:14px;margin:8px 0 0}
.b{padding:32px 40px}
.info{background:#EFF6FF;border-radius:12px;padding:20px;margin:20px 0}
.r{display:flex;justify-content:space-between;padding:10px 0;border-bottom:1px solid #DBEAFE;font-size:14px}
.r:last-child{border:none}.l{color:#6B6560}.v{font-weight:500;color:#1A1714}
.warn{background:#FFF7ED;border-left:3px solid #F59E0B;padding:12px 16px;border-radius:0 8px 8px 0;font-size:13px;color:#92400E;margin-top:16px}
.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
</style></head><body><div class="w">
<div class="h">
  <h1>🔄 Запись перенесена автоматически</h1>
  <p>Мы нашли для вас нового мастера</p>
</div>
<div class="b">
<p style="font-size:15px;color:#6B6560;margin:0 0 4px">Здравствуйте, <b style="color:#1A1714">%s</b>!</p>
<p style="font-size:14px;color:#6B6560;margin:0 0 4px">Ваш прежний мастер не смог вас принять. Мы автоматически нашли замену.</p>
<div class="info">
  <div class="r"><span class="l">Услуга</span><span class="v">%s</span></div>
  <div class="r"><span class="l">Новый мастер</span><span class="v">%s</span></div>
  <div class="r"><span class="l">Дата и время</span><span class="v">%s</span></div>
  <div class="r"><span class="l">Стоимость</span><span class="v">%.0f ₸</span></div>
</div>
<div class="warn">Если вас не устраивает новый мастер — вы можете отменить запись в личном кабинете.</div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		b.ClientName, b.ServiceName, b.MasterName,
		b.StartsAt.Format("02.01.2006 в 15:04"), b.Price)
}


type ReminderData struct {
	ClientName  string
	ClientEmail string
	MasterName  string
	ServiceName string
	StartsAt    time.Time
	DurationMin int
}

type PersonalizedData struct {
	ClientName      string
	ClientEmail     string
	AvgIntervalDays int
	LastVisit       time.Time
}

type SOSNotification struct {
	MasterEmail string
	MasterName  string
	ClientName  string
	SOSPrice    float64
	BasePrice   float64
	PrefStart   time.Time
	PrefEnd     time.Time
	ClientNote  string
	SOSID       string
}

type SOSResponse struct {
	ClientEmail string
	ClientName  string
	Accepted    bool
	MasterNote  string
	SOSPrice    float64
	PrefStart   time.Time
}

func upcomingReminderHTML(d *ReminderData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:#C9956A;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.r{display:flex;justify-content:space-between;padding:12px 0;border-bottom:1px solid #F4F2EF;font-size:14px}
.l{color:#6B6560}.v{font-weight:500}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
</style></head><body><div class="w">
<div class="h"><h1>⏰ Напоминание о записи</h1></div>
<div class="b">
<p style="color:#6B6560;font-size:15px;margin:0 0 20px">Здравствуйте, <b style="color:#1A1714">%s</b>! Ждём вас завтра.</p>
<div class="r"><span class="l">Услуга</span><span class="v">%s</span></div>
<div class="r"><span class="l">Мастер</span><span class="v">%s</span></div>
<div class="r"><span class="l">Дата и время</span><span class="v">%s</span></div>
<div class="r" style="border:none"><span class="l">Длительность</span><span class="v">%d мин</span></div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		d.ClientName, d.ServiceName, d.MasterName,
		d.StartsAt.Format("02.01.2006 в 15:04"), d.DurationMin)
}

func personalizedReminderHTML(d *PersonalizedData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:#1A1714;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.btn{display:inline-block;background:#C9956A;color:#fff;padding:14px 32px;border-radius:99px;text-decoration:none;font-weight:600;font-size:15px;margin-top:20px}
.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
</style></head><body><div class="w">
<div class="h"><h1>Не пора ли снова?</h1></div>
<div class="b">
<p style="font-size:15px;color:#6B6560">Здравствуйте, <b style="color:#1A1714">%s</b>!<br><br>
Вы обычно записываетесь каждые <b style="color:#1A1714">%d дней</b>. Похоже, пора снова!</p>
<div style="text-align:center"><a class="btn" href="#">Записаться онлайн</a></div>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		d.ClientName, d.AvgIntervalDays)
}

func sosToMasterHTML(d *SOSNotification) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:#EF4444;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.price{font-size:32px;font-weight:700;color:#EF4444;margin:16px 0}
.r{display:flex;justify-content:space-between;padding:10px 0;border-bottom:1px solid #F4F2EF;font-size:14px}
.l{color:#6B6560}.v{font-weight:500}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
</style></head><body><div class="w">
<div class="h"><h1>🆘 SOS Запрос!</h1></div>
<div class="b">
<p style="font-size:15px;color:#6B6560"><b style="color:#1A1714">%s</b> просит срочный приём.</p>
<div class="price">%.0f ₸ <span style="font-size:14px;color:#9E9890;font-weight:400">(+30%% за срочность)</span></div>
<div class="r"><span class="l">Желаемое время от</span><span class="v">%s</span></div>
<div class="r" style="border:none"><span class="l">до</span><span class="v">%s</span></div>
%s
<p style="margin-top:20px;font-size:14px;color:#6B6560">Войдите в кабинет мастера чтобы принять или отклонить.</p>
</div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		d.ClientName, d.SOSPrice,
		d.PrefStart.Format("02.01 в 15:04"),
		d.PrefEnd.Format("02.01 в 15:04"),
		func() string {
			if d.ClientNote == "" { return "" }
			return fmt.Sprintf(`<div style="background:#FFF7ED;border-left:3px solid #F59E0B;padding:12px;border-radius:0 8px 8px 0;margin-top:12px;font-size:13px;color:#92400E">💬 %s</div>`, d.ClientNote)
		}())
}

func sosResponseHTML(d *SOSResponse) string {
	bg, title, body := "#16A34A", "Запрос принят!", fmt.Sprintf("Мастер принял ваш SOS запрос. Стоимость: <b>%.0f ₸</b>.", d.SOSPrice)
	if !d.Accepted {
		bg, title, body = "#6B7280", "Запрос отклонён", "К сожалению, мастер не смог найти время."
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"><style>
body{font-family:-apple-system,sans-serif;background:#FAF9F7;margin:0;padding:40px 20px}
.w{max-width:560px;margin:0 auto;background:#fff;border-radius:16px;overflow:hidden}
.h{background:%s;padding:32px 40px}.h h1{color:#fff;font-size:22px;margin:0}
.b{padding:32px 40px}.f{padding:20px 40px;background:#F4F2EF;text-align:center;font-size:13px;color:#9E9890}
</style></head><body><div class="w">
<div class="h"><h1>%s</h1></div>
<div class="b"><p style="font-size:15px;color:#6B6560">Здравствуйте, <b style="color:#1A1714">%s</b>!<br><br>%s</p></div>
<div class="f">Beauty Dana</div>
</div></body></html>`,
		bg, title, d.ClientName, body)
}