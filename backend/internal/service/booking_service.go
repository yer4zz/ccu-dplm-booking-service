package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"booking-service/internal/model"
	"booking-service/internal/repository"
	"booking-service/pkg/notify"
)

var (
	ErrSlotUnavailable = errors.New("slot_unavailable")
	ErrServiceMismatch = errors.New("service_not_provided_by_master")
	ErrForbidden       = errors.New("forbidden")
	ErrAlreadyDone     = errors.New("booking_already_completed")
	ErrNotEnoughPoints = errors.New("not_enough_points")
)

type BookingService struct {
	Repo     *repository.Repo
	Notifier *notify.EmailNotifier
}

func NewBookingService(repo *repository.Repo, n *notify.EmailNotifier) *BookingService {
	return &BookingService{Repo: repo, Notifier: n}
}

func (s *BookingService) CalcPrice(
	ctx context.Context,
	clientID uuid.UUID,
	masterID uuid.UUID,
	serviceIDs []uuid.UUID,
	pointsToUse int,
) (*model.BookingPriceCalc, error) {
	calc := &model.BookingPriceCalc{ServiceCount: len(serviceIDs)}

	ids := make([]string, len(serviceIDs))
	for i, id := range serviceIDs {
		ids[i] = id.String()
	}

	rows, err := s.Repo.DB().Query(ctx, `
		SELECT
			s.id::text,
			COALESCE(ms.custom_price, s.price) as price
		FROM public.services s
		LEFT JOIN public.master_services ms
			ON ms.service_id = s.id AND ms.master_id = $1
		WHERE s.id = ANY($2::uuid[])
	`, masterID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var price float64
		rows.Scan(&id, &price)
		calc.BasePrice += price
	}

	isFirst, _ := s.Repo.IsFirstBooking(ctx, clientID)
	if isFirst {
		calc.IsFirstBooking = true
		calc.FirstDiscount  = calc.BasePrice * 0.30
	}
	if len(serviceIDs) > 1 {
		pct := float64(len(serviceIDs)-1) * 5.0
		if pct > 25 { pct = 25 }
		calc.MultiDiscount = calc.BasePrice * pct / 100
	}

	loyalty, _ := s.Repo.GetLoyaltyAccount(ctx, clientID)
	calc.PointsAvailable = loyalty.Balance
	if pointsToUse > 0 && pointsToUse <= loyalty.Balance {
		calc.PointsUsed    = pointsToUse
		calc.PointsDiscount = float64(pointsToUse) * 10
	}

	calc.TotalDiscount = calc.FirstDiscount + calc.MultiDiscount + calc.PointsDiscount
	calc.FinalPrice    = calc.BasePrice - calc.TotalDiscount
	if calc.FinalPrice < 0 {
		calc.FinalPrice = 0
	}
	return calc, nil
}

func (s *BookingService) Create(
	ctx context.Context,
	clientID uuid.UUID,
	req model.CreateBookingRequest,
) (*model.Booking, error) {

	if len(req.ServiceIDs) == 0 {
		return nil, ErrServiceMismatch
	}

	var totalDuration int
	for _, svcID := range req.ServiceIDs {
		_, duration, err := s.Repo.GetServicePrice(ctx, req.MasterID, svcID)
		if err != nil {
			fmt.Printf("[create] услуга %s не найдена: %v\n", svcID, err)
			return nil, ErrServiceMismatch
		}
		totalDuration += duration
	}

	endsAt := req.StartsAt.Add(time.Duration(totalDuration) * time.Minute)

	available, err := s.Repo.IsSlotAvailable(ctx, req.MasterID, req.StartsAt, endsAt)
	if err != nil || !available {
		return nil, ErrSlotUnavailable
	}

	calc, err := s.CalcPrice(ctx, clientID, req.MasterID, req.ServiceIDs, req.PointsUsed)
	if err != nil {
		return nil, err
	}

	booking, err := s.Repo.CreateBooking(ctx, model.Booking{
		ClientID:  clientID,
		MasterID:  req.MasterID,
		ServiceID: req.ServiceIDs[0],
		StartsAt:  req.StartsAt,
		EndsAt:    endsAt,
		Status:    model.StatusPending,
		PricePaid: calc.FinalPrice,
		Notes:     req.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: db conflict", ErrSlotUnavailable)
	}

	for _, svcID := range req.ServiceIDs {
		price, _, _ := s.Repo.GetServicePrice(ctx, req.MasterID, svcID)
		s.Repo.AddBookingService(ctx, booking.ID, svcID, price)
	}

	if calc.TotalDiscount > 0 {
		discType := "multi_service"
		if calc.IsFirstBooking {
			discType = "first_booking"
		}
		if calc.PointsUsed > 0 {
			discType = "loyalty_points"
		}
		s.Repo.SaveDiscount(ctx, model.Discount{
			BookingID:  booking.ID,
			Type:       discType,
			Percent:    (calc.TotalDiscount / calc.BasePrice) * 100,
			PointsUsed: calc.PointsUsed,
			Amount:     calc.TotalDiscount,
		})
		if calc.PointsUsed > 0 {
			s.Repo.DeductPoints(ctx, clientID, booking.ID, calc.PointsUsed)
		}
	}

	go func() {
		notif, err := s.Repo.GetBookingNotification(context.Background(), booking.ID)
		if err != nil {
			fmt.Printf("[booking] уведомление: %v\n", err)
			return
		}
		s.Notifier.BookingCreated(context.Background(), notif)
	}()

	return booking, nil
}

func (s *BookingService) Cancel(
	ctx context.Context,
	bookingID   uuid.UUID,
	requesterID uuid.UUID,
) error {
	b, err := s.Repo.GetBooking(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.ClientID != requesterID && b.MasterID != requesterID {
		return ErrForbidden
	}
	if b.Status == model.StatusCompleted {
		return ErrAlreadyDone
	}

	notif, _ := s.Repo.GetBookingNotification(ctx, bookingID)

	if err := s.Repo.UpdateStatus(ctx, bookingID, model.StatusCancelled); err != nil {
		return err
	}

	isMasterCancelling := b.MasterID == requesterID && b.ClientID != requesterID

	go func() {
		bgCtx := context.Background()

		if isMasterCancelling {
			altMasterID, findErr := s.Repo.FindAlternativeMaster(
				bgCtx, b.MasterID, b.ServiceID, b.StartsAt, b.EndsAt,
			)

			if findErr != nil || altMasterID == nil {
				fmt.Println("[cancel] альтернатива не найдена — отправляем отмену")
				if notif != nil {
					s.Notifier.BookingCancelled(bgCtx, notif)
				}
				return
			}

			newBooking, createErr := s.Repo.CreateBooking(bgCtx, model.Booking{
				ClientID:  b.ClientID,
				MasterID:  *altMasterID,
				ServiceID: b.ServiceID,
				StartsAt:  b.StartsAt,
				EndsAt:    b.EndsAt,
				Status:    model.StatusConfirmed,
				PricePaid: b.PricePaid,
				Notes:     b.Notes,
			})
			if createErr != nil {
				fmt.Printf("[cancel] ошибка создания новой записи: %v\n", createErr)
				if notif != nil {
					s.Notifier.BookingCancelled(bgCtx, notif)
				}
				return
			}

			if saveErr := s.Repo.SaveAutoReschedule(
				bgCtx, bookingID, newBooking.ID, b.ClientID,
			); saveErr != nil {
				fmt.Printf("[cancel] SaveAutoReschedule error: %v\n", saveErr)
			}

			newNotif, notifErr := s.Repo.GetBookingNotification(bgCtx, newBooking.ID)
			if notifErr == nil && newNotif != nil {
				s.Notifier.BookingAutoRescheduled(bgCtx, newNotif)
			}

			fmt.Printf("[cancel] запись перенесена к мастеру %s\n", *altMasterID)

		} else {
			if notif != nil {
				s.Notifier.BookingCancelled(bgCtx, notif)
			}
		}
	}()

	return nil
}



func (s *BookingService) Reschedule(
	ctx context.Context,
	bookingID uuid.UUID,
	requesterID uuid.UUID,
	req model.RescheduleRequest,
) (*model.Booking, error) {

	b, err := s.Repo.GetBooking(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.ClientID != requesterID {
		return nil, ErrForbidden
	}
	if b.Status == model.StatusCompleted || b.Status == model.StatusCancelled {
		return nil, errors.New("cannot_reschedule")
	}

	targetMasterID := b.MasterID
	if req.NewMasterID != nil {
		targetMasterID = *req.NewMasterID
	}

	services, err := s.Repo.GetBookingServiceIDs(ctx, bookingID)
	if err != nil || len(services) == 0 {
		services = []uuid.UUID{b.ServiceID}
	}

	var totalDuration int
	for _, svcID := range services {
		_, dur, err := s.Repo.GetServicePrice(ctx, targetMasterID, svcID)
		if err != nil {
			return nil, ErrServiceMismatch
		}
		totalDuration += dur
	}

	newEndsAt := req.NewStartsAt.Add(time.Duration(totalDuration) * time.Minute)

	available, err := s.Repo.IsSlotAvailableExclude(ctx, targetMasterID, req.NewStartsAt, newEndsAt, bookingID)
	if err != nil || !available {
		return nil, ErrSlotUnavailable
	}

	var totalPrice float64
	for _, svcID := range services {
		price, _, _ := s.Repo.GetServicePrice(ctx, targetMasterID, svcID)
		totalPrice += price
	}
	discountedPrice := totalPrice * 0.70

	err = s.Repo.RescheduleBooking(
		ctx, bookingID, req.NewMasterID,
		req.NewStartsAt, newEndsAt, discountedPrice,
	)
	if err != nil {
		return nil, err
	}

	s.Repo.SaveDiscount(ctx, model.Discount{
		BookingID: bookingID,
		Type:      "manual",
		Percent:   30,
		Amount:    totalPrice * 0.30,
	})

	updated, err := s.Repo.GetBooking(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	go func() {
		notif, err := s.Repo.GetBookingNotification(context.Background(), bookingID)
		if err != nil {
			return
		}
		s.Notifier.BookingRescheduled(context.Background(), notif)
	}()

	return updated, nil
}