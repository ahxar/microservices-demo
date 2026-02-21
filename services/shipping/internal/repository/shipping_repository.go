package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type ShippingRepository struct {
	db *sql.DB
}

type Shipment struct {
	ID                string
	OrderID           string
	TrackingNumber    string
	Carrier           string
	Service           string
	Status            string
	FromStreet        string
	FromCity          string
	FromState         string
	FromZip           string
	FromCountry       string
	ToStreet          string
	ToCity            string
	ToState           string
	ToZip             string
	ToCountry         string
	WeightGrams       int32
	ShippingCostCents int64
	Currency          string
	EstimatedDays     sql.NullInt32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TrackingEvent struct {
	ID          string
	ShipmentID  string
	Status      string
	Location    string
	Description string
	CreatedAt   time.Time
}

func NewShippingRepository(databaseURL string) (*ShippingRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ShippingRepository{db: db}, nil
}

func (r *ShippingRepository) Close() error {
	return r.db.Close()
}

func (r *ShippingRepository) CreateShipment(shipment *Shipment) (*Shipment, error) {
	query := `
		INSERT INTO shipments (
			order_id, tracking_number, carrier, service, status,
			from_street, from_city, from_state, from_zip, from_country,
			to_street, to_city, to_state, to_zip, to_country,
			weight_grams, shipping_cost_cents, currency, estimated_days
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		shipment.OrderID, shipment.TrackingNumber, shipment.Carrier, shipment.Service, shipment.Status,
		shipment.FromStreet, shipment.FromCity, shipment.FromState, shipment.FromZip, shipment.FromCountry,
		shipment.ToStreet, shipment.ToCity, shipment.ToState, shipment.ToZip, shipment.ToCountry,
		shipment.WeightGrams, shipment.ShippingCostCents, shipment.Currency, shipment.EstimatedDays,
	).Scan(&shipment.ID, &shipment.CreatedAt, &shipment.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	return shipment, nil
}

func (r *ShippingRepository) GetShipmentByTracking(trackingNumber string) (*Shipment, error) {
	shipment := &Shipment{}
	query := `
		SELECT id, order_id, tracking_number, carrier, service, status,
			   from_street, from_city, from_state, from_zip, from_country,
			   to_street, to_city, to_state, to_zip, to_country,
			   weight_grams, shipping_cost_cents, currency, estimated_days,
			   created_at, updated_at
		FROM shipments
		WHERE tracking_number = $1
	`

	err := r.db.QueryRow(query, trackingNumber).Scan(
		&shipment.ID, &shipment.OrderID, &shipment.TrackingNumber, &shipment.Carrier, &shipment.Service, &shipment.Status,
		&shipment.FromStreet, &shipment.FromCity, &shipment.FromState, &shipment.FromZip, &shipment.FromCountry,
		&shipment.ToStreet, &shipment.ToCity, &shipment.ToState, &shipment.ToZip, &shipment.ToCountry,
		&shipment.WeightGrams, &shipment.ShippingCostCents, &shipment.Currency, &shipment.EstimatedDays,
		&shipment.CreatedAt, &shipment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shipment not found")
		}
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return shipment, nil
}

func (r *ShippingRepository) GetShipmentByID(shipmentID string) (*Shipment, error) {
	shipment := &Shipment{}
	query := `
		SELECT id, order_id, tracking_number, carrier, service, status,
			   from_street, from_city, from_state, from_zip, from_country,
			   to_street, to_city, to_state, to_zip, to_country,
			   weight_grams, shipping_cost_cents, currency, estimated_days,
			   created_at, updated_at
		FROM shipments
		WHERE id = $1
	`

	err := r.db.QueryRow(query, shipmentID).Scan(
		&shipment.ID, &shipment.OrderID, &shipment.TrackingNumber, &shipment.Carrier, &shipment.Service, &shipment.Status,
		&shipment.FromStreet, &shipment.FromCity, &shipment.FromState, &shipment.FromZip, &shipment.FromCountry,
		&shipment.ToStreet, &shipment.ToCity, &shipment.ToState, &shipment.ToZip, &shipment.ToCountry,
		&shipment.WeightGrams, &shipment.ShippingCostCents, &shipment.Currency, &shipment.EstimatedDays,
		&shipment.CreatedAt, &shipment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shipment not found")
		}
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return shipment, nil
}

func (r *ShippingRepository) GetShipmentByOrderID(orderID string) (*Shipment, error) {
	shipment := &Shipment{}
	query := `
		SELECT id, order_id, tracking_number, carrier, service, status,
			   from_street, from_city, from_state, from_zip, from_country,
			   to_street, to_city, to_state, to_zip, to_country,
			   weight_grams, shipping_cost_cents, currency, estimated_days,
			   created_at, updated_at
		FROM shipments
		WHERE order_id = $1
	`

	err := r.db.QueryRow(query, orderID).Scan(
		&shipment.ID, &shipment.OrderID, &shipment.TrackingNumber, &shipment.Carrier, &shipment.Service, &shipment.Status,
		&shipment.FromStreet, &shipment.FromCity, &shipment.FromState, &shipment.FromZip, &shipment.FromCountry,
		&shipment.ToStreet, &shipment.ToCity, &shipment.ToState, &shipment.ToZip, &shipment.ToCountry,
		&shipment.WeightGrams, &shipment.ShippingCostCents, &shipment.Currency, &shipment.EstimatedDays,
		&shipment.CreatedAt, &shipment.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shipment not found")
		}
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return shipment, nil
}

func (r *ShippingRepository) UpdateShipmentStatus(shipmentID, status string) error {
	query := `UPDATE shipments SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, status, shipmentID)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}
	return nil
}

func (r *ShippingRepository) AddTrackingEvent(event *TrackingEvent) error {
	query := `
		INSERT INTO tracking_events (shipment_id, status, location, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		event.ShipmentID, event.Status, event.Location, event.Description,
	).Scan(&event.ID, &event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add tracking event: %w", err)
	}

	return nil
}

func (r *ShippingRepository) GetTrackingEvents(shipmentID string) ([]TrackingEvent, error) {
	query := `
		SELECT id, shipment_id, status, location, description, created_at
		FROM tracking_events
		WHERE shipment_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking events: %w", err)
	}
	defer rows.Close()

	var events []TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		err := rows.Scan(
			&event.ID, &event.ShipmentID, &event.Status,
			&event.Location, &event.Description, &event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
