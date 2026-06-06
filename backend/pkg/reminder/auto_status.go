package reminder

import (
	"context"
	"fmt"
)

func (r *ReminderJob) AutoUpdateBookingStatuses() {
	ctx := context.Background()

	
	res1, err := r.db.Exec(ctx, `
		update public.bookings
		set status = 'cancelled'
		where status = 'pending'
		  and ends_at < now()
	`)
	if err != nil {
		fmt.Printf("[auto-status] ошибка отмены pending: %v\n", err)
	} else {
		fmt.Printf("[auto-status] pending - cancelled: %d записей\n", res1.RowsAffected())
	}


	res2, err := r.db.Exec(ctx, `
		update public.bookings
		set status = 'completed'
		where status = 'confirmed'
		  and ends_at < now()
	`)
	if err != nil {
		fmt.Printf("[auto-status] ошибка завершения confirmed: %v\n", err)
	} else {
		fmt.Printf("[auto-status] confirmed - completed: %d записей\n", res2.RowsAffected())
	}
}