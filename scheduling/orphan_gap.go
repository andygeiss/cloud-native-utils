package scheduling

// OrphanGapReason represents the reason for an orphan gap.
type OrphanGapReason string

// Orphan gap reason keys (i18n keys).
const (
	OrphanGapReasonBetweenBookings OrphanGapReason = "reason.orphan_gap.between_bookings"
	OrphanGapReasonEndOfWindow     OrphanGapReason = "reason.orphan_gap.end_of_window"
	OrphanGapReasonStartOfWindow   OrphanGapReason = "reason.orphan_gap.start_of_window"
)

// OrphanGap represents an unusable 30-minute gap created by a selection.
type OrphanGap struct {
	// EndTime is the end time of the gap.
	EndTime TimeOfDay `json:"end_time"`
	// Reason is the i18n key explaining why this gap is an orphan.
	Reason OrphanGapReason `json:"reason"`
	// StartTime is the start time of the gap.
	StartTime TimeOfDay `json:"start_time"`
}

// OrphanGapDetector detects orphan gaps created by a selection.
type OrphanGapDetector struct {
	// dayHours defines the bookable window (opening hours).
	dayHours *DayHours
	// existingBookings are the already-booked time ranges.
	existingBookings []*TimeRange
	// slotDuration is the duration of each slot in minutes.
	slotDuration int
}

// NewOrphanGapDetector creates a new orphan gap detector.
func NewOrphanGapDetector(dayHours *DayHours, existingBookings []*TimeRange, slotDuration int) *OrphanGapDetector {
	return &OrphanGapDetector{
		dayHours:         dayHours,
		existingBookings: existingBookings,
		slotDuration:     slotDuration,
	}
}

// Detect finds all orphan gaps that would be created by the given selection.
// An orphan gap is a 30-minute slot that becomes unbookable because:
// 1. It's between the selection and the start of the booking window
// 2. It's between the selection and the end of the booking window
// 3. It's between the selection and an existing booking
//
// Returns nil if no orphan gaps are detected.
func (a *OrphanGapDetector) Detect(selection *TimeRange) []*OrphanGap {
	if selection == nil {
		return nil
	}

	gaps := make([]*OrphanGap, 0)

	// Check for gap at the start side (window start or booking before selection)
	startGap := a.detectStartOfWindowGap(selection)
	if startGap != nil {
		gaps = append(gaps, startGap)
	}

	// Check for gap at the end side (window end or booking after selection)
	endGap := a.detectEndOfWindowGap(selection)
	if endGap != nil {
		gaps = append(gaps, endGap)
	}

	if len(gaps) == 0 {
		return nil
	}

	return gaps
}

// detectStartOfWindowGap checks if there's an orphan gap between
// the start of the booking window and the selection.
// Only checks the gap from window start to selection start (ignores bookings).
func (a *OrphanGapDetector) detectStartOfWindowGap(selection *TimeRange) *OrphanGap {
	if a.dayHours.IsClosed {
		return nil
	}

	windowStart := a.dayHours.OpenTime
	selectionStart := selection.StartTime()

	// Find the nearest boundary before the selection
	// This could be the window start or the end of a booking
	nearestBefore := windowStart
	for _, booking := range a.existingBookings {
		if booking.EndTime().After(nearestBefore) && !booking.EndTime().After(selectionStart) {
			nearestBefore = booking.EndTime()
		}
	}

	// Calculate gap between nearest boundary and selection start
	gapMinutes := selectionStart.TotalMinutes() - nearestBefore.TotalMinutes()

	// An orphan gap exists if the gap is exactly one slot duration
	if gapMinutes == a.slotDuration {
		// Determine reason based on whether the boundary is window start or booking end
		reason := OrphanGapReasonStartOfWindow
		if !nearestBefore.Equal(windowStart) {
			reason = OrphanGapReasonBetweenBookings
		}
		return &OrphanGap{
			EndTime:   selectionStart,
			Reason:    reason,
			StartTime: nearestBefore,
		}
	}

	return nil
}

// detectEndOfWindowGap checks if there's an orphan gap between
// the selection and the end of the booking window.
func (a *OrphanGapDetector) detectEndOfWindowGap(selection *TimeRange) *OrphanGap {
	if a.dayHours.IsClosed {
		return nil
	}

	windowEnd := a.dayHours.CloseTime
	selectionEnd := selection.EndTime()

	// Find the nearest boundary after the selection
	// This could be the window end or the start of a booking
	nearestAfter := windowEnd
	for _, booking := range a.existingBookings {
		if booking.StartTime().Before(nearestAfter) && !booking.StartTime().Before(selectionEnd) {
			nearestAfter = booking.StartTime()
		}
	}

	// Calculate gap between selection end and nearest boundary
	gapMinutes := nearestAfter.TotalMinutes() - selectionEnd.TotalMinutes()

	// An orphan gap exists if the gap is exactly one slot duration
	if gapMinutes == a.slotDuration {
		// Determine reason based on whether the boundary is window end or booking start
		reason := OrphanGapReasonEndOfWindow
		if !nearestAfter.Equal(windowEnd) {
			reason = OrphanGapReasonBetweenBookings
		}
		return &OrphanGap{
			EndTime:   nearestAfter,
			Reason:    reason,
			StartTime: selectionEnd,
		}
	}

	return nil
}

// HasOrphanGaps returns true if the selection would create any orphan gaps.
func (a *OrphanGapDetector) HasOrphanGaps(selection *TimeRange) bool {
	gaps := a.Detect(selection)
	return len(gaps) > 0
}

// GetFirstOrphanGap returns the first detected orphan gap, or nil if none.
func (a *OrphanGapDetector) GetFirstOrphanGap(selection *TimeRange) *OrphanGap {
	gaps := a.Detect(selection)
	if len(gaps) == 0 {
		return nil
	}
	return gaps[0]
}
