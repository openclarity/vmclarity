package server

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/openclarity/vmclarity/scanner/server/store"
	"github.com/openclarity/vmclarity/scanner/types"
	"net/http"
	"time"
)

// TODO: this needs to be redone, overcomplicated :(

func (s *Server) SubmitScanEvent(ctx echo.Context, scanID types.ScanID) error {
	// Load request
	var event types.ScanEvent
	if err := ctx.Bind(&event); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to bind request: %v", err))
	}

	// Get scan
	scan, err := s.store.Scans().Get(scanID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return sendError(ctx, http.StatusNotFound, err.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// If the scan has finished (with/without errors, we dont care), no new events can be accepted
	switch scan.Status.State {
	case types.ScanStatusStateAborted, types.ScanStatusStateDone, types.ScanStatusStateFailed:
		return sendError(ctx, http.StatusBadRequest,
			fmt.Sprintf("cannot submit events for finished scans"))
	}

	// Extract scan event data
	eventInfo, err := event.EventInfo.ValueByDiscriminator()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	// Handle event
	switch eventData := eventInfo.(type) {
	case types.ScannerHandshakeEventInfo:
		return s.handleHandshake(ctx, scan, eventData)
	case types.ScannerFindingsEventInfo:
		return s.handleFindings(ctx, scan, eventData)
	case types.ScannerHeartbeatEventInfo:
		return s.handleHeartbeat(ctx, scan, eventData)
	default:
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) handleHandshake(ctx echo.Context, scan types.Scan, event types.ScannerHandshakeEventInfo) error {
	// Validate state of the scan
	switch scan.Status.State {
	case types.ScanStatusStatePending: // ok
	default: // scan in progress, not ok
		return sendError(ctx, http.StatusBadRequest,
			fmt.Sprintf("cannot send handshake event for non-Pending scans"))
	}

	// Update scan
	now := time.Now()
	handshakeMsg := "handshake between server and scanner succeeded; scan running"
	scan, err := s.store.Scans().Update(*scan.Id, types.Scan{
		Scanner: &types.ScannerInfo{
			Annotations: event.Annotations,
			Name:        event.Name,
		},
		StartTime: &now,
		Status: &types.ScanStatus{
			LastTransitionTime: time.Now(),
			Message:            &handshakeMsg,
			State:              types.ScanStatusStateInProgress,
		},
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}

func (s *Server) handleFindings(ctx echo.Context, scan types.Scan, event types.ScannerFindingsEventInfo) error {
	// Validate state of the scan
	switch scan.Status.State {
	case types.ScanStatusStateInProgress: // ok
	default: // scan pending, not ok
		return sendError(ctx, http.StatusBadRequest,
			fmt.Sprintf("perform handshake before submitting findings"))
	}

	// Add findings for scan
	_, err := s.store.ScanFindings().CreateMany(*scan.Id, event.Findings...)
	if err != nil {
		var checkErr *store.PreconditionFailedError
		if errors.As(err, &checkErr) {
			return sendError(ctx, http.StatusBadRequest, checkErr.Error())
		}
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}

func (s *Server) handleHeartbeat(ctx echo.Context, scan types.Scan, event types.ScannerHeartbeatEventInfo) error {
	// Validate state of the scan
	switch scan.Status.State {
	case types.ScanStatusStateInProgress: // ok
	default: // scan pending, not ok
		return sendError(ctx, http.StatusBadRequest,
			fmt.Sprintf("perform handshake before submitting heartbeats"))
	}

	// Check event state
	now := time.Now()
	statusMsg := event.Message
	statusState := scan.Status.State
	transitionTime := time.Now()
	endTime := (*time.Time)(nil)
	switch event.State {
	case types.ScannerHeartbeatOK:
		statusMsg = new(string)
		*statusMsg = "heartbeat received"

	case types.ScannerHeartbeatCancelled:
		statusState = types.ScanStatusStateAborted
		endTime = &now

	case types.ScannerHeartbeatErrored:
		statusState = types.ScanStatusStateFailed
		endTime = &now

	case types.ScannerHeartbeatCompleted:
		statusState = types.ScanStatusStateDone
		endTime = &now
	}

	// Add event summary data to scan if supplied
	summary := scan.Summary
	if event.Summary != nil {
		if summary == nil {
			summary = &types.ScanSummary{}
		}
		summary.Add(event.Summary)
	}

	// Update scan
	scan, err := s.store.Scans().Update(*scan.Id, types.Scan{
		EndTime: endTime,
		Status: &types.ScanStatus{
			LastTransitionTime: transitionTime,
			Message:            statusMsg,
			State:              statusState,
		},
		Summary: summary,
	})
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, err.Error())
	}

	return sendResponse(ctx, http.StatusCreated, scan)
}
