package entity

import "time"

type Subscription struct {
	ID          int64      `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int64      `db:"price"`
	UserID      string     `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type SubscriptionFilter struct {
	UserID      *string
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
}
